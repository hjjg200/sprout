package sprout

/*
 + SPROUT
 *
 * A sprout is a set of servers that handle HTTP requests
 */

type Sprout struct {
    servers []*Server
}

func New() *Sprout {

    // Check the OS
    switch runtime.GOOS {
    case "darwin", "windows", "linux":
    default:
        panic( SproutVariables().ErrorOSNotSupported() )
        return nil
    }

    return &Sprout{
        servers: make( []*Server, 0 )
    }

}
func( _sprout *Sprout ) StartAll() error {}
func( _sprout *Sprout ) StopAll() error {}
func( _sprout *Sprout ) AddServer( _server *Server ) {}