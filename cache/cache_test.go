package cache

import (
    "testing"
    "time"
    "os"
    "io"
)

func TestCache01( t *testing.T ) {

    chc := NewCache()

    create := func( name, content string ) {
        w, err := chc.Create( name, time.Now() )
        if err != nil {
            t.Error( err )
        }
        w.Write( []byte( content ) )
        w.Close()
    }

    create( "a.txt", "ABCD" )
    create( "b.txt", "DEFG" )
    chc.Flush()

    r, _ := chc.Files()[0].Open()
    r.Close()
    chc.Flush()

    create( "c.txt", "DDD" )
    create( "e/d.txt", "5345" )
    chc.Flush()

    files := chc.Files()
    for _, f := range files {
        print( f.Name )
        print( " - " )
        r, _ := f.Open()
        io.Copy( os.Stdout, r ); print( "\n" )
        r.Close()
    }
    chc.Flush()

    f, _ := os.OpenFile( "../test/TestCache01/a.zip", os.O_CREATE | os.O_WRONLY | os.O_TRUNC, 0644 )
    f.Write( chc.Data() )
    f.Close()

}