package network

import (
    "fmt"
    "net/http"
    "strings"

    "../i18n"
    "../environ"
)

type Request struct {
    body      *http.Request
    writer    *responseWriter
    localizer *i18n.Localizer
    space     *Space
    rsp       *Responder
    vars      []string
}

func NewRequest( w http.ResponseWriter, r *http.Request ) *Request {

    // New
    req := &Request{
        body: r,
    }
    
    req.writer = newResponseWriter( req, w )

    // return
    return req

}

func( req *Request ) Body() *http.Request {
    return req.body
}

func( req *Request ) Writer() http.ResponseWriter {
    return req.writer
}

func( req *Request ) Header() http.Header {
    return req.writer.Header()
}

func( req *Request ) String() string {
    return fmt.Sprintf(
        "%s %s %s <= %s",
        req.body.Method,
        req.body.Host + req.body.URL.Path,
        req.body.Proto,
        req.body.RemoteAddr,
    )
}

func( req *Request ) logStatus( code int ) {

    // Args
    args := []interface{}{
        req.String(),
        code,
    }

    // Log
    switch {
    case code >= 500:
        environ.Logger.Warnln( args... )
    default:
        environ.Logger.OKln( args... )
    }

}

func( req *Request ) Responder( code int ) *Responder {

    var rsp *Responder

    if req.rsp != nil {
        rsp = req.rsp
    } else {
        rsp = &Responder{
            req: req,
            writer: req.writer,
        }
        req.rsp = rsp
    }

    rsp.SetStatus( code )
    return rsp

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
