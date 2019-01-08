package sprout

import (
    "archive/zip"
    "bytes"
    "crypto/sha256"
    "errors"
    "fmt"
    "io"
    "io/ioutil"
    "os"
    "path/filepath"
    "strings"
    "time"
)

/*
 + fs.go
 |
 + Note that this file uses filepath for path operations while the others use path, as this file deals with the underlying filesystem.
 */

func ( s *Sprout ) LoadCache( fn string ) error {

    return nil
}

func ( s *Sprout ) BuildCache() error {

    /*
     | Prepare the Zip File
     */

    t       := time.Now().Format( EnvFilenameTimeFormat )
    fn      := envDirCache + "/" + t + ".tmp"
    f, err  := os.OpenFile( fn, os.O_WRONLY | os.O_CREATE, 0600 )
    if err != nil {
        return err
    }

    zw := zip.NewWriter( f )

    /*
     | Helper Methods
     */

    var foreach func ( string, func ( string ) error ) error
    foreach = func ( dir string, do func ( string ) error ) error {
        dir        = filepath.ToSlash( filepath.Clean( dir ) )
        fis, err2 := ioutil.ReadDir( dir )
        if err2 != nil {
            return err2
        }

        for _, fi := range fis {
            path := dir + "/" + fi.Name()
            err2 = do( path )
            if err2 != nil {
                return err2
            }
            if fi.IsDir() {
                err2 = foreach( path, do )
                if err2 != nil {
                    return err2
                }
            }
        }
        return nil
    }

    archive := func ( path string ) error {

        st, err3 := os.Stat( path )
        if err3 != nil {
            return err3
        }

        // Add Slash at the End If path Resolves to a Folder
        if st.IsDir() {
            path    = path + "/"
            _, err3 = zw.Create( path )
            if err3 != nil {
                return err3
            }
            return nil
        }

        // Create Node in the Zip
        w, err3 := zw.Create( path )
        if err3 != nil {
            return err3
        }

        pw, err3 := os.Open( path )
        if err3 != nil {
            return err3
        }
        defer pw.Close()

        // Assign to assets

        /*
        | Assign Asset to s.assets
        */

        dat, err := ioutil.ReadFile( path )
        if err != nil {
            return err
        }

        r := bytes.NewReader( dat )

        h := sha256.New()
        io.Copy( h, r )

        s.assets[path] = asset{
            modTime: st.ModTime(),
            reader: r,
            hash: fmt.Sprintf( "%x", h.Sum( nil ) ),
        }

        // Write to zip
        _, err3 = io.Copy( w, pw )
        if err3 != nil {
            return err3
        }

        return nil
    }

    /*
     | Archive Files
     */

    err = foreach( envDirAsset, s.ProcessAsset )
    if err != nil {
        return err
    }
    err = foreach( envDirAsset, archive )
    if err != nil {
        return err
    }

    /*
     | Change the Filename
     */

    h  := sha256.New()
    io.Copy( h, f )
    hs := fmt.Sprintf( "%x", h.Sum( nil ) )

    /*
     | End Archiving
     */

    zw.Close()
    f.Close()

    err = os.Rename( fn, envDirCache + "/" + t + "-" + hs[:6] + ".zip" )
    if err != nil {
        return err
    }

    return nil

}

func ( s *Sprout ) ProcessAsset( p string ) error {

    var (
        ErrInvalidAsset = errors.New( "sprout: the asset is invalid" )
    )

    p = filepath.ToSlash( filepath.Clean( p ) )

    if !strings.HasPrefix( p, envDirAsset + "/" ) {
        return ErrInvalidAsset
    }

    st, err := os.Stat( p )
    if err != nil {
        return err
    }
    if st.IsDir() {
        return nil
    }

    bs  := filepath.Base( p )
    dir := filepath.ToSlash( filepath.Dir( p ) )
    ext := filepath.Ext( p )

    switch ext {
    case ".sass", ".scss":
        css := dir + "/" + bs[:len( bs ) - 4] + "css"
        cmd := fmt.Sprintf( "sass %s %s", dir + "/" + bs, css )
        err := s.runCommand( cmd )
        if err != nil {
            return err
        }
    }

    return nil

}