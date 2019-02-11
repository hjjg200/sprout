package sprout

/*
 + VOLUME MANAGER
 *
 * A volume manager is a pseudo-static object used for managing Volume
 */

type volume_factory struct {}
var  static_volume_factory = &volume_factory{}

func VolumeFactory() *volume_factory {}
func( _volfac *volume_factory ) FromSource( _path string ) ( *Volume, error ) {}
func( _volfac *volume_factory ) New() ( *Volume ) {}