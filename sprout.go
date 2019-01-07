package sprout

import (
    "net/http"
)

var (
    envOS string
)

type route struct {
    rgx *regexp.Regexp
    hh  http.HandlerFunc
}

type Sprout struct {
    assets map[string] *bytes.Buffer
    routes []route
    mux           *http.ServeMux
    srvProduction *http.Server
    srvDev        *http.Server
}

func New() *Sprout {

    s := &sprout{}

    s.assets = make( map[string] *bytes.Buffer )
    s.routes = make( []route, 0 )

    s.mux    = http.NewServeMux()
    s.mux.Hanlde( "/", serveHTTP )

    err := s.sanityCheck()
    if err != nil {
        log.Fatalln( err )
    }

}

func ( s *Sprout ) sanityCheck() error {
    // check if there is any sass, scss if so check sass installed
    if err := checkOS(); err != nil {
        return err
    }
}

func ( s *Sprout ) checkOS() error {
    var (
        ErrNotSupportedOS = errors.New( "sprout: the OS is not supported" )
    )
    switch runtime.GOOS {
    case "windows", "linux", "darwin", "freebsd", "openbsd":
        envOS = runtime.GOOS
        return nil
    }
    return ErrNotSupportedOS
}

func ( s *Sprout ) doesCommandExist( cmd string ) bool {

    var (
        err error
        out bytes.Buffer
        e   *exec.Cmd
    )

    switch envOS {
    case "linux", "darwin", "freebsd", "openbsd":
        s := "if command -v " + cmd + " > /dev/null 2>&1; then echo 'true'; fi"
        e  = exec.Command( "bash", "-c", s )
    case "windows":
        s := "where /Q " + cmd + " & if %errorlevel%==0 echo true"
        e  = exec.Command( "cmd", "/C", s )
    }

    e.Stdout = &out
    err = e.Run()

    if err != nil {
        panic( err )
        return false
    }

    r := out.String()
    if r[:4] == "true" {
        return true
    }

    return false
}

func ( s *Sprout ) runCommand( cmd string ) error {
    var e *exec.Cmd
    switch envOS {
    case "linux", "darwin", "freebsd", "openbsd":
        e = exec.Command( "bash", "-c", cmd )
    case "windows":
        e = exec.Command( "cmd", "/C", cmd )
    }
    return e.Run()
}

func ( s *Sprout ) ProcessAsset( p string ) error {

    var (
        ErrInvalidAsset = errors.New( "sprout: the asset is invalid" )
    )

    if !strings.HasPrefix( p, "asset/" ) &&
       !strings.HasPrefix( p, "./asset/" ) {
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
    dir := filepath.Dir( p )
    ext := filepath.Ext( p )
    switch ext {
    case "sass", "scss":
        css := dir + "/" + bs[:len( bs ) - 4] + "css"
        cmd := fmt.Sprintf( "sass %s %s", dir + "/" + bs, css )
        err := s.runCommand( cmd )
        if err != nil {
            return err
        }
    }

}

func ( s *Sprout ) WriteAsset( w io.Writer, path string ) error {
    _, err := io.Copy( w, s.assets[path] )
    return err
}

func ( s *Sprout ) WriteJSON( w io.Writer, v interface{} ) error {
    return json.NewEncoder( w ).Encode( v )
}

func ( s *Sprout ) BuildCache() error {

    t       := time.Now().Format( "20060102-150405" )
    fn      := "./cache/" + t + ".tmp"
    f, err  := os.OpenFile( fn, os.O_WRONLY, 0600 )
    if err != nil {
        return err
    }
    defer f.Close()

    zw := zip.NewWriter( f )
    defer zw.Close()

    foreach := func ( dir string, do func ( string ) error ) error {
        dir        = filepath.Clean( dir )
        fis, err2 := ioutil.ReadDir( dir )
        if err2 != nil {
            return err2
        }

        for _, fi := range fis {
            path := dir + "/" + fi.Name()
            if fi.IsDir() {
                err2 = do( path )
                if err2 != nil {
                    return err2
                }
                return foreach( path, do )
            } else {
                err2 = do( path )
                if err2 != nil {
                    return err2
                }
            }
        }
        return nil
    }

    archive := func ( path string ) error {
        w, err3 := zw.Create( path )
        if err3 != nil {
            return err3
        }
        st, err3 := os.Stat( path )
        if err3 != nil {
            return err3
        }
        if st.IsDir() {
            return nil
        }
        pw, err3 := os.Open( path )
        if err3 != nil {
            return err3
        }
        defer pw.Close()

        // Assign to assets
        s.assets[path] = &bytes.Buffer{}
        _, err3 = io.Copy( s.assets[path], pw )
        if err3 != nil {
            return err3
        }

        // Write to zip
        _, err3 = io.Copy( w, pw )
        if err3 != nil {
            return err3
        }
    }

    //

    err = foreach( "asset/", ProcessAsset )
    if err != nil {
        return err
    }
    err = foreach( "asset/" archive )
    if err != nil {
        return err
    }

    h  := sha256.New()
    io.Copy( h, f )
    hs := fmt.Sprintf( "%x", h.Sum( nil ) )

    err = os.Rename( fn, "./cache/" + t + "-" + hs[:6] + ".zip" )
    if err != nil {
        return err
    }

    return nil

}

func ( s *Sprout ) serveHTTP( w http.ResponseWriter, r *http.Request ) {

    for _, route := s.routes {
        rgx := route.rgx
        hh  := route.hh

        if rgx.MatchString( r.URL.Path ) {
            hh( w, r )
            return
        }
    }

}

func ( s *Sprout ) AddRoute( rgxStr string, hh http.HandlerFunc ) error {

    rgx, err := regexp.Compile( rgxStr )
    if err != nil {
        return err
    }

    s.routes = append( s.routes, route{
        rgx: rgx,
        hh: hh,
    } )

}

func ( s *Sprout ) StartServer( addr string ) error {
    s.srvProduction = &http.Server{
        Addr: addr,
        Handler: s.mux,
    }
    return s.srvProduction.ListenAndServe()
}

func ( s *Sprout ) StartDevServer( addr string ) error {
    s.srvDev = &http.Server{
        Addr: addr,
        Handler: s.mux,
    }
    return s.srvDev.ListenAndServe()
}