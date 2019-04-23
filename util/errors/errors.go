package errors

import (
    "fmt"
    "path/filepath"
    "runtime"

    "github.com/hjjg200/sprout/environ"
)

type Error struct {
    typ string
    details string
    children []Error
}

func New( base Error, args ...interface{} ) Error {

    // Details
    Err := Error{ typ: base.typ }
    Err.Append( args... )

    // Return
    return Err

}

func caller() string {
    _, file, no, ok := runtime.Caller( 2 )
    if !ok {
        return ""
    }
    dir  := filepath.Base( filepath.Dir( file ) )
    file  = dir + "/" + filepath.Base( file )
    return fmt.Sprintf( "%s:%d - ", file, no )
}

func newType( typ string ) Error {
    return Error{ typ: typ }
}

func( Err Error ) Error() string {
    return fmt.Sprint( Err.typ + ": " + Err.details )
}

func( Err Error ) Has( other interface{} ) bool {
    Other, ok := other.( Error )
    if !ok {
        return false
    }
    for _, child := range Err.children {
        if child.typ == Other.typ {
            return true
        }
    }
    return false
}

func( Err Error ) Is( other interface{} ) bool {
    Other, ok := other.( Error )
    if !ok {
        return false
    }
    return Err.typ == Other.typ
}

func( Err Error ) Raise() {
    environ.Logger.Warnln( Err )
}

func( Err Error ) Append( args ...interface{} ) Error {

    // Newline
    if Err.details != "" {
        Err.details += "\n  "
    }

    Err.details += caller()

    // Details
    details := ""
    for i, arg := range args {
        if child, ok := arg.( Error ); ok {
            if i > 0 {
                details += "\n  "
            }
            details += child.Error()
            if Err.children == nil {
                Err.children = make( []Error, 0 )
            }
            Err.children = append( Err.children, child )
        } else {
            if i > 0 {
                details += " "
            }
            details += fmt.Sprint( arg )
        }
    }

    Err.details += details

    // Return
    return Err

}
