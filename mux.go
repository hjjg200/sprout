package sprout

import (
    "net/http"
    "regexp"
    "strings"
    "path"
    "os"
)

// Authenticator returns true if the given request contains suitable info to be authenticated
//   false otherwise
type Authenticator func( w http.ResponseWriter, r *http.Request ) bool

// HandlerFunc returns true when it handled the request and no other following handlers are needed
//   returns false when it could not handle the request
type HandlerFunc func( w http.ResponseWriter, r *http.Request ) bool
type Mux struct {
    parent  *Sprout
    handler HandlerFunc
}

const (
    MethodGet = 1 << iota
    MethodHead
    MethodPost
    MethodPut
    MethodPatch
    MethodDelete
    MethodConnect
    MethodOptions
    MethodTrace
)

func ( s *Sprout ) NewMux() *Mux {
    return &Mux{
        parent: s,
        handler: func ( w http.ResponseWriter, r *http.Request ) bool {
            return false
        },
    }
}

// interface http.Handler
func ( m *Mux ) ServeHTTP( w http.ResponseWriter, r *http.Request ) {

    // Check Locale
    _lcc        := len( m.parent.localizer.locales )
    _get_locale := func () string {
        __cookie, __err := r.Cookie( cookie_locale )
        if __err != nil {
            return ""
        }
        return __cookie.Value
    }
    _set_locale := func( _lc string ) {
        __cookie, __err := r.Cookie( cookie_locale )
        if __err != nil {
            __cookie = &http.Cookie{
                Name: cookie_locale,
                Value: _lc,
                Path: "/", // for every page
                MaxAge: 0, // persistent cookie
            }
        } else {
            __cookie.Value = _lc
            __cookie.Path  = "/"
        }
        http.SetCookie( w, __cookie )
    }
    _check_locale := func() {
        if _get_locale() == "" {
            _set_locale( m.parent.default_locale )
        }
    }

    _url := r.URL.Path

    if _lcc > 0 {
        if len( _url ) > 1 {
            _parts := strings.SplitN( _url[1:], "/", 2 )
            if m.parent.localizer.hasLocale( _parts[0] ) {
                // Set locale cookie to loccale
                _locale    := _parts[0]
                if len( _parts ) > 1 {
                    r.URL.Path = "/" + _parts[1]
                } else {
                    r.URL.Path = "/"
                }
                _set_locale( _locale )
            } else {
                // No locale in the url
                    // redirect to default or locale in cookie
                    // redirect if lcc > 1
                _check_locale()
            }
        } else {
            // Root
                // Check if cookie has locale if not laod default and redirect
                // redirect if lcc > 1
            _check_locale()
        }
    }

    m.handler( w, r )
}

/*
func ( m *Mux ) Append( other *Mux ) {
    if m.parent != other.parent { return }
    m.handler = func ( w http.ResponseWriter, r *http.Request ) bool {
        if m.handler( w, r ) { return true }
        return other.handler( w, r )
    }
}

func ( m *Mux ) Prepend( other *Mux ) {
    if m.parent != other.parent { return }
    m.handler = func ( w http.ResponseWriter, r *http.Request ) bool {
        if other.handler( w, r ) { return true }
        return m.handler( w, r )
    }
}
*/

func NotFound( w http.ResponseWriter, r *http.Request ) bool {
    WriteStatus( w, 404, "Not Found" )
    return true
}

