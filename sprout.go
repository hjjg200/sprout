package sprout

import (
    "bytes"
    "crypto/sha256"
    "encoding/json"
    "errors"
    "fmt"
    "io"
    "net/http"
    "os"
    "strconv"
    "time"

    "./log"
//  "./session"
)

/*
 + Private Variables
 */

const (
    envAppName = "sprout"
    envVersion = "pre-alpha 0.3"

    // Directory names must not contain slashes, dots, etc.
    envDirAsset = "asset"
    envDirCache = "cache"
)

var (
    envOS string
)

/*
 + Public Variables
 */

var (
    ErrNotSupportedOS = errors.New( "sprout: the OS is not supported" )
    ErrDirectory      = errors.New( "sprout: could not access a necessary directory" )
    ErrInvalidDirPath = errors.New( "sprout: the given path is invalid" )
)

var (
    EnvFilenameTimeFormat = "20060102-150405"
)

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
    assets  map[string] asset
    servers map[string] *Server
}

func New() *Sprout {

    s := &Sprout{}

    log.Infoln( "Preparing a new Sprout instance..." )

    s.assets  = make( map[string] asset )
    s.servers = make( map[string] *Server )

    prod, _ := s.NewServer( "production" )
    prod.Mux().WithCachedAssetServer()

    log.Infoln( "Loacting the latest cache..." )
    // Load the Latest Cache
    // Build Cache If None Found
    lcn, err := s.LatestCacheName()
    if err != nil {
        log.Infoln( "Could not load the latest cache, attempting to build one..." )
        lcn, err = s.BuildCache()
        if err != nil { log.Severeln( err ) }
        log.Infoln( "Successfully built a cache:", lcn )
    } else {
        err = s.LoadCache( lcn )
        if err != nil { log.Severeln( err ) }
        log.Infoln( "Loaded Cache:", lcn )
    }

    dev, _  := s.NewServer( "dev" )
    dev.Mux().WithRealtimeAssetServer()

    sanityCheck()
    return s

}

func sanityCheck() error {
    // check if there is any sass, scss if so check sass installed
    log.Infoln( "Doing a sanity check..." )
    if err := checkOS(); err != nil {
        log.Severeln( "Sanity check failed!" )
    }
    if err := ensureDirectories(); err != nil {
        log.Severeln( "Sanity check failed!" )
    }
    log.Infoln( "Everything looks fine!" )
    return nil
}

func ensureDirectories() error {
    log.Infoln( "Ensuring all the necessary directories..." )
    err := ensureDirectory( envDirAsset )
    if err != nil {
        log.Warnln( "Could not ensure all the directories!" )
        return err
    }
    err = ensureDirectory( envDirCache )
    if err != nil {
        log.Warnln( "Could not ensure all the directories!" )
        return err
    }
    log.Infoln( "Ensured all the directories!" )
    return nil
}

func ensureDirectory( p string ) error {
    log.Infoln( "Ensuring a directory...", p )
    st, err := os.Stat( p )
    switch {
    case os.IsNotExist( err ):
        err = os.Mkdir( p, 0750 )
        if err != nil {
            log.Warnln( "Error during ensuring the directory:", p, err )
            return err
        }
    case err != nil:
        log.Warnln( "Error during ensuring the directory:", p, err )
        return err
    case !st.IsDir():
        log.Warnln( "Error during ensuring the directory:", p, ErrDirectory )
        return ErrDirectory
    }
    log.Infoln( "Directory ready to go", p )
    return nil
}
/*
func FetchSession( w http.ResponseWriter, r *http.Request ) ( *session.Session, error ) {
    return session.Fetch( w, r )
}
*/
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
                line-height: 1.0;
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
/*
func ( s *Sprout ) AddRoute( rgxStr string, hh HTTPHandlerFunc ) error {

    rgx, err := regexp.Compile( rgxStr )
    if err != nil {
        panic( err )
        return err
    }

    s.routes = append( s.routes, route{
        rgx: rgx,
        hh: hh,
    } )

    log.Infoln( "Added a route:", rgxStr )
    return nil

}*/

/*
func ( s *Sprout ) StartServer( addr string ) error {

    hh := newHTTPHandler( s )
    hh  = hh.WithRoutes()
    hh  = hh.WithCachedAssetServer()

    log.Infoln( "Starting the production server..." )

    // Load the Latest Cache
    // Build Cache If None Found
    lcn, err := s.LatestCacheName()
    if err != nil {
        log.Infoln( "Could not load the latest cache, attempting to build one..." )
        lcn, err = s.BuildCache()
        if err != nil { return err }
        log.Infoln( "Successfully built a cache:", lcn )
    } else {
        err = s.LoadCache( lcn )
        if err != nil { return err }
        log.Infoln( "Loaded Cache:", lcn )
    }

    s.srvProduction = &http.Server{
        Addr: addr,
        Handler: hh,
    }

    log.Infoln( "The production server listens:", addr )
    return s.srvProduction.ListenAndServe()

}

func ( s *Sprout ) StartDevServer( addr string ) error {

    hh := newHTTPHandler( s )
    hh  = hh.WithRoutes()
    hh  = hh.WithRealtimeAssetServer()

    log.Infoln( "Starting the dev server..." )

    s.srvDev = &http.Server{
        Addr: addr,
        Handler: hh,
    }

    log.Infoln( "The dev server listens:", addr )
    return s.srvDev.ListenAndServe()

}
*/