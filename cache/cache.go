package cache

import (
    "archive/zip"
    "bytes"
    "io"
    "time"

    "github.com/hjjg200/sprout/util/errors"
    "github.com/hjjg200/together"
)

type Cache struct {
    data []byte
    buf  *bytes.Buffer
    hs   *together.HoldSwitch
    zrd  *zip.Reader
    zwr  *zip.Writer
    mode int
}

func NewCache() *Cache {

    chc := &Cache{
        data: nil,
        buf: nil,
        hs: nil,
        zrd: nil,
        zwr: nil,
        mode: -1,
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

func( ew *entryWriter ) Write( p []byte ) ( n int, err error ) {
    n, err = ew.wr.Write( p )
    return
}

func( ew *entryWriter ) Close() error {
    ew.parent.hs.Done( switchWrite )
    return nil
}

type entryReadCloser struct {
    parent *Cache
    rc io.ReadCloser
}

func( er *entryReadCloser ) Read( p []byte ) ( n int, err error ) {
    n, err = er.rc.Read( p )
    return
}

func( er *entryReadCloser ) Close() error {
    er.parent.hs.Done( switchRead )
    return er.rc.Close()
}

type entry struct {
    file   *zip.File
    parent *Cache

    zip.FileHeader
}

func( e *entry ) Open() ( io.ReadCloser, error ) {

    // Hold
    e.parent.hs.Add( switchRead, 1 ) // Close will call Done of HoldSwitch

    // Open
    rc, err := e.file.Open()
    if err != nil {
        return nil, errors.ErrIOError.Raise( "failed to access", e.Name )
    }
    return &entryReadCloser{
        rc: rc,
        parent: e.parent,
    }, nil

}

// Methods

func( chc *Cache ) Flush() error {
    if !chc.hs.IsEmpty() {
        return chc.hs.Close()
    }
    return nil
}

func( chc *Cache ) Data() []byte {

    chc.hs.Add( switchRead, 1 )
    defer chc.hs.Done( switchRead )

    // Make copy since chc.data is liable to changes at any moment
    cp := make( []byte, len( chc.data ) )
    copy( cp, chc.data )

    return cp

}

func( chc *Cache ) Files() []*entry {

    // Hold
    chc.hs.Add( switchRead, 1 )
    defer chc.hs.Done( switchRead )

    entries := make( []*entry, 0 )
    for _, f := range chc.zrd.File {
        entries = append( entries, &entry{
            FileHeader: f.FileHeader,
            file: f,
            parent: chc,
        } )
    }
    return entries

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
        return nil, errors.ErrIOError.Raise( "failed to write", path )
    }
    return &entryWriter{
        wr: wr,
        parent: chc,
    }, nil

}

func( chc *Cache ) beginWrite() {

    // Begin Write
    chc.mode = switchWrite
    chc.buf  = bytes.NewBuffer( nil )
    chc.zwr  = zip.NewWriter( chc.buf )

    // Read the previous content
    if len( chc.data ) > 0 {
        brd      := bytes.NewReader( chc.data )
        zrd, err := zip.NewReader( brd, brd.Size() )
        if err != nil {
            panic( err )
        }

        for _, f := range zrd.File {

            // File
            rc, err := f.Open()
            if err != nil {
                panic( err )
            }

            // Write
            w, err := chc.zwr.CreateHeader( &f.FileHeader )
            if err != nil {
                panic( err )
            }
            io.Copy( w, rc )
            rc.Close()

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
    chc.mode    = switchRead
    rd         := bytes.NewReader( chc.data )
    chc.zrd, _  = zip.NewReader( rd, rd.Size() )

}

func( chc *Cache ) endRead() {

    // End Read
    chc.zrd = nil

}
