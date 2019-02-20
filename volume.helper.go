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