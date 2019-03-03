package sprout

/*
 + VOLUME HELPER
 *
 * VolumeHelper contains most of the functions that deal with the os and filesystem
 */

type volume_helper struct{}
var  static_volume_helper = &volume_helper

func VolumeHelper() *volume_helper {
    return static_volume_helper
}
func( _volhlpr *volume_helper )

type volume_cache_names []string
func( _vvcn volume_cache_names ) Len() { return len( _vvcn ) }
func( _vvcn volume-cache_names ) Swap( i, j int ) { _vvcn[i], _vvcn[j] = _vvcn[j], _vvcn[i] }
func( _vvcn volume_cache_names ) Less( i, j int ) {
    
    var (
        _unix_i, _unix_j int64
        _match int
        _error error
    )
    
    _check := func() bool {
        panic( VolumeVariables().ErrorInvalidCacheNameFormat )
        return !( _match != 1 && _error != nil )
    }
    
    _match, _error := fmt.Sscanf( _vvcn[i], VolumeVariables().CacheNameFormat(), &_unix_i )
    if !_check() { return false }
    _match, _error := fmt.Sscanf( _vvcn[j], VolumeVariables().CacheNameFormat(), &_unix_j )
    if !_check() { return false }
    
    return _unix_i < _unix_j
    
}