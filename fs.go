package sprout

import (
    "archive/zip"
    "crypto/sha256"
    "errors"
    "fmt"
    "io"
    "io/ioutil"
    "os"
    "path/filepath"
    "sort"
    "strings"
    "time"
)

/*
 + fs.go
 |
 + Note that this file uses filepath for path operations while the others use path, as this file deals with the underlying filesystem.
 */

// For sorting cache names
type cacheNames []string

func ( cn cacheNames ) Len() int { return len( cn ) }
func ( cn cacheNames ) Swap( i, j int ) { cn[i], cn[j] = cn[j], cn[i] }
func ( cn cacheNames ) Less( i, j int ) bool {
    ti := parseTimeFromCacheName( cn[i] )
    tj := parseTimeFromCacheName( cn[j] )
    return ti.Sub( tj ) < 0
}

func formatCacheName( timeStr, hash string ) string {
    return timeStr + "-" + hash[:6] + ".zip"
}

func parseTimeFromCacheName( fn string ) time.Time {
    var (
        ErrInvalidCacheName = errors.New( "sprout: the given cache name is invalid" )
    )
    if len( fn ) < len( EnvFilenameTimeFormat ) {
        panic( ErrInvalidCacheName )
        return time.Now()
    }
    timeStr := fn[:len( EnvFilenameTimeFormat )]
    t, err  := time.ParseInLocation(
        EnvFilenameTimeFormat,
        timeStr,
        time.Now().Location(),
    )
    if err != nil {
        panic( err )
        return time.Now()
    }
    return t
}

func ( s *Sprout ) LatestCacheName() ( string, error ) {

    fis, err := ioutil.ReadDir( envDirCache )
    if err != nil {
        return "", err
    }

    cn := make( cacheNames, len( fis ) )
    for i := range fis {
        cn[i] = fis[i].Name()
    }

    sort.Sort( sort.Reverse( cn ) )
    return cn[0], nil

}

func ( s *Sprout ) LoadCache( fn string ) error {

    zr, err := zip.OpenReader( envDirCache + "/" + fn )
    if err != nil {
        return err
    }
    defer zr.Close()

    // Empty s.assets
    s.assets = make( map[string] asset )

    // Assign files
    for _, f := range zr.File {

        fn := f.Name
        // Continue if it is a directory
        if strings.HasSuffix( f.Name, "/" ) {
            continue
        }

        frc, err := f.Open()
        if err != nil {
            panic( err )
            continue
        }

        s.assets[fn] = makeAsset(
            f.Modified, frc,
        )

        frc.Close()

    }

    return nil

}

func ( s *Sprout ) BuildCache() error {

    /*
     | Prepare the Zip File
     */

    t       := time.Now().Format( EnvFilenameTimeFormat )
    fn      := envDirCache + "/" + t + ".tmp"
    f, err  := os.OpenFile( fn, os.O_RDWR | os.O_CREATE, 0600 )
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

        f, err := os.Open( path )
        if err != nil {
            return err
        }

        s.assets[path] = makeAsset(
            st.ModTime(), f,
        )

        f.Close()

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
     | End Archiving
     */

    zw.Close()
    _, err = f.Seek( 0, os.SEEK_SET )
    if err != nil {
        return err
    }

    /*
     | Change Filename
     */

    h  := sha256.New()
    io.Copy( h, f )
    hs := fmt.Sprintf( "%x", h.Sum( nil ) )
    f.Close()

    err = os.Rename( fn, envDirCache + "/" + formatCacheName( t, hs ) )
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