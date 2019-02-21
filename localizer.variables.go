package sprout

/*
 + LOCALIZER VARIABLES
 */

type localizer_variables struct {
    // GENERAL DEFINITIONS

    // DEFAULT DEFINITIONS
    default_threshold int
    default_left_delimiter string
    default_right_delimiter string

    // ERROR DEFINITIONS
    error_invalid_locale error
    error_invalid_delimiter error
    error_invalid_threshold error
}
var  static_localizer_variables = &localizer_variables{
    // GENERAL DEFINITIONS

    // DEFAULT DEFINITIONS
    default_threshold: 10,
    default_left_delimiter: "{%",
    default_right_delimiter: "%}",

    // ERROR DEFINITIONS
    error_invalid_locale: ErrorFactory().New( 500, "localizer:", "the given locale is invalid" ),
    error_invalid_delimiter: ErrorFactory().New( 500, "localizer:", "the given delimiters are invalid" ),
    error_invalid_threshold: ErrorFactory().New( 500, "localizer:", "the given threshold is invalid" ),
}

func LocaleVariables() *localizer_variables {
    return static_localizer_variables
}
// DEFAULT DEFINITIONS
func( _lcrvar *localizer_variables ) DefaultThreshold() int { return _lcrvar.default_threshold }
func( _lcrvar *localizer_variables ) DefaultLeftDelimiter() string { return _lcrvar.default_left_delimiter }
func( _lcrvar *localizer_variables ) DefaultRightDelimiter() string { return _lcrvar.default_right_delimiter }
// ERROR DEFINITIONS
func( _lcrvar *localizer_variables ) ErrorInvalidLocale() Error { return _lcrvar.error_invalid_locale }
func( _lcrvar *localizer_variables ) ErrorInvalidDelimiter() Error { return _lcrvar.error_invalid_delimiter }
func( _lcrvar *localizer_variables ) ErrorInvalidThreshold() Error { return _lcrvar.error_invalid_threshold }