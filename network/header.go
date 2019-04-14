package network

import (
    "io"
    "net/http"
)

type Header struct {
    body http.Header
    req *Request
}

func( h Header ) Add( key, value string ) {
    h.req.ensureOpen()
    h.body.Add( key, value )
}

func( h Header ) Del( key string ) {
    h.req.ensureOpen()
    h.body.Del( key )
}

func( h Header ) Get( key string ) string {
    h.req.ensureOpen()
    return h.body.Get( key )
}

func( h Header ) Set( key, value string ) {
    h.req.ensureOpen()
    h.body.Set( key, value )
}

func( h Header ) Write( w io.Writer ) error {
    h.req.ensureOpen()
    return h.body.Write( w )
}

func( h Header ) WriteSubset( w io.Writer, exclude map[string] bool ) error {
    h.req.ensureOpen()
    return h.body.WriteSubset( w, exclude )
}