package sprout

/*
 + ASSET
 *
 * An asset represents a single file
 */

type Asset struct {
    bytes     []byte
    hash      string
    name      string
    mime_type string
    mod_time  time.Time
}

//func MakeAsset( _name string, _bytes []byte, _mod_time time.Time ) *Asset {}
func( _asset *Asset ) Serve() Handler {

}
func( _asset *Asset ) Bytes() []byte {
    _bytes := make( []byte, len( _asset.bytes ) )
    copy( _bytes, _asset.bytes )
    return _bytes
}
func( _asset *Asset ) Hash() string {
    return _asset.hash
}
func( _asset *Asset ) Name() string {
    return _asset.name
}
func( _asset *Asset ) ModTime() time.Time {
    return _asset.mod_time
}
func( _asset *Asset ) MimeType() string {
    return _asset.mime_type
}