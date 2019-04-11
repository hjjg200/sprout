package sprout

import (
    "./network"
    "./volume"
    "./util"
    "testing"
)

func TestSprout01( t *testing.T ) {

    sprt := New()

    srv := network.NewServer()
    space := network.NewSpace( "127.0.0.1" )

    sprt.AddServer( srv )
    srv.AddSpace( space )
    srv.SetAddr( ":8002" )
    lgr := util.NewLogger()

    lgr.SetTimeType( util.LogTimeTypeSeconds )
    lgr.OKln( "OK" )
    lgr.Warnln( "Warn" )

    vol := volume.NewRealtimeVolume( "./test/TestSprout01" )
    space.SetVolume( vol )
    space.WithHandler( network.HandlerFactory.BasicAuth( func( id, pw string ) bool {
        return id == "root" && pw == "root"
    }, "" ) )
    space.WithRoute( "^/(index.html?)?$", space.TemplateHandler(
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
    space.WithRoute( "^/a([a-z])?([a-z])?([a-z])?e$", space.TemplateHandler(
        "template/abcde.html",
        func( req *network.Request ) interface{} {
            return map[string] []string {
                "vars": req.Vars(),
            }
        },
    ) )
    space.WithRoute( "^/stop$", func( req *network.Request ) bool {
        srv.Stop()
        return true
    } )
    space.WithRoute( "^/log$", func( req *network.Request ) bool {
        lgr.OKln( "OK" )
        return true
    } )
    space.WithAssetServer( "/asset/" )

    sprt.StartAll()

}