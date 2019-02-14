package sprout

/*
 + LOCALE
 *
 * A locale is a set of strings that is capable of localizing text on its own
 */

type Locale struct {
    language string
    locale   map[string] string
}

func( _locale *Locale ) Language() string {
    return _locale.language
}
func( _locale *Locale ) Locale() map[string] string {
    _copy := make( map[string] string )
    for _key, _value := range _locale.locale {
        _copy[_key] = _value
    }
    return _copy
}
func( _locale *Locale ) Localize( _source string ) string {}
func( _locale *Locale ) localize( _source string, _threshold int ) string {}