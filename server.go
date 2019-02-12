package sprout

/*
 + SERVER
 *
 * A server is home to multiple spaces.
 */

type Server struct {
	body   *http.Server
    spaces []*Space
    port   uint16
}

func( _server *Server ) ServeHTTP( _w http.ResponseWriter, _r *http.Request ) {} // interface http.Handler
func( _server *Server ) AddSpace( _space *Space ) error {}
func( _server *Server ) SetPort( _port uint16 ) {}
func( _server *Server ) Start() error {}
func( _server *Server ) StartTLS( _cert_file, _key_file string ) error {}
func( _server *Server ) Stop() error {}