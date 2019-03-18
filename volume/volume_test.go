package volume

import (
    "testing"
    "path/filepath"
    "os"
)

func Test1( t *testing.T ) {

    // /test/volume_test_01

    vol := NewVolume()
    err := filepath.Walk( "../test/volume_test_01/volume", vol.WalkFuncBasedOn( "../test/volume_test_01/volume" ) )
    if err != nil {
        t.Error( err )
    }
    f, err := os.OpenFile(
        "../test/volume_test_01/export.zip",
        os.O_CREATE | os.O_TRUNC | os.O_WRONLY,
        0644,
    )
    if err != nil {
        t.Error( err )
    }
    err = vol.ExportZip( f )
    if err != nil {
        t.Error( err )
    }

    t.Log( vol.i18n.L( "en", "{% ab %}" ) )

}