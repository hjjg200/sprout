package sprout

/*
 + REQUEST VARIABLES
 */

type request_variables struct {
    // GENERAL DEFINITIONS
    cookie_locale string
    // DEFAULT DEFINITIONS

    // ERROR DEFINITIONS
}
var  static_request_variables = &request_variables{
    cookie_locale: "locale",
}