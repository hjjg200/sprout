package sprout

/*
 + LOCALIZER VARIABLES
 */

type localizer_variables struct {
    // GENERAL DEFINITIONS
    default_threshold int

    // ERROR DEFINITIONS
    error_invalid_locale error
}
var  static_localizer_variables = &localizer_variables{
    // GENERAL DEFINITIONS
    default_threshold: 10,

    // ERROR DEFINITIONS
    error_invalid_locale: ErrorFactory().New( 500, "localizer:", "the given locale is invalid" ),
}

func LocaleVariables() *localizer_variables {
    return static_localizer_variables
}
func( _locvar *localizer_variables ) DefaultThreshold() int {
    return _locvar.default_threshold
}