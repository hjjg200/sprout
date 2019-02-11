package sprout

/*
 + ERROR
 *
 * An error is a struct that includes http status code, message, and error details
 */

type Error struct {
    code    int
    message string
    details error
}

func( _error Error ) Error() string {}
func( _error Error ) Copy() Error {}