package sprout

import (
    "./network"
    "./util"
)

type Sprout struct {
    servers []*network.Server
    logger  *util.Logger
}

func New() *Sprout {

    //
    return &Sprout{
        servers: make( []*network.Server, 0 ),
        logger: util.NewLogger(),
    }

}

// Getters & Setters

func( sprt *Sprout ) Servers() []*network.Server {
    return sprt.servers
}

func( sprt *Sprout ) SetServers( srvs []*network.Server ) {
    sprt.servers = srvs
}

func( sprt *Sprout ) AddServer( srv *network.Server ) {
    sprt.servers = append( sprt.servers, srv )
}

// General

func( sprt *Sprout ) StartAll() {
    for _, srv := range sprt.servers {
        err := srv.Start()
        if err != nil {
            sprt.logger.Warnln( err )
        }
    }
}

func( sprt *Sprout ) StopAll() {
    for _, srv := range sprt.servers {
        err := srv.Stop()
        if err != nil {
            sprt.logger.Warnln( err )
        }
    }
}