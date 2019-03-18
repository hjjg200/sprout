package sprout

type Sprout struct {
    servers []*Server
}

func New() *Sprout {
    
    //
    return &Sprout{
        servers: make( []*Server, 0 )
    }
    
}

// Getters & Setters

func( sprt *Sprout ) Servers() []*Server {
    return sprt.servers
}

func( sprt *Sprout ) SetServers( srvs []*Server ) {
    sprt.servers = srvs
}

func( sprt *Sprout ) AddServer( srv *Server ) {
    sprt.servers = append( sprt.servers, srv )
}
