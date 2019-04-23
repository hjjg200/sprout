package errors

import (
    "fmt"
    "testing"
)

func TestErrors01( t *testing.T ) {

    ErrErrorForTest := newType( "ErrErrorForTest" )
    fmt.Println( ErrErrorForTest )

    ErrErrorForTest = ErrErrorForTest.Append( "error 01!", "something went wrong" )
    fmt.Println( ErrErrorForTest )

    ErrErrorForTest = ErrErrorForTest.Append( "error 02!", ErrErrorForTest )
    fmt.Println( ErrErrorForTest )

    ErrErrorForTest = ErrErrorForTest.Append( "error 03!", ErrErrorForTest )
    fmt.Println( ErrErrorForTest )

}