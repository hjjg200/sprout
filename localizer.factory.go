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
func( _lcrfac *localizer_factory ) FromJSONs( _jsons [][]byte ) ( *Localizer, error ) {

    // _jsons may be empty
    _locales := make( map[string] *Locale )

    // Assign each locale
    for _, _json := range _jsons {
        _locale, _error := LocaleFactory().FromJSON( _json )
        if _error != nil {
            return nil, _error
        }
        _locales[_locale.language] = _locale
    }

    // Set the default locale
    var _default_locale string
    for _, _locale := range _locales {
        _default_locale = _locale.language
        break
    }

    // Assemble
    return &Localizer{
        locales: _locales,
        default_locale: _default_locale,
        threshold: LocalizerVariables().DefaultThreshold(),
        left_delimiter: LocalizerVariables().DefaultLeftDelimiter(),
        right_delimiter: LocalizerVariables().DefaultRightDelimiter(),
    }

}