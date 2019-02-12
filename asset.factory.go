package sprout

/*
 + ASSET FACTORY
 */

type asset_factory struct {}
var  static_asset_factory = &asset_factory{}

func AssetFactory() *asset_factory {
    return static_asset_factory
}
func( _astfac *asset_factory ) New( _name string, _bytes []byte, _mod_time time.Time ) *Asset {

    // Copy the slice as precaution
    _buffer := make( []byte, len( _bytes ) )
    copy( _buffer, _bytes )

    // Evaluate the hash
    _sha256 := sha256.New()
    _sha256.Write( _bytes )
    _hash := fmt.Sprintf( "%x", _sha256.Sum( nil ) )

    // Get the mime type
    _mime_type = mime.TypeByExtension( path.Ext( _name ) )
    if _mime_type == "" {
        _mime_type = "text/plain"
    }

    // Make
    return &Asset{
        bytes: _buffer,
        hash: _hash,
        name: _name,
        mime_type: _mime_type,
        mod_time: _mod_time,
    }

}
func( _astfac *asset_factory ) FromReader( _name string, _reader io.Reader, _mod_time time.Time ) *Asset {

}