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
    // Directory names must not contain slashes, dots, etc.
    envDirAsset    = "asset"
    envDirCache    = "cache"
    envDirLocale   = "locale"
    envDirTemplate = "template"
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