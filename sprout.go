package sprout

import (
    "strings"
    "time"

    "./network"
    "./environ"
)

type Sprout struct {
    servers []*network.Server
}

func New() *Sprout {

    //
    return &Sprout{
        servers: make( []*network.Server, 0 ),
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

    environ.Logger.OKln(
        strings.ToUpper( environ.AppName ),
        environ.AppVersion,
        "UP AND RUNNING since",
        time.Now().Unix(),
    )

    for _, srv := range sprt.servers {
        err := srv.Start()
        if err != nil {
            environ.Logger.Warnln( err )
        }
    }

}

func( sprt *Sprout ) StopAll() {
    for _, srv := range sprt.servers {
        err := srv.Stop()
        if err != nil {
            environ.Logger.Warnln( err )
        }
    }
}