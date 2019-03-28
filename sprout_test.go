package sprout

import (
    "./network"
    "./volume"
    "testing"
)

func TestSprout01( t *testing.T ) {

    sprt := New()

    srv := network.NewServer()
    space := network.NewSpace( "" )

    sprt.AddServer( srv )
    srv.AddSpace( space )
    srv.SetAddr( ":8002" )

    vol := volume.NewRealtimeVolume( "./test/TestSprout01" )
    space.SetVolume( vol )
    tmpl, _ := space.Volume().Template( "template/index.html" )
    space.WithRoute( "^/(index.html?)?$", network.HandlerFactory.Template(
        tmpl, func( req *network.Request ) interface{} {
           return map[string] string {
               "hello": "HELLO WORLD",
           }
        } ) )
    space.WithRoute( "^/stop$", func( req *network.Request ) bool {
        srv.Stop()
        return true
    } )

    sprt.StartAll()

}