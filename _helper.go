package sprout

func string_slice_includes( _ss []string, _ch string ) bool {
    for _, _i := range _ss {
        if _ch == _i {
            return true
        }
    }
    return false
}