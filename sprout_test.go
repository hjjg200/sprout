package sprout

import (
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
    s.AddRoute( "^/session$", testHandleHTTP2 )
    s.AddRoute( "^/(index.html?)?$", testHandleHTTP )

    go func() {
        testCheckError( t, s.StartServer( ":8080" ) )
    }()
    go func() {
        testCheckError( t, s.StartDevServer( ":8081" ) )
    }()

    <- closer
    ctx, _ := context.WithTimeout( context.Background(), 5 * time.Second )
    s.srvProduction.Shutdown( ctx )
    s.srvDev.Shutdown( ctx )
    t.Error( "d" )

}

func testHandleHTTP3( w http.ResponseWriter, r *http.Request ) bool {
    closer <- struct{}{}
    return true
}

func testHandleHTTP2( w http.ResponseWriter, r *http.Request ) bool {
    ss, err := FetchSession( w, r )
    if err != nil {
        w.Write( []byte( "session created!" ) )
        return true
    }
    w.Write( []byte( "Your sid: " + ss.SID() ) )
    return true
}

func testHandleHTTP( w http.ResponseWriter, r *http.Request ) bool {
    w.Write( []byte( "this is index.html" ) )
    return true
}