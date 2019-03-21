package network

import (
    "net/http"
    "strings"

    "../i18n"
)

type Request struct {
    body      *http.Request
    writer    http.ResponseWriter
    localizer *i18n.Localizer
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

func( req *Request ) Writer() http.ResponseWriter {
    return req.writer
}

func( req *Request ) Localizer() *i18n.Localizer {
    return req.localizer
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

        req.localizer, _ = i1.Localizer( lcName )
        return

    }

    req.localizer = nil

}

// Others

func( req *Request ) WriteStatus( code int ) {

}