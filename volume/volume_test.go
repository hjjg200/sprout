package volume

import (
    "testing"
    "os"
)

func TestVolume01( t *testing.T ) {

    // /test/TestVolume01

    vol := NewVolume()
    vol.ImportDirectory( "../test/TestVolume01/volume" )

    f, err := os.OpenFile(
        "../test/TestVolume01/export.zip",
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
    f.Write( chc.Data() )
    f.Close()

    t.Log( vol.i18n.L( "en", "{% ab %}" ) )

}

func TestVolume02( t *testing.T ) {

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