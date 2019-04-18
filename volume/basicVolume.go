package volume

import (
    "bytes"
    "encoding/json"
    "html/template"
    "io"
    "os"
    "path/filepath"
    "strings"
    "time"

    "../cache"
    "../i18n"
    "../util"
)

type BasicVolume struct {
    assets         map[string] *Asset
    i18n           *i18n.I18n
    localePath     map[string] string
    localePathMx   *util.MapMutex
    templates      *template.Template
    templatesClone *template.Template
}

func NewBasicVolume() *BasicVolume {
    vol := &BasicVolume{}
    vol.Reset()
    return vol
}

// Getters

func( vol *BasicVolume ) Asset( path string ) ( *Asset ) {
    ast := vol.assets[path]
    return ast
}

func( vol *BasicVolume ) Localizer( lcName string ) ( *i18n.Localizer ) {
    return vol.i18n.Localizer( lcName )
}

func( vol *BasicVolume ) Template( path string ) ( *template.Template ) {

    /*/
     +  Note that html/template#Template.Lookup( "" ) returns itself while
     + text/template#Template.Lookup( "" ) returns nil
     + provided the template was created as t := template.New( "" )
    /*/

    if tmpl := vol.templates.Lookup( path ); tmpl != nil {
        return tmpl
    }
    return nil
}

func( vol *BasicVolume ) I18n() *i18n.I18n {
    return vol.i18n
}

// General

func( vol *BasicVolume ) Reset() {
    vol.assets            = make( map[string] *Asset )
    vol.i18n              = i18n.New()
    vol.localePath        = make( map[string] string )
    vol.localePathMx      = util.NewMapMutex()
    vol.templates         = template.New( "" )
    vol.templatesClone, _ = vol.templates.Clone()
}

// Importers

func( vol *BasicVolume ) PutItem( path string, rd io.Reader, modTime time.Time ) error {

    switch typeByPath( path ) {
    case c_typeAsset:

        // Add
        ast := NewAsset( filepath.Base( path ), rd, modTime )
        vol.putAsset( path, ast )
        return nil

    case c_typeI18n:

        // Read
        buf := bytes.NewBuffer( nil )
        io.Copy( buf, rd )

        // Parse
        lc := i18n.NewLocale()
        err := lc.ParseJson( buf.Bytes() )
        if err != nil {
            return err
        }

        // Path
        vol.localePathMx.BeginWrite()
        vol.localePath[lc.Name()] = path
        vol.localePathMx.EndWrite()

        // Assign
        vol.PutLocale( lc )
        return nil

    case c_typeTemplate:

        // Read
        buf := bytes.NewBuffer( nil )
        io.Copy( buf, rd )

        // Add
        return vol.PutTemplate( path, buf.String() )

    default:
        return ErrInvalidPath.Append( path )
    }
    return nil

}

func( vol *BasicVolume ) PutAsset( path string, ast *Asset ) error {
    if !strings.HasPrefix( path, c_assetDirectory ) {
        return ErrInvalidPath.Append( path )
    }
    vol.putAsset( path, ast )
    return nil
}

func( vol *BasicVolume ) putAsset( path string, ast *Asset ) {
    buf := make( map[string] *Asset )
    for k, v := range vol.assets {
        buf[k] = v
    }
    buf[path] = ast
    vol.assets = buf
}

func( vol *BasicVolume ) PutLocale( lc *i18n.Locale ) {
    vol.i18n.PutLocale( lc )
}

func( vol *BasicVolume ) PutTemplate( path string, text string ) error {

    // Check path
    if !strings.HasPrefix( path, c_templateDirectory ) {
        return ErrInvalidPath.Append( path )
    }

    // Clone
    _, err := vol.templatesClone.New( path ).Parse( text )
    if err != nil {
        return ErrInvalidTemplate.Append( path, err )
    }

    buf, _ := vol.templatesClone.Clone()
    vol.templates = buf
    return nil

}

// Filesystem Related
//
// * Notice how the methods use osPath instaed of path to differentiate paths in volume and paths in operating system.

/*
 + walk
 *
 * filepath.WalkFunc
 * relpath must be volume-path-ready e.g.) templates/
 */

func( vol *BasicVolume ) ImportDirectory( path string ) error {
    return filepath.Walk( path, vol.walk( path ) )
}

func( vol *BasicVolume ) walk( basePath string ) filepath.WalkFunc {

    return func( osPath string, fi os.FileInfo, err error ) error {

        // Ignore dir
        if fi.IsDir() {
            return nil
        }

        // Rel
        relPath, relErr := filepath.Rel( basePath, osPath )
        if relErr != nil {
            return ErrInvalidPath.Append( relErr, "basePath:", basePath, "osPath:", osPath )
        }
        relPath = filepath.ToSlash( relPath )

        // Open
        f, err := os.Open( osPath )
        if err != nil {
            return ErrFileError.Append( osPath )
        }

        // Add and ignore invalid path error
        err = vol.PutItem( relPath, f, fi.ModTime() )
        if cast, ok := err.( util.Error ); ok {
            if !cast.Is( ErrInvalidPath ) {
                return err
            }
        } else {
            return err
        }

        // Close
        err = f.Close()
        if err != nil {
            return ErrFileError.Append( err )
        }

        return nil

    }

}

// cache.Porter

func( vol *BasicVolume ) Export() ( *cache.Cache, error ) {

    chc  := cache.NewCache()
    zero := time.Time{}

    // Asset
    for path, ast := range vol.assets {

        w, err := chc.Create( path, ast.modTime )
        if err != nil {
            return nil, ErrAssetExport.Append( path, err )
        }
        w.Write( ast.Bytes() )
        w.Close()

    }

    // i18n
    copy := make( map[string] string )
    vol.localePathMx.BeginRead()
    for k, v := range vol.localePath {
        copy[k] = v
    }
    vol.localePathMx.EndRead()

    for lcName, path := range copy {

        // Create
        w, err := chc.Create( path, zero )
        if err != nil {
            return nil, ErrI18nExport.Append( lcName, err )
        }

        //
        lc, ok := vol.i18n.Locale( lcName )
        if !ok {
            return nil, ErrI18nExport.Append( lcName, "the locale is unavailable" )
        }

        // Json
        jenc := json.NewEncoder( w )
        lcMap := map[string] interface{} {
            lcName: lc.Set(),
        }
        jenc.Encode( lcMap )
        w.Close()

    }

    // Template
    for _, tmpl := range vol.templates.Templates() {

        if tmpl.Name() == "" {
            /*
             *  Keep that in mind that text/template.Template.Templates() and html/template.Template.Templates()
             * work differently from each other. The html one includes itself while the other doesn't.
             */
            continue
        }

        // Create
        w, err := chc.Create( tmpl.Name(), zero )
        if err != nil {
            return nil, ErrTemplateExport.Append( tmpl.Name(), err )
        }

        println( tmpl.Name() )
        w.Write( []byte( tmpl.Tree.Root.String() ) )
        w.Close()

    }

    chc.Flush()
    return chc, nil

}

func( vol *BasicVolume ) Import( chc *cache.Cache ) error {

    for _, f := range chc.Files() {

        // Open
        rc, err := f.Open()
        if err != nil {
            return ErrZipImport.Append( f.Name, err )
        }

        // Put
        err = vol.PutItem( f.Name, rc, f.Modified )
        if err != nil {
            return ErrZipImport.Append( err )
        }

    }

    return nil

}