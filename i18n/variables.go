package i18n


var (
    // GENERAL DEFINITIONS

    // PRIVATE (CONSTANT) DEFINTIONS
    c_defaultThreshold int = 10
    c_defaultLeftDelimiter string = "{%"
    c_defaultRightDelimiter string = "%}"
    c_defaultCookie string = "lang"
    c_defaultQueryParameter string = "l"

    // ERROR DEFINITIONS
    ErrorInvalidLocale = util.NewError( 500, "the given locale is invalid" )
    ErrorInvalidDelimiter = util.NewError( 500, "the given delimiters are invalid" )
    ErrorInvalidThreshold = util.NewError( 500, "the given threshold is invalid" )
)