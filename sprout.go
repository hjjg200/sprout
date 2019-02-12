package sprout

/*
 + SPROUT
 *
 * A sprout is a set of servers that handle HTTP requests
 */

type Sprout struct {
    servers []*Server
}

func New() *Sprout {}
func( _sprout *Sprout ) StartAll() error {}
func( _sprout *Sprout ) StopAll() error {}
func( _sprout *Sprout ) AddServer( _server *Server ) {}