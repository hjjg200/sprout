package volume

import (
    "path/filepath"
)

func BaseWithoutExt( path string ) string {
    base := filepath.Base( path )
    ext  := filepath.Ext( path )
    if len( ext ) == 0 {
        return base
    }
    return base[:len( base ) - len( ext )]
}
