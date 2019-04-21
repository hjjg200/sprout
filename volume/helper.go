package volume

import (
    "path/filepath"
    "strings"

    "github.com/hjjg200/sprout/util"
)

func BaseWithoutExt( path string ) string {
    base := filepath.Base( path )
    ext  := filepath.Ext( path )
    if len( ext ) == 0 {
        return base
    }
    return base[:len( base ) - len( ext )]
}

func typeByPath( path string ) int {

    ext := util.String( filepath.Ext( path ) )

    switch {
    case strings.HasPrefix( path, c_assetDirectory + "/" ):
        return c_typeAsset
    case strings.HasPrefix( path, c_i18nDirectory + "/" ):
        if ext.IsIn( ".json" ) {
            return c_typeI18n
        }
    case strings.HasPrefix( path, c_templateDirectory + "/" ):
        if ext.IsIn( ".html", ".htm" ) {
            return c_typeTemplate
        }
    }

    return c_typeNull

}
