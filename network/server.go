package network

import (
    "context"
    "fmt"
    "net"
    "net/http"

    "github.com/hjjg200/sprout/util/errors"
)

type Server struct {
    addr     string
    body     *http.Server
    spaces   []*Space
    keyFile  string
    certFile string
}

func NewServer() *Server {
    srv := &Server{
        body: &http.Server{},
        spaces: make( []*Space, 0 ),
    }
    srv.body.Handler = srv
    return srv
}

func( srv *Server ) Addr() string {
    return srv.addr
}

func( srv *Server ) SetAddr( addr string ) {
    srv.addr = addr
}

func( srv *Server ) SetPort( port int16 ) {
    srv.addr = fmt.Sprintf( "0.0.0.0:%d", port )
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

func( srv *Server ) KeyFile() string {
    return srv.keyFile
}

func( srv *Server ) CertFile() string {
    return srv.certFile
}

func( srv *Server ) ConfigureTls( certFile, keyFile string ) {
    srv.certFile = certFile
    srv.keyFile  = keyFile
}

func( srv *Server ) DisableTls() {
    srv.ConfigureTls( "", "" )
}

// Server-related

func( srv *Server ) ServeRequest( req *Request ) {
    for _, spc := range srv.spaces {
        // Check host
        if spc.ContainsHost( req.body.Host ) {
            spc.serveRequest( req )
            return
        }
    }

    // Bad Request if not found
    req.PopBlank( 400 )
}

func( srv *Server ) ServeHTTP( w http.ResponseWriter, r *http.Request ) {
    srv.ServeRequest( NewRequest( w, r ) )
}

func( srv *Server ) Start() error {

    // Listener
    ln, err := net.Listen( "tcp", srv.addr )
    if err != nil {
        return errors.ErrServerOperation.Raise( "starting server", "addr:", srv.addr, "err:", err )
    }

    // Serve
    if srv.keyFile != "" && srv.certFile != "" {
        err = srv.body.ServeTLS( ln, srv.certFile, srv.keyFile )
    } else {
        err = srv.body.Serve( ln )
    }

    return errors.ErrServerExited.Raise( "addr:", srv.addr, "err:", err )

}

func( srv *Server ) Stop() error {

    //
    err := srv.body.Shutdown( context.Background() )
    if err != nil {
        return errors.ErrServerOperation.Raise( "err:", err )
    }
    return nil

}
