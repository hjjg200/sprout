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
func( _locale *Locale ) Localize( _source, _left_delimiter, _right_delimiter string, _threshold int  ) string {

    _follows := func( _2index int, _2substring string ) bool {
        // If the rest string is shorter than the substring
        if _2index + len( _2substring ) > len( _source ) {
            return false
        }
        // Check the first letter
        if _source[_2index] != _2substring[0] {
            return false
        }
        // Check the whole substring
        return _source[_2index:_2index + len( _2substring )] == _2substring
    }

    // Check the Threshold
    if _threshold < 1 {
        _threshold = LocalizerVariables().DefaultThreshold()
    }

    // Do
    for i := 0; i < _threshold; i++ {

        // Check if there is any left delimiter
        if strings.Contains( _source, _left_delimiter ) == false {
            break
        }

        var (
            _key        = ""
            _read       = false
            _found_keys = make( []string, 0 )
        )

        for _index := range _source {
            switch {
            case _follows( _index, _left_delimiter ):
                _key = ""
                _read = true
            case _follows( _index, _right_delimiter ):
                _key = ""
                _read = false
                _found_keys = append( _found_keys, _key )
            default:
                if _read {
                    _key += string( _source[_index] )
                }
            }
        }

        for _, _key := range _found_keys {
            _value, _ok := _locale.locales[_key]
            if _ok {
                _source = strings.Replace(
                    _source,
                    _left_delimiter + _key + _right_delimiter,
                    _value,
                    -1,
                )
            }
        }

    }

    return _source

}