package sprout

import (
    "./network"
    "testing"
)

func TestSprout1( t *testing.T ) {

    sprt := New()

    srv := network.NewServer()
    space := network.NewSpace( "" )

    sprt.AddServer( srv )
    srv.AddSpace( space )
    srv.SetAddr( ":8002" )

    space.WithRoute( "^/$", func( req *network.Request ) bool {
        println( "at root" )
        srv.Stop()
        return true
    } )

    sprt.StartAll()

}