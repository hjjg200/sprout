package sprout

/*
 + LOCALE VARIABLES
 */

type locale_variables struct {
    // GENERAL DEFINITIONS
    default_threshold int

    // ERROR DEFINITIONS
    error_invalid_locale error
}
var  static_locale_variables = &locale_variables{
    // GENERAL DEFINITIONS
    default_threshold: 10,

    // ERROR DEFINITIONS
    error_invalid_locale: Error{
        code: 500,
        details: ErrorFactory().New( "locale:", "the given locale is invalid" )
    },
}

func LocaleVariables() *locale_variables {
    return static_locale_variables
}
func( _locvar *locale_variables ) DefaultThreshold() int {
    return _locvar.default_threshold
}