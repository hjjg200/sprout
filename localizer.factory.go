package sprout

/*
 + LOCALIZER FACTORY
 */

type localizer_factory struct{}
var  static_localizer_factory = &localizer_factory

func LocalizerFactory() *localizer_factory {
    return static_localizer_factory
}
func( _lcrfac *localizer_factory ) FromDiectory( _path string ) ( *Localizer, error ) {
    
}