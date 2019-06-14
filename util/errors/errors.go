package errors

import (
    "fmt"
    "path/filepath"
    "runtime"
)

type errorString struct {
    s string
}

func( e errorString ) Error() string {
    return e.s
}

func caller( skip int ) string {
    _, file, no, ok := runtime.Caller( skip )
    if !ok {
        return ""
    }
    dir  := filepath.Base( filepath.Dir( file ) )
    file  = dir + "/" + filepath.Base( file )
    return fmt.Sprintf( "%s:%d", file, no )
}

func Stack( err error ) error {
    clr := caller( 2 )
    return errorString{
        s: err.Error() + "\n  at " + clr,
    }
}