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

    lcn, _ := s.BuildCache()
    s.LoadCache( lcn )
    prod, _ := s.Server( "production" )
    prod.Mux().WithRoute( MethodGet, regexp.MustCompile( "^/close$" ), testHandleHTTP3 )
    df := func() interface{} {
        return map[string] map[string] string {
            "data2": map[string] string{
                "1": "ok",
                "2": "nice",
            },
        }
    }
    prod.Mux().WithRoute(
        MethodGet, regexp.MustCompile( "^/(index.html?)?$" ),
        s.ServeCachedTemplate( "template/index.html", df ),
    )
    prod.Mux().WithRoute(
        MethodGet, regexp.MustCompile( "^/hello$" ),
        s.ServeCachedAsset( "asset/hello.html" ),
    )
    prod.Mux().WithHandlerFunc( NotFound )

    go func() {
        testCheckError( t, prod.Start( ":8080" ) )
    }()

    fmt.Println( s.localizer.localize( "{%ui.button%} is good but {%ab%} is bad", "en-us", 3 ) )

    <- closer
    prod.Stop()
    t.Error( "d" )

}

func testHandleHTTP3( req *Request ) bool {
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