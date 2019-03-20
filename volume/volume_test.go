package volume

import (
    "testing"
    "os"
)

func Test1( t *testing.T ) {

    // /test/volume_test_01

    vol := NewVolume()
    vol.ImportDirectory( "../test/volume_test_01/volume" )

    f, err := os.OpenFile(
        "../test/volume_test_01/export.zip",
        os.O_CREATE | os.O_TRUNC | os.O_WRONLY,
        0644,
    )
    if err != nil {
        t.Error( err )
    }
    chc, err := vol.Export()
    if err != nil {
        t.Error( err )
    }
    t.Log( len( chc.Files() ) )
    f.Write( chc.Data() )
    f.Close()
    
    t.Log( vol.i18n.L( "en", "{% ab %}" ) )

}

func Test2( t *testing.T ) {
    
    css, err := CompileScss( `
body {
    div {
        .white {
            color: white;
        }
    }
}
` )
    if err != nil {
        t.Log( err )
    }
    t.Log( css )
    
}