package sprout

import (
    "fmt"
    "regexp"
    "os"
    "testing"
    "runtime"
    "net/http"
)

var closer chan struct{}

func testCheckError( t *testing.T, err error ) {
    if err != nil {
        t.Error( err )
    }
}

func TestSprout( t *testing.T ) {

    if runtime.GOOS == "windows" {
        os.Chdir( "D:\\sprout" )
    } else if runtime.GOOS == "darwin" {
        os.Chdir( "/Users/anton/sprout" )
    }

    closer = make( chan struct{}, 1 )

    s := New()

    prod, _ := s.Server( "production" )
    prod.Mux().WithRoute( MethodGet, regexp.MustCompile( "^/close$" ), testHandleHTTP3 )
    prod.Mux().WithRoute( MethodGet, regexp.MustCompile( "^/(index.html?)?$" ), testHandleHTTP )
    prod.Mux().WithHandlerFunc( NotFound )

    go func() {
        testCheckError( t, prod.Start( ":8080" ) )
    }()

    fmt.Println( s.localizer.localize( "{%ui.button%} is good but {%ab%} is bad", "en-us", 3 ) )

    <- closer
    prod.Stop()
    t.Error( "d" )

}

func testHandleHTTP3( w http.ResponseWriter, r *http.Request ) bool {
    closer <- struct{}{}
    return true
}
/*
func testHandleHTTP2( w http.ResponseWriter, r *http.Request ) bool {
    ss, err := FetchSession( w, r )
    if err != nil {
        w.Write( []byte( "session created!" ) )
        return true
    }
    w.Write( []byte( "Your sid: " + ss.SID() ) )
    return true
}
*/
func testHandleHTTP( w http.ResponseWriter, r *http.Request ) bool {
    w.Write( []byte( "this is index.html" ) )
    return true
}