package volume

import (
    "bytes"
    "encoding/json"
    "html/template"
    "io"
    "os"
    "path/filepath"
    "strings"
    "sync"
    "time"

    "github.com/hjjg200/sprout/cache"
    "github.com/hjjg200/sprout/environ"
    "github.com/hjjg200/sprout/i18n"
    "github.com/hjjg200/sprout/util/errors"
)

type BasicVolume struct {
    // assets
    //  wrties: seldom
    //  reads: frequent
    //  size: big
    //  => RWMutex
    // localePath
    //  writes: seldom
    //  reads: seldom
    //  size: small
    //  => buffer
    assets         map[string] *Asset
    assetsMu       sync.RWMutex
    i18n           *i18n.I18n
    localePath     map[string] string
    templates      *template.Template
    templatesClone *template.Template
    fallback       Volume
}

func NewBasicVolume() *BasicVolume {
    vol := &BasicVolume{}
    vol.Reset()
    vol.SetFallback( DefaultVolume )
    return vol
}

func newBasicVolume() *BasicVolume {
    vol := &BasicVolume{}
    vol.Reset()
    return vol
}

// Getters

func( vol *BasicVolume ) Asset( path string ) ( *Asset ) {

    // Lock
    vol.assetsMu.RLock()
    defer vol.assetsMu.RUnlock()

    ast := vol.assets[path]
    if ast == nil && vol.fallback != nil {
        ast = vol.fallback.Asset( path )
    }
    return ast

}

func( vol *BasicVolume ) Localizer( lcName string ) ( *i18n.Localizer ) {
    lczr := vol.i18n.Localizer( lcName )
    if lczr == nil && vol.fallback != nil {
        lczr = vol.fallback.Localizer( lcName )
    }
    return lczr
}

func( vol *BasicVolume ) Template( path string ) ( *template.Template ) {

    /*/
     +  Note that html/template#Template.Lookup( "" ) returns itself while
     + text/template#Template.Lookup( "" ) returns nil
     + provided the template was created as t := template.New( "" )
    /*/

    tmpl := vol.templates.Lookup( path )
    if tmpl == nil && vol.fallback != nil {
        tmpl = vol.fallback.Template( path )
    }
    return tmpl
}

func( vol *BasicVolume ) I18n() *i18n.I18n {
    i1 := vol.i18n
    if i1 == nil && vol.fallback != nil {
        i1 = vol.fallback.I18n()
    }
    return i1
}

func( vol *BasicVolume ) SetFallback( flb Volume ) {
    vol.fallback = flb
}

// General

func( vol *BasicVolume ) Reset() {

    // Members
    vol.assets            = make( map[string] *Asset )
    vol.i18n              = i18n.New()
    vol.localePath        = make( map[string] string )
    vol.templates         = template.New( "" )
    vol.templatesClone, _ = vol.templates.Clone()
    vol.fallback          = nil

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
            return errors.ErrInvalidObject.Raise( err )
        }

        // Path
        buf2 := make( map[string] string )
        for k, v := range vol.localePath {
            buf2[k] = v
        }
        buf2[lc.Name()] = path
        vol.localePath = buf2

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
        return errors.ErrInvalidPath.Raise( path )
    }
    return nil

}

func( vol *BasicVolume ) PutAsset( path string, ast *Asset ) error {
    if !strings.HasPrefix( path, c_assetDirectory ) {
        return errors.ErrInvalidPath.Raise( path )
    }
    vol.putAsset( path, ast )
    return nil
}

func( vol *BasicVolume ) putAsset( path string, ast *Asset ) {

    // Put
    vol.assetsMu.Lock()
    vol.assets[path] = ast
    vol.assetsMu.Unlock()

    // Compile
    cmp, ok := DefaultCompilers.OutputOf( path )
    if ok {
        cmpAst, err := DefaultCompilers.Compile( ast )
        if err != nil {
            environ.Logger.Warnln( errors.ErrCompileFailure.Raise( err ) )
        }
        vol.putAsset( cmp, cmpAst )
    }

}

