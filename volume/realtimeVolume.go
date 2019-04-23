package volume

import (
    "html/template"
    "os"
    "path/filepath"
    "time"

    "github.com/hjjg200/sprout/i18n"
    "github.com/hjjg200/sprout/cache"
    "github.com/hjjg200/sprout/util/errors"
)

type RealtimeVolume struct {
    vol *BasicVolume
    srcPath string
    modTime map[string] time.Time
}

// VVOLUME METHODS

func NewRealtimeVolume( srcPath string ) *RealtimeVolume {

    srcPath = filepath.ToSlash( filepath.Clean( srcPath ) )

    vol := NewBasicVolume()
    vol.ImportDirectory( srcPath )

    return &RealtimeVolume{
        vol: vol,
        srcPath: srcPath,
        modTime: make( map[string] time.Time ),
    }

}

func( rtv *RealtimeVolume ) abs( path string ) string {
    return rtv.srcPath + "/" + path
}

func( rtv *RealtimeVolume ) validate( path string ) error {

    absPath := rtv.abs( path )

    fi, err := os.Stat( absPath )
    if err != nil {
        if os.IsNotExist( err ) {
            // If Compiled
            if in, ok := DefaultCompilers.InputOf( path ); ok {
                var err2 error
                for _, i := range in {
                    err2 = rtv.validate( i )
                    if err2 == nil { return nil }
                }
                return err2
            } else {
                // Remove item if there is any in the underlying volume
                if rtv.vol.HasItem( path ) {
                    return rtv.vol.RemoveItem( path )
                }
            }
            return errors.ErrNotFound.Raise( path )
        }
        return errors.ErrIOError.Raise( path, err )
    }

    mt, ok := rtv.modTime[path]
    if ok {
        if fi.ModTime().Sub( mt ) <= 0 {
            return nil
        }
    }

    // Write to modtTime
    buf := make( map[string] time.Time )
    for k, v := range rtv.modTime { buf[k] = v }
    buf[path] = fi.ModTime()
    rtv.modTime = buf

    // Write
    f, err := os.Open( absPath )
    if err != nil {
        return errors.ErrIOError.Raise( absPath, err )
    }
    defer f.Close()



    // Put
    return rtv.vol.PutItem( path, f, fi.ModTime() )

}

func( rtv *RealtimeVolume ) validateTemplates() error {

    return filepath.Walk( rtv.srcPath + "/" + c_templateDirectory, func( osPath string, fi os.FileInfo, err error ) error {

        // Ignore dir
        if fi.IsDir() {
            return nil
        }

        // Rel
        relPath, relErr := filepath.Rel( rtv.srcPath, osPath )
        if relErr != nil {
            return errors.ErrInvalidPath.Raise( "relErr:", relErr, "basePath:", rtv.srcPath, "osPath:", osPath )
        }
        relPath = filepath.ToSlash( relPath )

        // Add and ignore invalid path error
        err = rtv.validate( relPath )
        if !errors.ErrInvalidPath.Is( err ) {
            return err
        }

        return nil

    } )

}

func( rtv *RealtimeVolume ) validateI18n() error {

    for _, path := range rtv.vol.localePath {
        err := rtv.validate( path )
        if !errors.ErrNotFound.Is( err ) && err != nil {
            return nil
        }
    }
    return nil

}

func( rtv *RealtimeVolume ) walkI18nDirectory() error {

    i18nDir := rtv.abs( c_i18nDirectory )

    { // Ensure the i18n Directory
        fi, err := os.Stat( i18nDir )
        if err != nil {
            return errors.ErrIOError.Raise( i18nDir, err )
        } else if !fi.IsDir() {
            return errors.ErrIOError.Raise( i18nDir, "it is not a directory" )
        }
    }

    return filepath.Walk( i18nDir, func( absPath string, fi os.FileInfo, err error ) error {

        // Rel
        relPath, relErr := filepath.Rel( rtv.srcPath, absPath )
        if relErr != nil {
            return errors.ErrInvalidPath.Raise( relErr, relPath )
        }
        if fi.IsDir() {
            return nil
        }

        relPath = filepath.ToSlash( relPath )

        return rtv.validate( relPath )

    } )

}

func( rtv *RealtimeVolume ) Asset( path string ) ( *Asset ) {
    err := rtv.validate( path )
    if !errors.ErrNotFound.Is( err ) && err != nil {
        return nil
    }
    return rtv.vol.Asset( path )
}

func( rtv *RealtimeVolume ) I18n() ( *i18n.I18n ) {
    err := rtv.walkI18nDirectory()
    if !errors.ErrNotFound.Is( err ) && err != nil {
        return nil
    }
    return rtv.vol.I18n()
}

func( rtv *RealtimeVolume ) Localizer( lcName string ) ( *i18n.Localizer ) {

    // Valiate
    path, ok := rtv.vol.localePath[lcName]
    if ok {
        err := rtv.validate( path )
        if !errors.ErrNotFound.Is( err ) && err != nil {
            return nil
        }
    } else {
        // Walk
        err := rtv.walkI18nDirectory()
        if !errors.ErrNotFound.Is( err ) && err != nil {
            return nil
        }
    }

    return rtv.vol.Localizer( lcName )

}

func( rtv *RealtimeVolume ) Template( path string ) ( *template.Template ) {
    err := rtv.validateTemplates()
    if !errors.ErrNotFound.Is( err ) && err != nil {
        return nil
    }
    return rtv.vol.Template( path )
}

func( rtv *RealtimeVolume ) SetFallback( vol Volume ) {
    rtv.vol.SetFallback( vol )
}

func( rtv *RealtimeVolume ) Export() ( *cache.Cache, error ) {
    return nil, nil
}

func( rtv *RealtimeVolume ) Import( chc *cache.Cache ) error {
    return nil
}
