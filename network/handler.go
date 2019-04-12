package network

import (
    "net/http"
)

type Handler func( *Request ) int

func( hnd Handler ) ServeHTTP( w http.ResponseWriter, r *http.Request ) {
    req := NewRequest( w, r )
    hnd( req )
}