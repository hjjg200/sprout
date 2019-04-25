package network

import (
    "net/http"
    "strings"
    "time"
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

// From net/http
func checkIfModifiedSince( r *http.Request, modtime time.Time ) ( yes, ok bool ) {
    if r.Method != "GET" && r.Method != "HEAD" {
        return false, false
    }
    ims      := r.Header.Get( "If-Modified-Since" )
    zeroTime := time.Time{}
    if ims == "" || zeroTime == modtime {
        return false, false
    }
    t, err := http.ParseTime( ims )
    if err != nil {
        return false, false
    }
    // The Date-Modified header truncates sub-second precision, so
    // use mtime < t+1s instead of mtime <= t to check for unmodified.
    if modtime.Before( t.Add( 1 * time.Second ) ) {
        return false, true
    }
    return true, true
}