func( vol *BasicVolume ) PutLocale( lc *i18n.Locale ) {
    vol.i18n.PutLocale( lc )
}

func( vol *BasicVolume ) PutTemplate( path string, text string ) error {

    // Check path
    if !strings.HasPrefix( path, c_templateDirectory ) {
        return errors.ErrInvalidPath.Raise( path )
    }

    // Clone
    _, err := vol.templatesClone.New( path ).Parse( text )
    if err != nil {
        return errors.ErrInvalidObject.Raise( "invalid template", path, err )
    }

    // Put
    vol.templates, _ = vol.templatesClone.Clone()

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
    path = filepath.ToSlash( filepath.Clean( path ) )
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
            return errors.ErrInvalidPath.Raise( "relErr:", relErr, "basePath:", basePath, "osPath:", osPath )
        }
        relPath = filepath.ToSlash( relPath )

        // Open
        f, err := os.Open( osPath )
        if err != nil {
            return errors.ErrIOError.Raise( osPath )
        }

        // Add and ignore invalid path error
        err = vol.PutItem( relPath, f, fi.ModTime() )
        if !errors.ErrInvalidPath.Is( err ) {
            // Ignore invalid path errors
            return err
        }

        // Close
        err = f.Close()
        if err != nil {
            return errors.ErrIOError.Raise( err )
        }

        return nil

    }

}

//

func( vol *BasicVolume ) HasItem( path string ) bool {

    for p := range vol.assets {
        if p == path { return true }
    }
    for _, p := range vol.localePath {
        if p == path { return true }
    }
    for _, t := range vol.templatesClone.Templates() {
        if t.Name() == path { return true }
    }
    return false

}

func( vol *BasicVolume ) RemoveItem( path string ) error {

    switch typeByPath( path ) {
    case c_typeAsset:

        if _, ok := vol.assets[path]; ok {
            vol.assetsMu.Lock()
            delete( vol.assets, path )
            vol.assetsMu.Unlock()
        }
        return nil

    case c_typeI18n:

        for lcName, p := range vol.localePath {
            if p == path {
                vol.i18n.RemoveLocale( lcName )
            }
        }
        return nil

    case c_typeTemplate:

        removed := template.New( "" )
        for _, t := range vol.templatesClone.Templates() {
            if t.Name() != path {
                if tr := t.Tree; tr != nil {
                    removed.New( t.Name() ).Parse( tr.Root.String() )
                }
            }
        }
        vol.templatesClone, _ = removed.Clone()
        vol.templates, _      = vol.templatesClone.Clone()
        return nil

    default:
        return errors.ErrInvalidPath.Raise( path )
    }
    return errors.ErrNotFound.Raise( path )

}

// cache.Porter

func( vol *BasicVolume ) Export() ( *cache.Cache, error ) {

    chc  := cache.NewCache()
    zero := time.Time{}

    // Asset
    for path, ast := range vol.assets {

        w, err := chc.Create( path, ast.modTime )
        if err != nil {
            return nil, errors.ErrExportFailure.Raise( path, err )
        }
        w.Write( ast.Bytes() )
        w.Close()

    }

    // i18n
    copy := make( map[string] string )
    for k, v := range vol.localePath {
        copy[k] = v
    }

    for lcName, path := range copy {

        // Create
        w, err := chc.Create( path, zero )
        if err != nil {
            return nil, errors.ErrExportFailure.Raise( lcName, err )
        }

        //
        lc, ok := vol.i18n.Locale( lcName )
        if !ok {
            return nil, errors.ErrExportFailure.Raise( lcName, "the locale is unavailable" )
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
            return nil, errors.ErrExportFailure.Raise( tmpl.Name(), err )
        }

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
            return errors.ErrImportFailure.Raise( f.Name, err )
        }

        // Put
        err = vol.PutItem( f.Name, rc, f.Modified )
        if err != nil {
            return errors.ErrImportFailure.Raise( err )
        }

    }

    return nil

}
