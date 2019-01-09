package sprout

import (
    "path/filepath"
    "fmt"
    "context"
    "time"
    "os"
    "testing"
    "net/http"
)

var closer chan struct{}

func testCheckError( t *testing.T, err error ) {
    if err != nil {
        t.Error( err )
    }
}

func TestSprout( t *testing.T ) {

    os.Chdir( "D:\\sprout" )

    closer = make( chan struct{}, 1 )

    s := New()
    s.AddRoute( "^/close$", testHandleHTTP3 )
    s.AddRoute( "^/hello_world$", testHandleHTTP2 )
    s.AddRoute( "^/", testHandleHTTP )

    testCheckError( t, s.BuildCache() )
    fmt.Println( s.LoadCache( "asset.zip" ) )

    go func() {
        testCheckError( t, s.StartServer( ":8080" ) )
    }()

    <- closer
    ctx, _ := context.WithTimeout( context.Background(), 5 * time.Second )
    s.srvProduction.Shutdown( ctx )

}

func testHandleHTTP3( w http.ResponseWriter, r *http.Request ) {
    closer <- struct{}{}
}

func testHandleHTTP2( w http.ResponseWriter, r *http.Request ) {
    w.Write( []byte( "hello world" ) )
}

func testHandleHTTP( w http.ResponseWriter, r *http.Request ) {
    p, _ := filepath.Abs( "./" )
    fmt.Println( "PATH:", p )
    w.Write( []byte( p ) )
}