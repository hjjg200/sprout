package sprout

import (
    "bytes"
    "crypto/sha256"
    "encoding/json"
    "errors"
    "fmt"
    "io"
    "log"
    "net/http"
    "regexp"
    "strconv"
    "time"
)

/*
 + Private Variables
 */

const (
    envAppName = "sprout"
    envVersion = "pre-alpha 0.1"

    // Directory names must not contain slashes, dots, etc.
    envDirAsset = "asset"
    envDirCache = "cache"
)

var (
    envOS     string
    envLogger *log.Logger
)

/*
 + Public Variables
 */

var (
    ErrNotSupportedOS = errors.New( "sprout: the OS is not supported" )
)

var (
    EnvFilenameTimeFormat = "20060102-150405"
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

func makeAsset( mt time.Time, r io.Reader ) asset {
    h   := sha256.New()
    mts := strconv.FormatInt( mt.Unix(), 10 )
    h.Write( []byte( mts ) )
    buf := bytes.NewBuffer( nil )
    io.Copy( buf, r )
    return asset{
        modTime: mt,
        reader: bytes.NewReader( buf.Bytes() ),
        hash: fmt.Sprintf( "%x", h.Sum( nil ) ),
    }
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
    c := fmt.Sprint( code )
    t := `<!doctype html>
<html>
    <head>
        <title>` + c + " " + msg + `</title>
        <style>
            html {
                font-family: sans-serif;
                padding: 0;
            }
            body {
                color: hsl( 220, 5%, 45% );
                text-align: center;
                padding: 10px;
                margin: 0;
            }
            div {
                border: 1px dashed hsl( 220, 5%, 88% );
                padding: 20px;
                margin: 0 auto;
                max-width: 300px;
                text-align: left;
            }
            h1, h2, h3 {
                display: block;
                margin: 0 0 5px 0;
            }
            footer {
                color: hsl( 220, 5%, 68% );
                font-family: monospace;
                font-size: 1em;
                text-align: right;
                line-height: 1.3;
            }
        </style>
    </head>
    <body>
        <div>
            <h1>` + c + `</h1>
            <h3>` + msg + `</h3>
            <footer>` + envAppName + " " + envVersion + `<br />on ` + envOS + `</footer>
        </div>
    </body>
</html>`
    w.Write( []byte( t ) )
}

func WriteJSON( w io.Writer, v interface{} ) error {
    return json.NewEncoder( w ).Encode( v )
}

func ( s *Sprout ) AddRoute( rgxStr string, hh http.HandlerFunc ) error {

    rgx, err := regexp.Compile( rgxStr )
    if err != nil {
        panic( err )
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

    // Load the Latest Cache
    // Build Cache If None Found

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