// Creates a symlink-like handler for target directory
//   Example: WithSymlink( "/home/www/somefolder/", "/link/" )
func ( m *Mux ) WithSymlink( target, link string ) {

    switch {
    case target == "",
        link == "",
        link[0] != '/': // link must start with a slash
        panic( ErrInvalidDirPath )
        return
    }

    target = path.Clean( target )
    link   = path.Clean( link )

    m.WithHandlerFunc( func ( w http.ResponseWriter, r *http.Request ) bool {

        url := r.URL.Path
        // Must not contain dotdot
        // Must have link as prefix
        if isSafeFileURL( url ) && strings.HasPrefix( url, link ) {

            // Prepend the target path to url
            var rel string // relative path
            if len( url ) > len( link ) {
                rel = url[len( link ):]
            } else {
                rel = ""
            }
            // rel is likely to have a slash at the beginning
            //   that slash gets removed while being cleaned below
            //   since two slashes become one slash
            p := path.Clean( target + "/" + rel ) // the file we are looking for
            b := path.Base( p )

            st, err := os.Stat( p )
            // Not found
            if os.IsNotExist( err ) {
                WriteStatus( w, 404, "Not Found" )
                return true
            }
            // Other errors
            if err != nil {
                // Status Internal Server Error
                WriteStatus( w, 500, "Internal Server Error" )
                return true
            }

            // When it's a directory
            // Later handle this with some option like: forbidDirectoryAccess
            if st.IsDir() {
                WriteStatus( w, 403, "Forbidden" )
                return true
            }

            f, err := os.Open( p )
            if err != nil {
                // Status Internal Server Error
                WriteStatus( w, 500, "Internal Server Error" )
                return true
            }

            // Set it to octet stream so that it won't be executed or compiled
            w.Header().Set( "Content-Type", "application/octet-stream" )
            w.Header().Set( "Content-Transfer-Encoding", "Binary" )
            w.Header().Set( "Content-Disposition", "attachment; filename=\"" + b + "\"" )

            http.ServeContent( w, r, b, st.ModTime(), f )
            f.Close()
            return true
        }

        return false

    } )

}

/*
 + The Asset-oriented part
 */

func ( m *Mux ) WithRealtimeAssetServer() {

    m.WithHandlerFunc( func ( w http.ResponseWriter, r *http.Request ) bool {

        url := r.URL.Path
        if isSafeAssetURL( url ) {
            // Remove the first slash at the beginning
            p   := path.Clean( url[1:] )
            b   := path.Base( p )
            ext := strings.ToLower( path.Ext( p ) )

            // Whitelist of Asset Extensions
            //   This is temporary security measure
            //   Liable to being removed or modified
            //   Later use some config var like: whitelistedExtensions

            found := false
            for _, i := range m.parent.whitelistedExtensions {
                if i == ext { found = true; break; }
            }
            if !found {
                // Status Not Found
                WriteStatus( w, 404, "Not Found" )
                return true
            }

            st, err := os.Stat( p )
            // Not found in the asset folder
            if os.IsNotExist( err ) {
                WriteStatus( w, 404, "Not Found" )
                return true
            }
            if err != nil {
                // Status Internal Server Error
                WriteStatus( w, 500, "Internal Server Error" )
                return true
            }
            if st.IsDir() {
                WriteStatus( w, 403, "Forbidden" )
                return true
            }

            // Process the asset
            err = m.parent.ProcessAsset( p )
            if err != nil {
                panic( err )
                WriteStatus( w, 500, "Internal Server Error" )
                return true
            }

            f, err := os.Open( p )
            if err != nil {
                // Status Internal Server Error
                WriteStatus( w, 500, "Internal Server Error" )
                return true
            }

            http.ServeContent( w, r, b, st.ModTime(), f )
            f.Close()
            return true

        }

        return false

    } )

}

func ( m *Mux ) WithCachedAssetServer() {

    m.WithHandlerFunc( func ( w http.ResponseWriter, r *http.Request ) bool {

        url := r.URL.Path
        if isSafeAssetURL( url ) {
            // Remove the first slash at the beginning
            p   := path.Clean( url[1:] )
            b   := path.Base( p )
            ext := path.Ext( b )

            // Whitelist of Asset Extensions
            //   This is temporary security measure
            //   Liable to being removed or modified
            switch ext {
            case ".css", ".js":
            default:
                // Status Not Found
                WriteStatus( w, 404, "Not Found" )
                return true
            }

            a, ok := m.parent.assets[p]
            if ok {
                // Check if Version Is Set
                v := r.FormValue( "v" )
                if v == "" || v != a.hash[:6] {
                    http.Redirect(
                        w, r, url + "?v=" + a.hash[:6],
                        http.StatusFound,
                    )
                    return true
                }

                // Serve Content is the Version Is Set
                http.ServeContent( w, r, b, a.modTime, a.reader )
            } else {
                // Status Not Found
                WriteStatus( w, 404, "Not Found" )
            }
            return true
        }

        return false

    } )

}

