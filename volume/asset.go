package volume

import (
    "bytes"
    "crypto/md5"
    "fmt"
    "io"
    "mime"
    "path/filepath"
    "time"
)

type Asset struct {
    mimeType string
    *bytes.Reader
    modTime time.Time
    version string // first 6 letter of md5 hash of unix time string of modTime
}

/*
 + MakeAsset
 *
 * @param r - reader that contains the content for the asset
 * @param mt - modified time of the asset
 */

func NewAsset( name string, r io.Reader, mt time.Time ) *Asset {

    // Basics
    ast := &Asset{
        mimeType: mime.TypeByExtension( filepath.Ext( name ) ),
        modTime: mt,
    }

    // Check mime
    if ast.mimeType == "" {
        ast.mimeType = "text/plain"
    }

    // Version
    hash := md5.New()
    hash.Write( []byte( fmt.Sprint( mt.Unix() ) ) )
    ast.version = fmt.Sprintf( "%x", hash.Sum( nil ) )[:6]

    // Data
    buf := bytes.NewBuffer( nil )
    io.Copy( buf, r )
    ast.Reader = bytes.NewReader( buf.Bytes() )

    return ast

}
