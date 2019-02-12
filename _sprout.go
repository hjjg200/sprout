package sprout

import (
    "bytes"
    "crypto/sha256"
    "encoding/json"
    "errors"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "strconv"
    "html/template"
    "time"

    "./log"
//  "./session"
)

/*
 + Private Variables
 */

const (
    envAppName = "sprout"
    envVersion = "pre-alpha 0.6"

    // Directory names must not contain slashes, dots, etc.
    envDirAsset    = "asset"
    envDirCache    = "cache"
    envDirLocale   = "locale"
    envDirTemplate = "template"
)

var (
    envOS string
)

/*
 + Public Variables
 */

var (
    ErrNotSupportedOS = errors.New( "sprout: the OS is not supported" )
    ErrDirectory      = errors.New( "sprout: could not access a necessary directory" )
    ErrInvalidDirPath = errors.New( "sprout: the given path is invalid" )
)

var (
    EnvFilenameTimeFormat = "20060102-150405"
)

type asset struct {
    modTime time.Time
    data    []byte

    // sha256 hash of the file
    hash    string
}

func makeAsset( mt time.Time, r io.Reader ) asset {
    h   := sha256.New()
    mts := strconv.FormatInt( mt.Unix(), 10 )
    h.Write( []byte( mts ) )
    buf := bytes.NewBuffer( nil )
    io.Copy( buf, r )
    return asset{
        modTime: mt,
        data: buf.Bytes(),
        hash: fmt.Sprintf( "%x", h.Sum( nil ) ),
    }
}

type Sprout struct {
    cwd       string
    assets    map[string] asset
    templates *template.Template
    servers   map[string] *Server
    localizer *localizer
    default_locale string

    whitelistedExtensions []string
}

func New() *Sprout {

    s      := &Sprout{}
    // Cwd
    _cwd, _ := filepath.Abs( "./" )
    s.cwd    = filepath.ToSlash( filepath.Clean( _cwd ) )

    sanityCheck()

    // Assign the default whitelisted extensions
    s.whitelistedExtensions = make( []string, len( defaultWhitelistedExtensions ) )
    copy( s.whitelistedExtensions, defaultWhitelistedExtensions )

    log.Infoln( "Preparing a new Sprout instance..." )

    s.assets  = make( map[string] asset )
    s.servers = make( map[string] *Server )

    prod, _ := s.NewServer( "production" )
    prod.Mux().WithCachedAssetServer()

    log.Infoln( "Loacting the latest cache..." )
    // Load the Latest Cache
    // Build Cache If None Found
    _lcn, _err := s.LatestCacheName()
    if _err != nil {
        log.Infoln( "Could not load the latest cache, attempting to build one..." )
        _lcn, _err = s.BuildCache()
        if _err != nil { log.Severeln( _err ) }
        log.Infoln( "Successfully built a cache:", _lcn )
    }

    _err = s.LoadCache( _lcn )
    if _err != nil { log.Severeln( _err ) }
    log.Infoln( "Loaded Cache:", _lcn )

    dev, _  := s.NewServer( "dev" )
    dev.Mux().WithRealtimeAssetServer()

    // Localizer
    log.Infoln( "Attempting to load locales" )
    if len( s.localizer.locales ) > 0 {
        for _k, _ := range s.localizer.locales {
            s.SetDefaultLocale( _k )
            break
        }
        log.Infoln( "Successfully loaded", len( s.localizer.locales ), "locales" )
    } else {
        log.Infoln( "No locale file has been found" )
    }

    return s

}

func ( sp *Sprout ) SetDefaultLocale( locale string ) error {
    _, _ok := sp.localizer.locales[locale]
    if !_ok {
        log.Infoln( locale, "is not a valid locale. The default locale has not been changed." )
        return ErrInvalidLocale
    }
    log.Infoln( "Changed the default locale to", locale )
    sp.default_locale = locale
    return nil
}

func sanityCheck() error {
    // check if there is any sass, scss if so check sass installed
    log.Infoln( "Doing a sanity check..." )
    if err := checkOS(); err != nil {
        log.Severeln( "Sanity check failed!" )
    }
    if err := ensureDirectories(); err != nil {
        log.Severeln( "Sanity check failed!" )
    }
    log.Infoln( "Everything looks fine!" )
    return nil
}

func ensureDirectories() error {
    log.Infoln( "Ensuring all the necessary directories..." )
    _dirs_to_ensure := []string{
        envDirAsset, envDirCache, envDirLocale, envDirTemplate,
    }
    for _, _dir := range _dirs_to_ensure {
        _err := ensureDirectory( _dir )
        if _err != nil {
            log.Warnln( "Could not ensure all the directories!" )
            return _err
        }
    }
    log.Infoln( "Ensured all the directories!" )
    return nil
}

func ensureDirectory( p string ) error {
    log.Infoln( "Ensuring a directory...", p )
    st, err := os.Stat( p )
    switch {
    case os.IsNotExist( err ):
        err = os.Mkdir( p, 0750 )
        if err != nil {
            log.Warnln( "Error during ensuring the directory:", p, err )
            return err
        }
    case err != nil:
        log.Warnln( "Error during ensuring the directory:", p, err )
        return err
    case !st.IsDir():
        log.Warnln( "Error during ensuring the directory:", p, ErrDirectory )
        return ErrDirectory
    }
    log.Infoln( "Directory ready to go", p )
    return nil
}
