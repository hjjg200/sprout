package sprout

/*
 + VOLUME VARIABLES
 */

type volume_variables struct {
    // GENERAL DEFINITIONS
    template_extensions []string
    template_left_delimiter string
    template_right_delimiter string

    // DEFAULT DEFINITIONS
    default_whitelisted_extensions []string

    // ERROR DEFINITIONS
    error_no_available_cache error
}
var  static_volume_variables = &volume_variables{
    // GENERAL DEFINITIONS
    template_extensions: []string{
        ".html", ".htm",
    },
    template_left_delimiter: "{{",
    template_right_delimiter: "}}",

    // DEFAULT DEFINITIONS
    default_whitelisted_extensions: []string{
        ".css", ".js", ".jpg", ".jpeg", ".png", ".gif", ".ico", ".icn", ".html", ".htm",
    },

    // ERROR DEFINITIONS
    error_no_available_cache: Error{
        code: 500,
        details: ErrorFactory().New( "volume:", "there is no available cache" )
    },
}