func isSafeFileURL( url string ) bool {
    return !containsDotDot( url )
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

/*
 + The Route-oriented part
 */

func makeMethodChecker( mflag int ) map[string] bool {

    _flag := make( map[string] bool )
    _flag[http.MethodGet]     = MethodGet     & mflag == MethodGet
    _flag[http.MethodHead]    = MethodHead    & mflag == MethodHead
    _flag[http.MethodPost]    = MethodPost    & mflag == MethodPost
    _flag[http.MethodPut]     = MethodPut     & mflag == MethodPut
    _flag[http.MethodPatch]   = MethodPatch   & mflag == MethodPatch
    _flag[http.MethodDelete]  = MethodDelete  & mflag == MethodDelete
    _flag[http.MethodConnect] = MethodConnect & mflag == MethodConnect
    _flag[http.MethodOptions] = MethodOptions & mflag == MethodOptions
    _flag[http.MethodTrace]   = MethodTrace   & mflag == MethodTrace

    return _flag

}

func ( m *Mux ) WithRoute( mflag int, rgx *regexp.Regexp, hf HandlerFunc ) {

    _flag := makeMethodChecker( mflag )
    m.WithHandlerFunc( func ( w http.ResponseWriter, r *http.Request ) bool {
        if _flag[r.Method] {
            if rgx.MatchString( r.URL.Path ) {
                hf( w, r )
                return true
            }
        }
        return false
    } )

}
/*
func ( m *Mux ) WithRoute2( mflag int, rgx *regexp.Regexp, hf HandlerFunc ) {

    _flag       := makeMethodChecker( mflag )
    _lcc        := len( m.parent.localizer.locales )
    _get_locale := func ( r *http.Request ) string {
        __cookie, __err := r.Cookie( cookie_locale )
        if __err != nil {
            return ""
        }
        return __cookie.Value
    }
    _set_locale := func ( w http.ResponseWriter, r *http.Request, lc string ) {
        __cookie, __err := r.Cookie( cookie_locale )
        if __err != nil {
            __cookie = &http.Cookie{
                Name: cookie_locale,
                Value: lc,
                Path: "/", // for every page
                MaxAge: 0, // persistent cookie
            }
        } else {
            __cookie.Value = lc
        }
        http.SetCookie( w, __cookie )
    }
    _check_locale := func ( w http.ResponseWriter, r *http.Request ) {
        if _get_locale( r ) == "" {
            _set_locale( w, r, m.parent.default_locale )
        }
    }

    m.WithHandlerFunc( func ( w http.ResponseWriter, r *http.Request ) bool {

        __url := r.URL.Path

        if _flag[r.Method] {

            if _lcc > 0 {
                if len( __url ) > 1 {
                    __parts := strings.SplitN( __url[1:], "/", 2 )
                    if m.parent.localizer.hasLocale( __parts[0] ) {
                        // Set locale cookie to loccale
                        __locale := __parts[0]
                        __url     = "/" + __parts[1]
                        _set_locale( w, r, __locale )
                    } else {
                        // No locale in the url
                            // redirect to default or locale in cookie
                            // redirect if lcc > 1
                        _check_locale( w, r )
                    }
                } else {
                    // Root
                        // Check if cookie has locale if not laod default and redirect
                        // redirect if lcc > 1
                    _check_locale( w, r )
                }
            }

            if rgx.MatchString( __url ) {
                hf( w, r )
                return true
            }

        }

        return false
    } )

}*/

func ( m *Mux ) WithTemplate( mflag int, rgx *regexp.Regexp, hf HandlerFunc ) {

    _flag := makeMethodChecker( mflag )

    m.WithHandlerFunc( func ( w http.ResponseWriter, r *http.Request ) bool {
        if _flag[r.Method] {

            __url := r.URL.Path
            if rgx.MatchString( __url ) {
                hf( w, r )
                return true
            }
        }
        return false
    } )

}

func ( m *Mux ) WithAuthenticator( auther Authenticator, realm string ) {

    m.WithHandlerFunc( func ( w http.ResponseWriter, r *http.Request ) bool {
        if auther( w, r ) {
            // returns false so that following handlers can handle the request
            return false
        }
        w.Header().Set( "WWW-Authenticate", "Basic realm=\"" + realm + "\"" )
        WriteStatus( w, 401, "Unauthorized" )
        // returns true since this is the last stop the request will reach
        return true
    } )

}

func ( m *Mux ) WithHandlerFunc( hf HandlerFunc ) {
    mh := m.handler
    m.handler = func ( w http.ResponseWriter, r *http.Request ) bool {
        if mh( w, r ) { return true }
        return hf( w, r )
    }
}