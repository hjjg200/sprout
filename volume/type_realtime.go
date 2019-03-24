package volume

type RealtimeVolume struct {
    cache   *cache.Cache
    srcPath string
}

// ENTRY

type RealtimeEntry interface {
    Body() interface{}
    Validate()
}

type RealtimeAsset struct {
    srcPath string
    body *Asset
    newBody *Asset
     
}

// VVOLUME METHODS

func NewRealtimeVolume( srcPath string ) *RealtimeVolume {
    
}

func( rtv *RealtimeVolume ) Asset( path string ) *Asset {
    
}

func( rtv *RealtimeVolume ) I18n() *I18n {
    
}

func( rtv *RealtimeVolume ) Template( path string ) *template.Template {

}