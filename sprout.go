package sprout

import (
    "fmt"
    "strings"
    "sync"
    "time"

    "github.com/hjjg200/sprout/network"
    "github.com/hjjg200/sprout/environ"
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

    now := time.Now()
    environ.Logger.OKln(
        strings.ToUpper( "🍀" + environ.AppName ),
        environ.AppVersion,
        "UP AND RUNNING since",
        fmt.Sprintf( "%d.%03d", now.Unix(), now.Nanosecond() / 1e+6 ),
    )

    var wg sync.WaitGroup
    wg.Add( len( sprt.servers ) )

    for _, srv := range sprt.servers {
        go func( i *network.Server ) {
            err := i.Start()
            if err != nil {
                environ.Logger.Warnln( err )
            }
            wg.Done()
        }( srv )
    }

    wg.Wait()

}

func( sprt *Sprout ) StopAll() {
    for _, srv := range sprt.servers {
        err := srv.Stop()
        if err != nil {
            environ.Logger.Warnln( err )
        }
    }
}
