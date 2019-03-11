package sprout

/*
 + VOLUME VARIABLES
 */

type volume_variables struct {
    // GENERAL DEFINITIONS
    template_extensions []string
    template_left_delimiter string
    template_right_delimiter string
    cache_name_format string
    max_asset_size int64 // in bytes
    whitelisted_extensions []string

    // CONSTANT DEFINITIONS
    const_asset_directory string
    const_template_directory string
    const_locale_directory string
    const_locale_extension string

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
    cache_name_format: "cache-%d",
    max_asset_size: 10 << 20, // 10MB
    whitelisted_extensions: []string{
        ".css", ".js", ".jpg", ".jpeg", ".png", ".gif", ".ico", ".icn", ".html", ".htm",
    },

    // CONSTANT DEFINITIONS
    const_asset_directory: "asset",
    const_template_directory: "tempplate",
    const_locale_directory: "locale",
    const_locale_extension: ".json",

    // ERROR DEFINITIONS
    error_no_available_cache: Error{
        code: 500,
        details: ErrorFactory().New( "volume:", "there is no available cache" )
    },
}