package network

import (
    "net/http"

    "../environ"
)

type responseWriter struct {
    body        http.ResponseWriter
    req         *Request
    status      int
    wroteHeader bool
}

func newResponseWriter( w http.ResponseWriter ) *responseWriter {
    return &responseWriter{
        body: w,
        status: 0,
        wroteHeader: false,
    }
}

func( rw *responseWriter ) Header() http.Header {
    return rw.body.Header()
}

func( rw *responseWriter ) Write( p []byte ) ( int, error ) {
    if !rw.wroteHeader {
        rw.wroteHeader = true
        rw.status      = 200
    } else {
        // Write is essential to call WriteHeader
        rw.body.WriteHeader( rw.status )
    }
    return rw.body.Write( p )
}

func( rw *responseWriter ) WriteHeader( code int ) {
    if !rw.wroteHeader {
        rw.wroteHeader = true
        rw.status      = code
    } else {
        if code != rw.status {
            environ.Logger.Panicln( ErrDifferentStatusCode )
        }
    }
}