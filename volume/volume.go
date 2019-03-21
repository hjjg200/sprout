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

type Volume struct {
    assets map[string] *Asset
    i18n *i18n.I18n
    templates *template.Template
}

func NewVolume() *Volume {
    vol := &Volume{}
    vol.Reset()
    return vol
}

// Getters

func( vol *Volume ) Asset( path string ) ( *Asset, bool ) {
    ast, ok := vol.assets[path]
    return ast, ok
}
func( vol *Volume ) Localizer( lcName string ) ( *i18n.Localizer, bool ) {
    lczr, err := i18n.NewLocalizer( vol.i18n, lcName )
    if err != nil {
        return nil, false
    }
    return lczr, true
}
func( vol *Volume ) Template( path string ) ( *template.Template, bool ) {
    if tmpl := vol.templates.Lookup( path ); tmpl != nil {
        return tmpl, true
    }
    return nil, false
}
func( vol *Volume ) I18n() *i18n.I18n {
    return vol.i18n
}

// General

func( vol *Volume ) Localize( lcName, src string ) string {
    return vol.i18n.Localize( lcName, src )
}
func( vol *Volume ) ExecuteTemplate( path string, data interface{} ) ( string, error ) {

    // Get
    tmpl, ok := vol.Template( path )
    if !ok {
        return "", ErrTemplateNonExistent.Append( path )
    }

    // Read
    buf := bytes.NewBuffer( nil )
    err := tmpl.Execute( buf, data )
    if err != nil {
        return "", ErrTemplateExecError.Append( err )
    }

    return buf.String(), nil

}

func( vol *Volume ) Reset() {
    vol.assets    = make( map[string] *Asset )
    vol.i18n      = i18n.New()
    vol.templates = template.New( "" )
}

// Importers

func( vol *Volume ) PutItem( path string, rd io.Reader, modTime time.Time ) error {

    switch {
    case strings.HasPrefix( path, c_assetDirectory ):
        ast := NewAsset( filepath.Base( path ), rd, modTime )
        return vol.PutAsset( path, ast )
    case strings.HasPrefix( path, c_i18nDirectory ):

        // Read
        buf := bytes.NewBuffer( nil )
        io.Copy( buf, rd )

        // Parse
        lc := i18n.NewLocale()
        err := lc.ParseJson( buf.Bytes() )
        if err != nil {
            return err
        }

        // Assign
        return vol.PutLocale( lc )

    case strings.HasPrefix( path, c_templateDirectory ):

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

func( vol *Volume ) PutAsset( path string, ast *Asset ) error {
    if _, ok := vol.assets[path]; ok {
        return ErrOccupiedPath.Append( path )
    }
    if !strings.HasPrefix( path, c_assetDirectory ) {
        return ErrInvalidPath.Append( path )
    }
    vol.assets[path] = ast
    return nil
}

func( vol *Volume ) PutLocale( lc *i18n.Locale ) error {
    return vol.i18n.AddLocale( lc )
}

func( vol *Volume ) PutTemplate( path string, text string ) error {

    // Check if exists
    if _, ok := vol.Template( path ); ok {
        return ErrOccupiedPath.Append( path )
    }

    // Check path
    if !strings.HasPrefix( path, c_templateDirectory ) {
        return ErrInvalidPath.Append( path )
    }

    // Parse
    _, err := vol.templates.New( path ).Parse( text )
    if err != nil {
        ErrInvalidTemplate.Append( path, err )
    }

    return nil
}

// Filesystem Related
//
// * Notice how the methods use osPath instaed of path to differentiate paths in volume and paths in operating system.

func( vol *Volume ) ImportLocaleDirectory( osPath string ) error {
    return vol.i18n.ImportDirectory( osPath )
}

/*
 + walkFuncBasedOn
 *
 * filepath.WalkFunc
 * relpath must be volume-path-ready e.g.) templates/
 */

func( vol *Volume ) ImportDirectory( path string ) error {
    vol.Reset()
    return filepath.Walk( path, vol.walkFuncBasedOn( path ) )
}

func( vol *Volume ) walkFuncBasedOn( basePath string ) filepath.WalkFunc {

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

func( vol *Volume ) Export() ( *cache.Cache, error ) {

    chc  := cache.NewCache()
    zero := time.Time{}

    // Asset
    for path, ast := range vol.assets {

        w, err := chc.Create( path, ast.modTime )
        if err != nil {
            return nil, ErrAssetExport.Append( path, err )
        }
        ast.Seek( 0, io.SeekStart )
        io.Copy( w, ast )
        w.Close()

    }

    // i18n
    for lcName, lc := range vol.i18n.Locales() {

        // Create
        w, err := chc.Create( c_i18nDirectory + "/" + lcName + ".json", zero )
        if err != nil {
            return nil, ErrI18nExport.Append( lcName, err )
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

func( vol *Volume ) Import( chc *cache.Cache ) error {

    vol.Reset()

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