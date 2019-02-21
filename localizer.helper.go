package sprout

/*
 + LOCALIZER HELPER
 */

type accept_language_entry struct{
    language string
    q_factor float64
}
type accept_language_entries []accept_language_entry

func ( _entries accept_language_entries ) Len() int { return len( _entries ) }
func ( _entries accept_language_entries ) Swap( i, j int ) {
    _entries[i], _entries[j] = _entries[j], _entries[i]
}
func ( _entries accept_language_entries ) Less( i, j int ) bool {
    return _entries[i].q_factor < _entries[j].q_factor
}