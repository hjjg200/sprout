package volume

import (
    "fmt"
    "io"
    "os"
    "testing"
)

func TestRealtimeVolume01( t *testing.T ) {

    testDir := "../test/TestRealtimeVolume01"

    write := func( p string, c string ) {
        f, err := os.OpenFile( p, os.O_WRONLY | os.O_TRUNC | os.O_CREATE, 0644 )
        if err != nil {
            return
        }
        defer f.Close()
        f.Write( []byte( c ) )
    }

    rtv := NewRealtimeVolume( testDir )

    write( testDir + "/" + "asset/a.css", `body {
    line-height: 1.5; /* FIRST WRITE */
}` )
    ast, _ := rtv.Asset( "asset/a.css" )
    io.Copy( os.Stdout, ast )
    fmt.Print( "\n" )

    write( testDir + "/" + "asset/a.css", `body {
    line-height: 2.22; /* SECOND WRITE */
}` )
    ast, _ = rtv.Asset( "asset/a.css" )
    io.Copy( os.Stdout, ast )
    fmt.Print( "\n" )

    write( testDir + "/" + "asset/a.css", `body {
    color: red;
    line-height: 2.22; /* THIRD WRITE */
}` )
    ast, _ = rtv.Asset( "asset/a.css" )
    io.Copy( os.Stdout, ast )
    fmt.Print( "\n" )

    // Localizer
    write( testDir + "/" + "i18n/en.json", `{
    "en": {
        "ABC": "first"
    }
}` )
    lc, ok := rtv.Localizer( "en" )
    if !ok {
        t.Error( "no localizer" )
        fmt.Println( rtv.I18n().LocaleNames() )
        return
    }
    fmt.Println( lc.L( `{% ABC %}` ) )
    write( testDir + "/" + "i18n/en.json", `{
    "en": {
        "ABC": "second"
    }
}` )
    lc, ok = rtv.Localizer( "en" )
    fmt.Println( lc.L( `{% ABC %}` ) )
    write( testDir + "/" + "i18n/en.json", `{
    "en": {
        "ABC": "third"
    }
}` )
    lc, ok = rtv.Localizer( "en" )
    fmt.Println( lc.L( `{% ABC %}` ) )

}