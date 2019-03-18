package network

type Server struct {
    body *http.Server
    spaces []*Space
}

func NewServer() *Server {
    srv := &Server{
        body: &http.Server{},
        spaces: make( []*Space, 0 )
    }
    srv.body.Handler = srv
    return srv
}

func( srv *Server ) Body() *http.Server {
    return srv.body
}

func( srv *Server ) Spaces() []*Space {
    return srv.spaces
}

func( srv *Server ) SetSpaces( spcs []*Space ) {
    srv.spaces = spcs
}

func( srv *Server ) AddSpace( spc *Space ) {
    srv.spaces = append( srv.spaces, spc )
}

// Server-related

func( srv *Server ) ServeRequest( req *Request ) {
    for _, spc := range srv.spaces {
        // Check
        if spc.ContainsAlias( req.body.URL.Host ) {
            spc.ServeRequest( req )
            return
        }
    }
}

func( srv *Server ) ServeHTTP( w http.ResponseWriter, r *http.Request ) {
    srv.ServeRequest( NewRequest( w, r ) )
}

func( srv *Server ) Start( addr string ) error {
    
    // Listener
    ln, err := net.Listen( "tcp", addr )
    if err != nil {
        return ErrStartingServer.Append( "addr:", addr, "err:", err )
    }
    
    // Serve
    err = srv.body.Serve( ln )
    return ErrServerExited.Append( "addr:", addr, "err:", err )
    
}

func( srv *Server ) StartTLS( addr, certFile, keyFile string ) error {
    
    // Listener
    ln, err := net.Listen( "tcp", addr )
    if err != nil {
        return ErrStartingServer.Append( "addr:", addr, "err:", err )
    }
    
    // Serve
    err = srv.body.ServeTLS( ln, certFile, keyFile )
    return ErrServerExited.Append( "addr:", addr, "err:", err )
    
}

func( srv *Server ) Stop() error {
    
    //
    err := srv.body.Shutdown( context.Background() )
    if err != nil {
        return ErrStoppingServer.Append( "err:", err )
    }
    return nil
    
}