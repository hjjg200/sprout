package util

import (
    "fmt"
    "strings"
)

type Error struct {
    id int
    code int
    children []Error
    detail string
}

var errorIncrement = -1

func NewError( code int, args... interface{} ) Error {
    var msg string
    for i := range args {
        if i > 0 {
            msg += " "
        }
        msg += fmt.Sprint( args[i] )
    }

    err := Error{
        code: code,
        detail: msg,
    }
    err.renewID()

    return err
}

func ErrorHasPrefix( err error, Err2 Error ) bool {
    if Err1, ok := err.( Error ); ok {
        return strings.HasPrefix( Err1.detail, Err2.detail )
    }
    return false
}

func( err *Error ) renewID() {
    errorIncrement++
    err.id = errorIncrement
}

func( err1 Error ) Has( err2 Error ) bool {
    if err1.id == err2.id {
        return true
    }
    for i := range err1.children {
        if err1.children[i].id == err2.id {
            return true
        }
    }
    return false
}

func( err1 Error ) Is( err2 Error ) bool {
    return err1.id == err2.id
}

func( err Error ) Append( args... interface{} ) Error {
    err.detail += "\n  "
    for i := range args {
        err.detail += " " + fmt.Sprint( args[i] )
        switch cast := args[i].( type ) {
        case Error:
            err.children = append( err.children, cast )
        }
    }
    return err
}

func( err Error ) String() string {
    return fmt.Sprintf( "%d %s: %s", err.code, HttpStatusMessages[err.code], err.detail )
}

func( err Error ) Error() string {
    return err.String()
}