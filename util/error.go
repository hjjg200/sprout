package util

import (
    "fmt"
    "strings"
)

type Error struct {
    code int
    detail string
}

func NewError( code int, args... interface{} ) Error {
    var msg string
    for i := range args {
        if i > 0 {
            msg += " "
        }
        msg += fmt.Sprint( args[i] )
    }

    return Error{
        code: code,
        detail: msg,
    }
}

func ErrorHasPrefix( err error, Err2 Error ) bool {
    if Err1, ok := err.( Error ); ok {
        return strings.HasPrefix( Err1.detail, Err2.detail )
    }
    return false
}

func( err Error ) Append( args... interface{} ) Error {
    for i := range args {
        err.detail += " " + fmt.Sprint( args[i] )
    }
    return err
}

func( err Error ) String() string {
    return fmt.Sprintf( "%d %s: %s", err.code, HttpStatusMessages[err.code], err.detail )
}

func( err Error ) Error() string {
    return err.String()
}