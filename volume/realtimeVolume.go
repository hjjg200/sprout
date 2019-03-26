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

func( rtv *RealtimeVolume ) validate( path string ) bool {
    
}

func( rtv *RealtimeVolume ) Asset( path string ) ( *Asset, error ) {

    var (
        ast *Asset
    )

    err := validateAsset( path )
    if err != nil {
        return nil, err
    }

    ast, _ := rtv.vol.Asset( path )
    return ast, nil

}

func( rtv *RealtimeVolume ) validateAsset( path string ) error {

    absPath := rtv.abs( path )
    base    := filepath.Base( path )

    fi, err := os.Stat( absPath )
    if err != nil {
        return ErrFileError.Append( absPath, err )
    }

    mt, ok := rtv.modTime[path]
    if ok {
        if fi.ModTime().Sub( mt ) <= 0 {
            // No need to validate
            return nil
        }
    }

    // Write to modtTime
    buf := make( map[string] time.Time )
    for k, v := range rtv.modTime { buf[k] = v }
    buf[path] = fi.ModTime()
    rtv.modTime = buf

    f, err := os.Open( absPath )
    if err != nil {
        return ErrFileError.Append( absPath, err )
    }
    defer f.Close()

    ast := NewAsset( base, f, fi.ModTime() )
    return rtv.vol.PutAsset( path, ast )

}