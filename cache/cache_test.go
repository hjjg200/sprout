package cache

import (
    "testing"
    "time"
    "os"
    "io"
)

func Test1( t *testing.T ) {

    chc := NewCache()

    w, _ := chc.Create( "a.txt", time.Now() )
    w.Write( []byte( "AB" ) )

    w.Close()

    r, _ := chc.Open( "a.txt" )
    io.Copy( os.Stdout, r )
    r.Close()

    t.Error( "A" )

}