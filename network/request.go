package network

import (
    "io"
    "net/http"
    "strings"
    "time"

    "../i18n"
    "../environ"
)

type Request struct {
    body      *http.Request
    closed    bool
    writer    http.ResponseWriter
    localizer *i18n.Localizer
    vars      []string
}

func NewRequest( w http.ResponseWriter, r *http.Request ) *Request {

    // New
    req := &Request{
        body: r,
        writer: w,
    }

    // return
    return req

}

func( req *Request ) Body() *http.Request {
    return req.body
}

func( req *Request ) Header() http.Header {
    return req.writer.Header()
}

func( req *Request ) setStatus( code int ) {
    req.writer.WriteHeader( code )
    req.logStatus( code )
}

func( req *Request ) logStatus( code int ) {

    // Args
    args := []interface{}{
        req.body.Method,
        req.body.Host + req.body.URL.Path,
        req.body.Proto,
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

// Respond

func( req *Request ) Responder( code int ) http.ResponseWriter {
    req.setStatus( code )
    return req.writer
}

func( req *Request ) Respond( code int, content string ) {
    req.RespondContent( code, "text/html;charset=utf-8", content )
}

func( req *Request ) RespondContent( code int, mimeType, content string ) {
    req.setStatus( code )
    req.writer.Header().Set( "content-type", mimeType )
    req.writer.Write( []byte( content ) )
}

func( req *Request ) RespondText( code int, content string ) {
    req.RespondContent( code, "text/plain;charset=utf-8", content )
}

func( req *Request ) RespondJson( code int, obj interface{} ) {

}

func( req *Request ) RespondFile( name string, modTime time.Time, rdskr io.ReadSeeker ) {
    req.logStatus( 200 )
    http.ServeContent( req.writer, req.body, name, modTime, rdskr )
}