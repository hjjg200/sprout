package network

type Server struct {
    body *http.Server
    spaces []*Space
}

func NewServer() *Server {
    srv := &Server{
        spaces: make( []*Space, nil )
    }
    srv.body = &http.Server{
        Handler: srv,
    }
    return srv
}

func( srv *Server ) ServeHTTP( w http.ResponseWriter, r *http.Request ) {
    
}