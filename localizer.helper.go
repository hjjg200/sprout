package sprout

/*
 + LOCALIZER HELPER
 */

type localizer_helper struct{}
var  static_localizer_helper = &localizer_helper{}

func LocalizerHelper() *localizer_helper {
    return static_localizer_helper
}
func( _lcrhlpr *localizer_helper ) IsValidLocaleName( _locale string ) bool {

   // https://tools.ietf.org/html/rfc3066#page-2

    var (
        _alphanum = false
        _primary_len = 0
        _sub_len = 0
    )
    for i := range _locale {
        switch {
        case _locale[i] >= 'A' && _locale[i] <= 'Z',
             _locale[i] >= 'a' && _locale[i] <= 'z',
             _locale[i] >= '0' && _locale[i] <= '9':
            if !_alphanum {
                if _locale[i] >= '0' && _locale[i] <= '9' {
                    return false
                }
                _sub_len++
            } else {
                _primary_len++
            }
        case _locale[i] == '-':
            if _alphanum { return false }
            _alphanum = true
        default:
            return false
        }
    }
    if _primary_len > 8 || _sub_len > 8 {
        return false
    }
    return true
}

type accept_language_entry struct{
    language string
    q_factor float64
}
type accept_language_entries []accept_language_entry

func( _entries accept_language_entries ) Len() int { return len( _entries ) }
func( _entries accept_language_entries ) Swap( i, j int ) {
    _entries[i], _entries[j] = _entries[j], _entries[i]
}
func( _entries accept_language_entries ) Less( i, j int ) bool {
    return _entries[i].q_factor < _entries[j].q_factor
}