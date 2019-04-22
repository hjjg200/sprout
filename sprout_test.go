package sprout

import (
    "bytes"
    "github.com/hjjg200/sprout/network"
    "github.com/hjjg200/sprout/volume"
    "testing"
)

func TestSprout01( t *testing.T ) {

    sprt := New()

    srv := network.NewServer()
    srv2 := network.NewServer()
    space := network.NewSpace()
    space2 := network.NewSpace()

    sprt.AddServer( srv )
    sprt.AddServer( srv2 )
    srv.AddSpace( space )
    space.AddAlias( "127.0.0.1" )
    srv.SetPort( 8002 )
    srv2.AddSpace( space2 )
    srv2.SetPort( 8003 )

    vol := volume.NewRealtimeVolume( "./test/TestSprout01" )

    space2.WithReverseProxy( "http://127.0.0.1:8002" )

    space.SetVolume( vol )
    space.WithHandler( network.HandlerFactory.BasicAuth( func( id, pw string ) bool {
        return id == "root" && pw == "root"
    }, "" ) )
    space.WithRoute( "^/(index.html?)?$", []string{ "GET" }, space.TemplateHandler(
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
    space.WithRoute( "^/a([a-z])?([a-z])?([a-z])?e$", []string{ "GET" }, space.TemplateHandler(
        "template/abcde.html",
        func( req *network.Request ) interface{} {
            return map[string] []string {
                "vars": req.Vars(),
            }
        },
    ) )
    space.WithRoute( "^/stop$", []string{ "GET" }, func( req *network.Request ) bool {
        srv.Stop()
        return true
    } )
    space.WithRoute( "^/error$", []string{ "GET" }, network.HandlerFactory.Status( 500 ) )
    space.WithRoute(
        "^/some.txt$",
        []string{ "GET" },
        func( req *network.Request ) bool {
            ast := space.Volume().Asset( "asset/some.txt" )
            rdskr := bytes.NewReader( ast.Bytes() )
            req.PopAttachment( ast.Name(), ast.ModTime(), rdskr )
            return true
        },
    )
    space.WithAssetServer( "/asset/" )
    space.WithHandler( network.HandlerFactory.Status( 404 ) )

    sprt.StartAll()

}