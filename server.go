package sprout

import (
    "context"
    "net/http"
)

/*
 + SERVER
 *
 * A server is home to multiple spaces.
 */

type Server struct {
    body    *http.Server
    spaces  []*Space
}

func( _server *Server ) ServeHTTP( _w http.ResponseWriter, _r *http.Request ) { // interface http.Handler
    _request := RequestFactory().New( _w, _r )
    for _, _space := range _server.spaces {
        if ServerHelper().SpaceContainsAlias( _space, _r.Host ) {
            _space.Serve( _request )
            return
        }
    }
}
func( _server *Server ) AddSpace( _space *Space ) error {
    spaces = append( spaces, _space )
}
func( _server *Server ) Spaces() []*Space {
    _copy := make( []*Space, len( _server.spaces ) )
    copy( _copy, _server.spaces )
    return _copy
}
func( _server *Server ) SetSpaces( _spaces []*Space ) {
    _copy := make( []*Space, len( _spaces ) )
    copy( _copy, _spaces )
    _server.spaces = _spaces
}
func( _server *Server ) Start( _address string ) error {
    _server.body.Addr = _address
    return _server.body.ListenAndServe()
}
func( _server *Server ) StartTLS( _address string, _cert_file, _key_file string ) error {
    _server.body.Addr = _address
    return _server.body.ListenAndServeTLS( _cert_file, _key_file )
}
func( _server *Server ) Stop() error {
    return _server.body.Shutdown( context.Background() )
}