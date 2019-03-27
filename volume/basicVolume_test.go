package volume

import (
    "testing"
    "os"
)

func TestBasicVolume01( t *testing.T ) {

    // /test/TestBasicVolume01

    vol := NewBasicVolume()
    vol.ImportDirectory( "../test/TestBasicVolume01/volume" )

    f, err := os.OpenFile(
        "../test/TestBasicVolume01/export.zip",
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
