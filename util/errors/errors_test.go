package errors

import (
    "fmt"
    "testing"
)

func TestErrors01( t *testing.T ) {

    ErrErrorForTest := newType( "ErrErrorForTest" )
    fmt.Println( ErrErrorForTest )

    ErrErrorForTest = ErrErrorForTest.Raise( "error 01!", "something went wrong" )
    fmt.Println( ErrErrorForTest )

    ErrErrorForTest = ErrErrorForTest.Raise( "error 02!", ErrErrorForTest )
    fmt.Println( ErrErrorForTest )

    ErrErrorForTest = ErrErrorForTest.Raise( "error 03!", ErrErrorForTest )
    fmt.Println( ErrErrorForTest )

}