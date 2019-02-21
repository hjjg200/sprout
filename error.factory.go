package sprout

/*
 + ERROR FACTORY
 */

type error_factory struct{}
var  static_error_factory = &error_factory{}

func ErrorFactory() *error_factory {
    return static_error_factory
}
func( _errfac *error_factory ) New( _code int, _args interface{}... ) Error {

    var _message string
    for i := range _args {
        if i > 0 {
            _message += " "
        }
        _message += fmt.Sprint( _args[i] )
    }

    return Error{
        code: _code,
        details: errors.New( _message ),
    }

}