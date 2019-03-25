package volume

type RealtimeVolume struct {
    srcPath string
    entries map[string] interface{}
    newEntries map[string] interface{}
    modTime map[string] time.Time
}

// VVOLUME METHODS

func NewRealtimeVolume( srcPath string ) *RealtimeVolume {

}

func( rtv *RealtimeVolume ) Asset( path string ) *Asset {

    var (
        ast *Asset
    )

    //
    if entry, ok := rtv.entries[path]; ok && entry != nil {
        ast = entry.( *Asset )
        go validateAsset( path )
        return ast
    }

    validateAsset( path )
    if _, ok := rtv.entries[path]; ok {
        return rtv.entries[path].( *Asset )
    }

    return nil

}

func( rtv *RealtimeVolume ) validateAsset( path string ) {

    absPath := rtv.srcPath + "/" + path
    base    := filepath.Base( path )

    fi, err := os.Stat( absPath )
    if err != nil {
        return
    }

    mt, ok := rtv.modTime[path]
    if !ok {
        return
    }

    if fi.ModTime().Sub( mt ) <= 0 {
        return
    }

    f, err := os.Open( absPath )
    if err != nil {
        return
    }

    entries[path] = NewAsset( base, f, fi.ModTime() )
    f.Close()

}