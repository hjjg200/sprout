package network

import (
    "net/http"
    "strings"
)

func MakeMethodChecker( ss []string ) map[string] bool {

    fl := make( map[string] bool )
    fl[http.MethodGet]     = false
    fl[http.MethodHead]    = false
    fl[http.MethodPost]    = false
    fl[http.MethodPut]     = false
    fl[http.MethodPatch]   = false
    fl[http.MethodDelete]  = false
    fl[http.MethodConnect] = false
    fl[http.MethodOptions] = false
    fl[http.MethodTrace]   = false

    for _, s := range ss {
        fl[strings.ToUpper( s )] = true
    }

    return fl

}