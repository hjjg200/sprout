package volume

import (
    "fmt"
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

    fmt.Println( "--- CSS" )
    write( testDir + "/" + "asset/a.css", `body {
    line-height: 1.5; /* FIRST WRITE */
}` )
    ast := rtv.Asset( "asset/a.css" )
    fmt.Println( string( ast.Bytes() ) )

    write( testDir + "/" + "asset/a.css", `body {
    line-height: 2.22; /* SECOND WRITE */
}` )
    ast = rtv.Asset( "asset/a.css" )
    fmt.Println( string( ast.Bytes() ) )

    write( testDir + "/" + "asset/a.css", `body {
    color: red;
    line-height: 2.22; /* THIRD WRITE */
}` )
    ast = rtv.Asset( "asset/a.css" )
    fmt.Println( string( ast.Bytes() ) )

    // Localizer
    fmt.Println( "--- LOCALIZER" )
    write( testDir + "/" + "i18n/en.json", `{
    "en": {
        "ABC": "first"
    }
}` )
    lc := rtv.Localizer( "en" )
    if lc == nil {
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
    lc = rtv.Localizer( "en" )
    fmt.Println( lc.L( `{% ABC %}` ) )
    write( testDir + "/" + "i18n/en.json", `{
    "en": {
        "ABC": "third"
    }
}` )
    lc = rtv.Localizer( "en" )
    fmt.Println( lc.L( `{% ABC %}` ) )

    // Compile
    fmt.Println( "--- COMPILE" )
    
    ast = rtv.Asset( "asset/b.css" )
    if ast == nil {
        t.Error( "not found" )
        return
    }
    fmt.Println( string( ast.Bytes() ) )
    
}