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
func( _locfac *locale_factory ) FromJSON( _json []byte ) ( *Locale, error ) {

    // Unmarshal JSON Data
    var _interface interface{}
    _err := json.Unmarshal( _json, &_interface )
    if _err != nil {
        return nil, _err
    }

    // Check the Language
    _value := reflect.ValueOf( _interface )
        // There must be one key under the root
    if len( _value.MapKeys() ) > 1 {
        return nil, LocaleVariables().ErrorTooManyJSONKeys()
    }
        // Language is the key name of the key
    _language := _value.MapKeys()[0].String()
        // And the key must be a map
    _map := _value.MapIndex( _value.MapKeys()[0] )
    switch _cast := _map.Interface().( type ) {
    default:
        _reflect_cast := reflect.ValueOf( _cast )
        if _reflect_cast.Kind() != reflect.Map {
            return nil, LocaleVariables().ErrorMalformedJSON()
        }
    }

    // Parse the first map
    var (
        _locale = make( map[string] string )
        _recursive_do func( string, reflect.Value ) // Need to be pre-declared to be used as a recursive function
    )
    _recursive_do = func( _2base_key string, _2map reflect.Value ) {
        for _, _2map_key := range _2map.MapKeys() {
            _2element := _2map.MapIndex( _2map_key )
            // Make the next base key
            var _2current_key
            if _2base_key != "" {
                _2current_key = _2base_key + "." + _2map_key.String()
            } else {
                _2current_key = _2base_key.String()
            }
            // Assign value
            switch _2cast := _2element.Interface().( type ) {
            case string:
                _locale[_2current_key] = _2cast
            default:
                _2reflect_cast := reflect.ValueOf( _2cast )
                if _2reflect_cast.Kind() == reflect.Map {
                    _recursive_do( _2current_key, _2reflect_cast )
                }
            }
        }
    }
    // Call the recursive function with an empty key
    _recursive_do( "", _map )

    return &Locale{
        Language: strings.ToLower( _language ),
        locale: _locale,
    }, nil

}