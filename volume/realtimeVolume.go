package volume

import (
    "html/template"
    "os"
    "path/filepath"
    "time"

    "../i18n"
)

type RealtimeVolume struct {
    vol *Volume
    srcPath string
    modTime map[string] time.Time
}

// VVOLUME METHODS

func NewRealtimeVolume( srcPath string ) *RealtimeVolume {
    return &RealtimeVolume{
        vol: NewVolume(),
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
        return ErrFileError.Append( path, err )
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
        return ErrFileError.Append( absPath, err )
    }
    defer f.Close()

    // Put
    return rtv.vol.PutItem( path, f, fi.ModTime() )

}

func( rtv *RealtimeVolume ) validateI18n() error {

    for _, path := range rtv.vol.localePath {
        err := rtv.validate( path )
        if err != nil {
            return err
        }
    }
    return nil

}

func( rtv *RealtimeVolume ) walkAndValidate( path string ) {

    // MkdirAll

    filepath.Walk( rtv.abs( path ), func( absPath string, fi os.FileInfo, err error ) error {

        // Rel
        relPath, relErr := filepath.Rel( rtv.srcPath, absPath )
        if relErr != nil {
            return relErr
        }

        //
        if fi.IsDir() {
            return nil
        }

        return rtv.validate( relPath )

    } )

}

func( rtv *RealtimeVolume ) Asset( path string ) ( *Asset, bool ) {
    err := rtv.validate( path )
    if err != nil {
        return nil, false
    }
    return rtv.vol.Asset( path )
}

func( rtv *RealtimeVolume ) I18n() ( *i18n.I18n ) {
    err := rtv.validateI18n()
    if err != nil {
        return nil
    }
    return rtv.vol.I18n()
}

func( rtv *RealtimeVolume ) Localizer( lcName string ) ( *i18n.Localizer, bool ) {

    // Locate
    path, ok := rtv.vol.localePath[lcName]
    if !ok {
        rtv.walkAndValidate( c_i18nDirectory )
        path, ok = rtv.vol.localePath[lcName]
        if !ok {
            return nil, false
        }
    }

    // Valiate
    err := rtv.validate( path )
    if err != nil {
        return nil, false
    }
    return rtv.vol.Localizer( lcName )

}

func( rtv *RealtimeVolume ) Template( path string ) ( *template.Template, bool ) {
    err := rtv.validate( path )
    if err != nil {
        return nil, false
    }
    return rtv.vol.Template( path )
}