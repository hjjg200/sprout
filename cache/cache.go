package cache

import (
    "archive/zip"
    "bytes"
    "io"
    "time"

    "github.com/hjjg200/together"
)

type Cache struct {
    data []byte
    buf  *bytes.Buffer
    hs   *together.HoldSwitch
    zrd  *zip.Reader
    zwr  *zip.Writer
}

func NewCache() *Cache {

    chc := &Cache{
        data: nil,
        buf: nil,
        hs: nil,
        zrd: nil,
        zwr: nil,
    }

    hs := together.NewHoldSwitch()
    chc.hs = hs
    hs.Handlers( switchRead, chc.beginRead, chc.endRead )
    hs.Handlers( switchWrite, chc.beginWrite, chc.endWrite )

    return chc

}

// Zip entry wrapper

type entryWriter struct {
    parent *Cache
    wr io.Writer
}

type entryReadCloser struct {
    parent *Cache
    rc io.ReadCloser
}

func( ew *entryWriter ) Write( p []byte ) ( n int, err error ) {
    n, err = ew.wr.Write( p )
    return
}

func( ew *entryWriter ) Close() error {
    ew.parent.hs.Done( switchWrite )
    return nil
}

func( er *entryReadCloser ) Read( p []byte ) ( n int, err error ) {
    n, err = er.rc.Read( p )
    return
}

func( er *entryReadCloser ) Close() error {
    er.parent.hs.Done( switchRead )
    return er.rc.Close()
}

// Methods

func( chc *Cache ) Open( path string ) ( io.ReadCloser, error ) {

    // Hold
    chc.hs.Add( switchRead, 1 ) // Close will call Done of HoldSwitch

    // Fetch
    for _, f := range chc.zrd.File {
        if f.Name == path {
            rc, err := f.Open()
            if err != nil {
                return nil, ErrEntryAccessFailed.Append( path )
            }
            return &entryReadCloser{
                rc: rc,
                parent: chc,
            }, nil
        }
    }

    // Error
    return nil, ErrEntryNotFound.Append( path )

}

func( chc *Cache ) Create( path string, mt time.Time ) ( io.WriteCloser, error ) {

    // Hold
    chc.hs.Add( switchWrite, 1 ) // Close will call Done of HoldSwitch

    // Fetch
    wr, err := chc.zwr.CreateHeader( &zip.FileHeader{
        Name: path,
        Modified: mt,
    } )
    if err != nil {
        return nil, ErrEntryWriteFail.Append( path )
    }
    return &entryWriter{
        wr: wr,
        parent: chc,
    }, nil

}

func( chc *Cache ) beginWrite() {

    // Begin Write
    chc.buf = bytes.NewBuffer( nil )
    chc.zwr = zip.NewWriter( chc.buf )

    // Read the previous content
    if len( chc.data ) > 0 {
        brd    := bytes.NewReader( chc.data )
        zrd, _ := zip.NewReader( brd, brd.Size() )

        for _, f := range zrd.File {

            // File
            r, err := f.Open()
            if err != nil {
                panic( err )
            }

            // Write
            w, err := chc.zwr.CreateHeader( &f.FileHeader )
            if err != nil {
                panic( err )
            }
            io.Copy( w, r )

        }
    }

}

func( chc *Cache ) endWrite() {

    // End Writing to Zip
    chc.zwr.Close()
    chc.data = chc.buf.Bytes()
    chc.buf  = nil
    chc.zwr  = nil

}

func( chc *Cache ) beginRead() {

    // Begin Read
    rd         := bytes.NewReader( chc.data )
    chc.zrd, _  = zip.NewReader( rd, rd.Size() )

}

func( chc *Cache ) endRead() {

    // End Read
    chc.zrd = nil

}