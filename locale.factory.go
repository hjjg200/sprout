package sprout

/*
 + LOCALE FACTORY
 *
 * A locale manager is a pseudo-static class that handles locale-related things.
 */

type locale_factory struct {}
var  static_locale_factory = &locale_factory{}

func LocaleFactory() *locale_factory {
    return static_locale_factory
}
func( _locfac *locale_factory ) FromJSON( _json string ) ( *Locale, error ) {}
func( _locfac *locale_factory ) from( _source interface{} ) ( *Locale, error ) {}