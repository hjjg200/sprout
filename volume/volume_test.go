package volume

import (
    "testing"
)

func TestVolume01( t *testing.T ) {

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