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
    "sort"
    "strings"
    "html/template"
    "time"

    "github.com/hjjg200/go/together"
)

/*
 + fs.go
 |
 + Note that this file uses filepath for path operations while the others use path, as this file deals with the underlying filesystem.
 */

var fileHoldGroup = together.NewHoldGroup()

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

    if len( fis ) == 0 {
        return "", ErrNoAvailableCache
    }

    cn := make( cacheNames, len( fis ) )
    for i := range fis {
        cn[i] = fis[i].Name()
    }

    sort.Sort( sort.Reverse( cn ) )
    return cn[0], nil

}

func ( s *Sprout ) LoadCache( fn string ) error {

    zr, _err := zip.OpenReader( envDirCache + "/" + fn )
    if _err != nil {
        return _err
    }
    defer zr.Close()

    // Empty caches
    s.assets = make( map[string] asset )
    s.templates = template.New( "" )
    s.templates.Delims( template_left_delimiter, template_right_delimiter )
    s.localizer = s.newLocalizer()

    // Assign files
    for _, f := range zr.File {

        fn   := f.Name
        _ext := filepath.Ext( fn )
        // Continue if it is a directory
        if strings.HasSuffix( f.Name, "/" ) {
            continue
        }

        frc, _err := f.Open()
        if _err != nil {
            panic( _err )
            continue
        }

        switch {
        case strings.HasPrefix( fn, envDirAsset ):
            s.assets[fn] = makeAsset(
                f.Modified, frc,
            )
        case strings.HasPrefix( fn, envDirLocale ):
            if _ext != ".json" {
                continue
            }
            _buf := bytes.Buffer{}
            _buf.ReadFrom( frc )
            _base   := filepath.Base( fn )
            _ext    := filepath.Ext( _base )
            _locale := _base[:len( _base ) - len( _ext )]
            _err     = s.localizer.appendLocale( _locale, _buf.Bytes() )
            if _err != nil {
                panic( _err )
            }
        case strings.HasPrefix( fn, envDirTemplate ):
            if !string_slice_includes( template_extensions, _ext ) {
                continue
            }
            _buf := bytes.Buffer{}
            _buf.ReadFrom( frc )
            _, _err = s.templates.New( fn ).Parse( _buf.String() )
            if _err != nil {
                panic( _err )
            }
        }

        frc.Close()

    }

    return nil

}

func ( s *Sprout ) BuildCache() ( string, error ) {

    /*
     | Prepare the Zip File
     */

    t       := time.Now().Format( EnvFilenameTimeFormat )
    fn      := envDirCache + "/" + t + ".tmp"
    f, err  := os.OpenFile( fn, os.O_RDWR | os.O_CREATE, 0600 )
    if err != nil {
        return "", err
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
        fh, err3 := zip.FileInfoHeader( st )
        fh.Name = path
        if err3 != nil {
            return err3
        }
        w, err3 := zw.CreateHeader( fh )
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
        /*
        f, err := os.Open( path )
        if err != nil {
            return err
        }

        s.assets[path] = makeAsset(
            st.ModTime(), f,
        )

        f.Close()*/

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
        return "", err
    }

    _dirs_to_cache := []string{
        envDirAsset, envDirTemplate, envDirLocale,
    }
    for _, _dir := range _dirs_to_cache {
        err = foreach( _dir, archive )
        if err != nil {
            return "", err
        }
    }

    /*
     | End Archiving
     */

    zw.Close()
    _, err = f.Seek( 0, os.SEEK_SET )
    if err != nil {
        return "", err
    }

    /*
     | Change Filename
     */

    h  := sha256.New()
    io.Copy( h, f )
    hs := fmt.Sprintf( "%x", h.Sum( nil ) )
    f.Close()

    cn := formatCacheName( t, hs )
    err = os.Rename( fn, envDirCache + "/" + cn )
    if err != nil {
        return "", err
    }

    return cn, nil

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

    // If the file exists

        // Lock mutex so that a file won't be processed multiple times at the same time
    fileHoldGroup.HoldAt( p )
    defer fileHoldGroup.UnholdAt( p )

    bs  := filepath.Base( p )
    dir := filepath.ToSlash( filepath.Dir( p ) )
    ext := strings.ToLower( filepath.Ext( p ) )

    switch ext {
    case ".sass", ".scss":
        css := dir + "/" + bs[:len( bs ) - 5] + ".css"
        cmd := fmt.Sprintf( "sass %s %s", dir + "/" + bs, css )
        err := s.runCommand( cmd )
        if err != nil {
            return err
        }
    case ".css":

        // If it is a css file, look for a sass file

        pb   := dir + "/" + bs[:len( bs ) - 4]
        scss := pb + ".scss"
        sass := pb + ".sass"
        // Look for scss
        st, err = os.Stat( scss )
        if err == nil && !st.IsDir() {
            cmd := fmt.Sprintf( "sass %s %s", scss, pb + ".css" )
            err  = s.runCommand( cmd )
            if err != nil {
                return err
            }
            return nil
        }
        // Look for sass
        st, err = os.Stat( sass )
        if err == nil && !st.IsDir() {
            cmd := fmt.Sprintf( "sass %s %s", sass, pb + ".css" )
            err  = s.runCommand( cmd )
            if err != nil {
                return err
            }
            return nil
        }

        // Return nil since a css doesn't have to have a matching sass file
        return nil

    }

    return nil

}