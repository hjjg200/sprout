package errors

import (
    "fmt"
    "path/filepath"
    "runtime"
)

type Error struct {
    typ string
    children []interface{}
}

/*

[ WARN ]    24.553 - CompileFailure: at volume.Compile:16 - int(1) 12, error(2) someError
    at volume.Do:124 - int(1), int(2)
    at volume.Make:221 - float32(1)
  ProcessingError: at d

*/

func Append( base interface{}, args ...interface{} ) Error {

    var Err Error

    if cast, ok := base.( Error ); ok {
        Err = cast.append( args... )
    } else if cast, ok := base.( error ); ok {
        Err = Error{ typ: "ErrUnknown" }
        args = append( []interface{}{ cast.Error() }, args... )
        Err = Err.append( args... )
    } else {
        Err = Error{}
    }

    // Return
    return Err

}

func caller() string {
    _, file, no, ok := runtime.Caller( 3 )
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

func( Err Error ) append( args ...interface{} ) Error {

    // Newline
    if Err.details != "" {
        Err.details += "\n    "
    }

    Err.details += "at " + caller()

    // Details
    details := ""
    for i, arg := range args {
        typ := fmt.Sprintf( "%T(%d)", arg, i )
        if child, ok := arg.( Error ); ok {
            if i > 0 {
                details += "\n  "
            }
            details += child.Error()
            if Err.children == nil {
                Err.children = make( []Error, 0 )
            }
            Err.children = append( Err.children, child )
            Err.children = append( Err.children, child.children... )
        } else {
            if i > 0 {
                details += " "
            }
            details += typ + ": "
            details += fmt.Sprint( arg )
        }
    }

    Err.details += details

    // Return
    return Err

}
