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
    name string
    mimeType string
    bytes []byte
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
        name: name,
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
    ast.bytes = buf.Bytes()

    return ast

}

func( ast *Asset ) Bytes() []byte {
    return ast.bytes
}

func( ast *Asset ) Name() string {
    return ast.name
}

func( ast *Asset ) MimeType() string {
    return ast.mimeType
}

func( ast *Asset ) ModTime() time.Time {
    return ast.modTime
}
