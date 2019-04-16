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
    finalized   bool
}

func newResponseWriter( w http.ResponseWriter ) *responseWriter {
    return &responseWriter{
        body: w,
        status: 0,
        wroteHeader: false,
        finalized: false,
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
        rw.finalized = true
        rw.body.WriteHeader( rw.status )
    }
    return rw.body.Write( p )
}

func( rw *responseWriter ) WriteHeader( status int ) {
    if !rw.wroteHeader {
        rw.wroteHeader = true
        rw.status      = status
    } else {
        if status != rw.status {
            if rw.finalized {
                environ.Logger.Panicln( ErrDifferentStatusCode )
            } else {
                rw.status = status
            }
        }
    }
}