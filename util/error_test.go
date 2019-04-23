package util

import (
    "testing"
)

func TestError01( t *testing.T ) {

    err := NewError( 500, "Error1" )
    err2 := NewError( 500, "error2" )

    t.Log( err.Has( err2 ) )

    err3 := err.Raise( err2, "appended" )

    t.Log( err3.Has( err2 ) )

    t.Log( err3.Is( err ) )

}