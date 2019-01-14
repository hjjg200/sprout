package sprout

import (
    "context"
    "errors"
    "net"
    "net/http"
)

type Server struct {
    body *http.Server
    mux  *Mux
}

var (
    ErrServerExists = errors.New( "sprout: the server already exists" )
)

func ( s *Sprout ) NewServer( key string ) ( *Server, error ) {
    srv, ok := s.servers[key]
    if ok {
        return nil, ErrServerExists
    }
    srv = &Server{
        body: &http.Server{},
        mux: s.NewMux(),
    }
    s.servers[key] = srv
    return srv, nil
}

func ( s *Sprout ) Server( key string ) ( *Server, bool ) {
    srv, ok := s.servers[key]
    return srv, ok
}

func ( srv *Server ) Mux() *Mux {
    return srv.mux
}

func ( srv *Server ) SetMux( m *Mux ) {
    srv.mux = m
}

func ( srv *Server ) Start( addr string ) error {
    l, err := net.Listen( "tcp", addr )
    if err != nil {
        return err
    }
    srv.body.Handler = srv.mux
    return srv.body.Serve( l )
}

func ( srv *Server ) StartTLS( addr, certFile, keyFile string ) error {
    l, err := net.Listen( "tcp", addr )
    if err != nil {
        return err
    }
    srv.body.Handler = srv.mux
    return srv.body.ServeTLS( l, certFile, keyFile )
}

func ( srv *Server ) Stop() error {
    return srv.body.Shutdown( context.Background() )
}