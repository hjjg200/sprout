package util

import (
    "fmt"
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