package volume

type RealtimeVolume struct {
    srcPath string
    modTime map[string] time.Time
}

// VVOLUME METHODS

func NewRealtimeVolume( srcPath string ) *RealtimeVolume {

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
    absPath := rtv.abs( c_i18nDirectory )
    
}

func( rtv *RealtimeVolume ) Asset( path string ) ( *Asset ) {
    err := rtv.validate( path )
    if err != nil {
        return nil
    }
    return rtv.vol.Asset( path )
}

func( rtv *RealtimeVolume ) I18n() ( *i18n.I18n ) {
    
}

func( rtv *RealtimeVolume ) Localizer( lcName string ) ( *i18n.Localizer ) {
    err := rtv.validate( path )
    if err != nil {
        return nil
    }
    return rtv.vol.Template( path )
}

func( rtv *RealtimeVolume ) Template( path string ) ( *template.Template ) {
    err := rtv.validate( path )
    if err != nil {
        return nil
    }
    return rtv.vol.Template( path )
}