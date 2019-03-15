package i18n

import (
    "../util"
)

var (
    // GENERAL DEFINITIONS

    // PRIVATE (CONSTANT) DEFINTIONS
    c_defaultThreshold int = 10
    c_defaultLeftDelimiter string = "{%"
    c_defaultRightDelimiter string = "%}"
    c_defaultCookie string = "lang"
    c_defaultQueryParameter string = "l"

    // ERROR DEFINITIONS
    ErrInvalidLocale = util.NewError( 500, "the given locale is invalid" )
    ErrInvalidDelimiters = util.NewError( 500, "the given delimiters are invalid" )
    ErrInvalidThreshold = util.NewError( 500, "the given threshold is invalid" )
    ErrInvalidParameter = util.NewError( 500, "the given parameter is invalid" )
    ErrMalformedJson = util.NewError( 500, "the given JSON is malformed" )
    ErrLocaleNonExistent = util.NewError( 500, "the specified locale does not exist" )
)