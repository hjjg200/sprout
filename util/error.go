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

func( err Error ) Error() string {
    return err.detail
}