package volume

import (
    "path/filepath"
    "strings"

    "../system"
)

/*
 + fs.go
 *
 * fs.go has methods that deal with the file system
 */


/*
 + ProcessAsset
 *
 * @param path - path of the file to process
 * @return error
 */

func ProcessAsset( path string ) error {

    bwoext := BaseWithoutExt( path )
    base   := filepath.Base( path )
    ext    := strings.ToLower( filepath.Ext( path ) )
    dir    := filepath.ToSlash( filepath.Dir( path ) )

    switch ext {
    case ".scss", ".sass":

        if !fl_sassInstalled {
            return ErrUnableToProcessAsset.Append( path )
        }

        cmdArgs := []string{
            "sass",
            dir + "/" + base,
            dir + "/" + bwoext + ".css",
        }

        err := system.Exec( cmdArgs... )
        if err != nil {
            return ErrUnableToProcessAsset.Append( path )
        }

    }

    return nil

}