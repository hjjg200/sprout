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

func newResponseWriter( req *Request, w http.ResponseWriter ) *responseWriter {
    return &responseWriter{
        body: w,
        req: req,
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
        rw.req.logStatus( rw.status )
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
                environ.Logger.OKln( "switched status", rw.req.String(), rw.status, "=>", status )
                rw.status = status
            }
        }
    }
}