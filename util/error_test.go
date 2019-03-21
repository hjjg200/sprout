package util

import (
    "testing"
)

func Test1( t *testing.T ) {

    err := NewError( 500, "Error1" )
    err2 := NewError( 500, "error2" )

    t.Log( err.Has( err2 ) )

    err3 := err.Append( err2, "appended" )

    t.Log( err3.Has( err2 ) )

    t.Log( err3.Is( err ) )

}