package sprout

import (
    "path/filepath"
    "fmt"
    "context"
    "time"
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

    closer = make( chan struct{}, 1 )

    s := New()
    err := s.AddRoute( "^/close$", testHandleHTTP3 )
    testCheckError( t, err )
    err = s.AddRoute( "^/hello_world$", testHandleHTTP2 )
    testCheckError( t, err )
    err = s.AddRoute( "^/", testHandleHTTP )
    testCheckError( t, err )

    testCheckError( t, s.BuildCache() )

    go func() {
        testCheckError( t, s.StartDevServer( ":8080" ) )
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