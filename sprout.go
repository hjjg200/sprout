package sprout

import (
    "bytes"
    "encoding/json"
    "log"
    "io"
    "net/http"
    "regexp"
    "time"
)

var (
    envOS string
)

type route struct {
    rgx *regexp.Regexp
    hh  http.HandlerFunc
}

type asset struct {
    modTime time.Time
    reader  *bytes.Reader

    // sha256 hash of the file
    hash    string
}

type Sprout struct {
    assets        map[string] asset
    routes        []route
    srvProduction *http.Server
    srvDev        *http.Server
}

func New() *Sprout {

    s := &Sprout{}

    s.assets = make( map[string] asset )
    s.routes = make( []route, 0 )

    err := sanityCheck()
    if err != nil {
        log.Fatalln( err )
    }

    return s

}

func sanityCheck() error {
    // check if there is any sass, scss if so check sass installed
    if err := checkOS(); err != nil {
        return err
    }
    return nil
}

func WriteStatus( w http.ResponseWriter, code int, msg string ) {
    w.Header().Set( "Content-Type", "text/html" )


}

func WriteJSON( w io.Writer, v interface{} ) error {
    return json.NewEncoder( w ).Encode( v )
}

func ( s *Sprout ) AddRoute( rgxStr string, hh http.HandlerFunc ) error {

    rgx, err := regexp.Compile( rgxStr )
    if err != nil {
        return err
    }

    s.routes = append( s.routes, route{
        rgx: rgx,
        hh: hh,
    } )

    return nil

}

func ( s *Sprout ) StartServer( addr string ) error {
    hh := newHTTPHandler( s )
    hh  = hh.WithRoutes()
    hh  = hh.WithCachedAssetServer()

    s.srvProduction = &http.Server{
        Addr: addr,
        Handler: hh,
    }
    return s.srvProduction.ListenAndServe()
}

func ( s *Sprout ) StartDevServer( addr string ) error {
    hh := newHTTPHandler( s )
    hh  = hh.WithRoutes()
    hh  = hh.WithRealtimeAssetServer()

    s.srvDev = &http.Server{
        Addr: addr,
        Handler: hh,
    }
    return s.srvDev.ListenAndServe()
}