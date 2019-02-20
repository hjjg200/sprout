package sprout

/*
 + SERVER FACTORY
 *
 * A server factory is a pseudo-static member that is responsible for making servers
 */

type server_factory struct {}
var  static_server_factory = &server_factory{}

func ServerFactory() *server_factory {
    return static_server_factory
}
func( _srvfac *server_factory ) New() *Server {
    return &Server{
        body: &http.Server{
            Handler: _server,
        },
        spaces: make( []*Space, 0 ),
    }
}