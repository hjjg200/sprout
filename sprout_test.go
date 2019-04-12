package sprout

import (
    "./network"
    "./volume"
    "./util"
    "./environ"
    "testing"
)

func init() {
    environ.Logger.SetTimeType( util.LogTimeTypeSeconds )
}

func TestSprout01( t *testing.T ) {

    sprt := New()

    srv := network.NewServer()
    space := network.NewSpace()

    sprt.AddServer( srv )
    srv.AddSpace( space )
    srv.SetPort( 8002 )

    vol := volume.NewRealtimeVolume( "./test/TestSprout01" )
    space.SetVolume( vol )
    space.WithHandler( network.HandlerFactory.BasicAuth( func( id, pw string ) bool {
        return id == "root" && pw == "root"
    }, "" ) )
    space.WithRoute( "^/(index.html?)?$", network.MethodGet, space.TemplateHandler(
        "template/index.html",
        func( req *network.Request ) interface{} {
            return map[string] interface{} {
                "hello": []string{
                    "abc",
                    "def",
                    "ghi",
                },
            }
        },
    ) )
    space.WithRoute( "^/a([a-z])?([a-z])?([a-z])?e$", network.MethodGet, space.TemplateHandler(
        "template/abcde.html",
        func( req *network.Request ) interface{} {
            return map[string] []string {
                "vars": req.Vars(),
            }
        },
    ) )
    space.WithRoute( "^/stop$", network.MethodGet, func( req *network.Request ) int {
        srv.Stop()
        return 200
    } )
    space.WithRoute( "^/error$", network.MethodGet, network.HandlerFactory.Status( 500 ) )
    space.WithAssetServer( "/asset/" )
    space.WithHandler( network.HandlerFactory.Status( 404 ) )

    sprt.StartAll()

}