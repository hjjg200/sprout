package volume

import (
    "testing"
)

func TestVolume01( t *testing.T ) {

    css, err := CompileScss( []byte( `
body {
    div {
        .white {
            color: white;
        }
    }
}
` ) )
    if err != nil {
        t.Log( err )
    }
    t.Log( string( css ) )

}