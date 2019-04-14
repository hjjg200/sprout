package network

import (
    "fmt"
    "net/http"
    "strings"
    "runtime"

    "../i18n"
    "../util"
    "../environ"
)

type Request struct {
    body      *http.Request
    closed    bool
    writer    http.ResponseWriter
    header    Header
    localizer *i18n.Localizer
    vars      []string
}

func NewRequest( w http.ResponseWriter, r *http.Request ) *Request {

    // New
    req := &Request{
        body: r,
        writer: w,
    }
    
    // Header
    req.header = Header{
        body: w.Header(),
        req: req,
    }

    // return
    return req

}

func( req *Request ) Body() *http.Request {
    return req.body
}

func( req *Request ) ensureOpen() {
    if req.closed {
        environ.Logger.Panicln( ErrRequestClosed )
    }
}

func( req *Request ) Write( p []byte ) ( int, error ) {
    req.ensureOpen()
    return req.writer.Write( p )
}

func( req *Request ) Header() Header {
    req.ensureOpen()
    return req.header
}

func( req *Request ) Closed() bool {
    return req.closed
}

func( req *Request ) Close( status int ) error {
    
    req.ensureOpen()
    req.closed = true
    req.writer.WriteHeader( status )
    
    // Args
    args := []interface{}{
        req.body.Method,
        req.body.Host + req.body.URL.Path,
        req.body.Proto,
        status,
    }

    // Log
    switch {
    case status >= 500:
        environ.Logger.Warnln( args... )
    default:
        environ.Logger.OKln( args... )
    }
    
    return nil
    
}

func( req *Request ) Localizer() *i18n.Localizer {
    return req.localizer
}

func( req *Request ) Vars() []string {
    return req.vars
}

func( req *Request ) PopulateLocalizer( i1 *i18n.I18n ) {

    // Check locale
    switch i1.NumLocale() {
    case 0:
    default:
        lcName, err := i1.ParseUrlPath( req.body.URL )
        if err != nil {
            lcName, err = i1.ParseUrlQuery( req.body.URL )
            if err != nil {
                lcName, err = i1.ParseCookies( req.body.Cookies() )
                if err != nil {
                    lcName, err = i1.ParseAcceptLanguage( req.body.Header.Get( "accept-language" ) )
                    if err != nil {
                        lcName = i1.DefaultLocale()
                    }
                }
            }
        } else { // When the URL path contains the locale at its beginning
            // Remove the locale from the URL for further processing
            path  := req.body.URL.Path[1:]
            split := strings.SplitN( path, "/", 2 )

            // Set path
            req.body.URL.Path = "/"
            if len( split ) == 2 {
                req.body.URL.Path += split[1]
            }
        }
        req.localizer = i1.Localizer( lcName )
        return

    }

    req.localizer = nil

}

// Others

func( req *Request ) WriteStatus( code int ) {

    // Content
    c   := fmt.Sprint( code )
    msg := util.HttpStatusMessages[code]
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
            <footer>` + environ.AppName + " " + environ.AppVersion + `<br />on ` + runtime.GOOS + `</footer>
        </div>
    </body>
</html>`

    req.Write( []byte( t ) )

}