package sprout

import (
    "net/http"
    "strings"
    "path"
    "os"
)

type HTTPHandler struct {
    parent    *Sprout
    serveHTTP http.HandlerFunc
}

func newHTTPHandler( s *Sprout ) *HTTPHandler {
    return &HTTPHandler{
        parent: s,
        serveHTTP: func ( w http.ResponseWriter, r *http.Request ) {
            // If none of the handlers handled the request
            // This is the fallback
            WriteStatus( w, 404, "Not Found" )
        },
    }
}

func ( hh *HTTPHandler ) ServeHTTP( w http.ResponseWriter, r *http.Request ) {
    hh.serveHTTP( w, r )
}

func ( hh *HTTPHandler ) WithRealtimeAssetServer() *HTTPHandler {

    return hh.withNewHandlerFunc( func ( w http.ResponseWriter, r *http.Request ) {

        url := r.URL.Path
        if isSafeAssetURL( url ) {
            // Remove the first slash at the beginning
            p     := path.Clean( url[1:] )
            b     := path.Base( p )
            ext   := path.Ext( p )

            // Whitelist of Asset Extensions
            //   This is temporary security measure
            //   Liable to being removed or modified
            switch ext {
            case ".css", ".js":
            default:
                // Status Not Found
                WriteStatus( w, 404, "Not Found" )
                return
            }

            f, err := os.Open( p )
            if err != nil {
                // Status Internal Server Error
                WriteStatus( w, 500, "Internal Server Error" )
                return
            }
            st, err := f.Stat()
            if err != nil {
                // Status Internal Server Error
                WriteStatus( w, 500, "Internal Server Error" )
                return
            }

            http.ServeContent( w, r, b, st.ModTime(), f )
            return

        }

        hh.serveHTTP( w, r )

    } )

}

func ( hh *HTTPHandler ) WithCachedAssetServer() *HTTPHandler {

    return hh.withNewHandlerFunc( func ( w http.ResponseWriter, r *http.Request ) {

        url := r.URL.Path
        if isSafeAssetURL( url ) {
            // Remove the first slash at the beginning
            p     := path.Clean( url[1:] )
            b     := path.Base( p )
            ext   := path.Ext( p )

            // Whitelist of Asset Extensions
            //   This is temporary security measure
            //   Liable to being removed or modified
            switch ext {
            case ".css", ".js":
            default:
                // Status Not Found
                WriteStatus( w, 404, "Not Found" )
                return
            }

            a, ok := hh.parent.assets[p]
            if ok {
                // Check if Version Is Set
                v := r.FormValue( "v" )
                if v == "" || v != a.hash[:6] {
                    http.Redirect( w, r, url + "?v=" + a.hash[:6], http.StatusFound )
                    return
                }

                // Serve Content is the Version Is Set
                http.ServeContent( w, r, b, a.modTime, a.reader )
            } else {
                // Status Not Found
                WriteStatus( w, 404, "Not Found" )
            }

            return
        }

        hh.serveHTTP( w, r )

    } )

}

func isSafeAssetURL( url string ) bool {
    if strings.HasPrefix( url, "/asset/" ) {
        if !containsDotDot( url ) {
            return true
        }
    }
    return false
}

// From net/http
func containsDotDot( v string ) bool {
    if !strings.Contains( v, ".." ) {
        return false
    }
    for _, ent := range strings.FieldsFunc( v, isSlashRune ) {
        if ent == ".." {
            return true
        }
    }
    return false
}
func isSlashRune(r rune) bool { return r == '/' || r == '\\' }

func ( hh *HTTPHandler ) WithRoutes() *HTTPHandler {

    return hh.withNewHandlerFunc( func ( w http.ResponseWriter, r *http.Request ) {
        for _, route := range hh.parent.routes {
            rgx := route.rgx
            rhh := route.hh

            if rgx.MatchString( r.URL.Path ) {
                rhh( w, r )
                return
            }
        }
        hh.serveHTTP( w, r )
    } )

}

func ( hh *HTTPHandler ) withNewHandlerFunc( hf http.HandlerFunc ) *HTTPHandler {
    return &HTTPHandler{
        parent: hh.parent,
        serveHTTP: hf,
    }
}