package sprout

/*
 + ERROR
 *
 * An error is a struct that includes http status code, message, and error details
 */

type Error struct {
    code    int
    details error
}

func( _error Error ) Code() int {
    return _error.code
}
func( _error Error ) Details() error {
    return _error.details
}
func( _error Error ) Error() string {
    return fmt.Sprintf(
        "%d %s - %s",
        _error.code,
        SproutVariables().HTTPStatusMessages()( _error.code ),
        _error.details.Error(),
    )
}