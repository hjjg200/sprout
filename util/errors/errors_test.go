package errors

import (
    "fmt"
    "testing"
)

func TestErrors01( t *testing.T ) {

    ErrIO := newType( "ErrIO" )
    ErrFunc := newType( "ErrFunc" )
    ErrFormat := newType( "ErrFormat" )
    ErrValue := newType( "ErrValue" )

    err := Append( ErrValue, 22 ); fmt.Println( err )
    fmt.Println( "---" )
    err = Append( err, 30, 35 ); fmt.Println( err )
    fmt.Println( "---" )
    err = Append( ErrFormat, "given foramt blah", err ); fmt.Println( err )
    fmt.Println( "---" )
    err = Append( err, "YYYY-mm-dd" ); fmt.Println( err )
    fmt.Println( "---" )
    err = Append( ErrFunc, "error here", err ); fmt.Println( err )
    fmt.Println( "---" )
    err = Append( ErrIO, "failed operation", err ); fmt.Println( err )



}