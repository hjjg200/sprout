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

func MakeAsset( _name string, _bytes []byte, _mod_time time.Time ) *Asset {}
func( _asset *Asset ) Serve() Handler {}
func( _asset *Asset ) Bytes() []byte {}
func( _asset *Asset ) Hash() string {}
func( _asset *Asset ) Name() string {}
func( _asset *Asset ) ModTime() time.Time {}
func( _asset *Asset ) MimeType() string {}