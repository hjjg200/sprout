package sprout

import (
    "./network"
    "./volume"
    "testing"
    "html/template"
)

func TestSprout01( t *testing.T ) {

    sprt := New()

    srv := network.NewServer()
    space := network.NewSpace( "127.0.0.1" )

    sprt.AddServer( srv )
    srv.AddSpace( space )
    srv.SetAddr( ":8002" )

    vol := volume.NewRealtimeVolume( "./test/TestSprout01" )
    space.SetVolume( vol )
    space.WithHandler( network.HandlerFactory.BasicAuth( func( id, pw string ) bool {
        return id == "root" && pw == "root"
    }, "" ) )
    space.WithRoute( "^/(index.html?)?$", network.HandlerFactory.Template(
        func() *template.Template {
            tmpl, _ := space.Volume().Template( "template/index.html" )
            return tmpl
        }, func( req *network.Request ) interface{} {
           return map[string] interface{} {
               "hello": []string{
                   "abc",
                   "def",
                   "ghi",
               },
           }
        } ) )
    space.WithRoute( "^/stop$", func( req *network.Request ) bool {
        srv.Stop()
        return true
    } )
    space.WithAssetServer( "/asset/" )

    sprt.StartAll()

}