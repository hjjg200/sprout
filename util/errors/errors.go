package errors

import (
    "fmt"
    "path/filepath"
    "runtime"
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

func( Err Error ) Has( other Error ) bool {
    for _, child := range Err.children {
        if child.typ == other.typ {
            return true
        }
    }
    return false
}

func( Err Error ) Is( other Error ) bool {
    return Err.typ == other.typ
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