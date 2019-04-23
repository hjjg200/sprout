package network

import (
    "net/http"
    "net/http/httputil"
    "net/url"
    "strings"
    "regexp"

    "github.com/hjjg200/sprout/environ"
    "github.com/hjjg200/sprout/util/errors"
    "github.com/hjjg200/sprout/volume"
)

type Space struct {
    aliases  []string
    handlers []Handler
    volume   volume.Volume
}

func NewSpace() *Space {
    return &Space{
        aliases: make( []string, 0 ),
        handlers: make( []Handler, 0 ),
    }
}

func( spc *Space ) Aliases() []string {
    return spc.aliases
}

func( spc *Space ) SetAliases( aliases []string ) {
    spc.aliases = aliases
}

func( spc *Space ) AddAlias( alias string ) {
    spc.aliases = append( spc.aliases, alias )
}

func( spc *Space ) Handlers() []Handler {
    return spc.handlers
}

func( spc *Space ) SetHandlers( handlers []Handler ) {
    spc.handlers = handlers
}

func( spc *Space ) AddHandler( handler Handler ) {
    spc.handlers = append( spc.handlers, handler )
}

func( spc *Space ) Volume() volume.Volume {
    return spc.volume
}

func( spc *Space ) SetVolume( vol volume.Volume ) {
    spc.volume = vol
}

// VOLUME-RELATED

func( spc *Space ) AssetHandler( path string ) Handler {
    return func( req *Request ) bool {
        return HandlerFactory.Asset( spc.volume.Asset( path ) )( req )
    }
}

func( spc *Space ) TemplateHandler( path string, dataFunc func( *Request ) interface{} ) Handler {
    return func( req *Request ) bool {
        return HandlerFactory.Template( spc.volume.Template( path ), dataFunc )( req )
    }
}

// GENERAL

func( spc *Space ) ServeRequest( req *Request ) {

    if !spc.ContainsHost( req.body.Host ) {
        // Bad Request
        req.PopBlank( 400 )
        return
    }

    spc.serveRequest( req )

}

func( spc *Space ) serveRequest( req *Request ) {

    // Set parent
    req.space = spc

    // Check locale with the space's i18n
    if spc.volume != nil {
        if spc.volume.I18n() != nil {
            req.PopulateLocalizer( spc.volume.I18n() )
        }
    }

    // Handle
    for _, handler := range spc.handlers {
        if done := handler( req ); done {
            break
        }
    }

}

func( spc *Space ) ServeHTTP( w http.ResponseWriter, r *http.Request ) {
    spc.ServeRequest( NewRequest( w, r ) )
}

func( spc *Space ) ContainsHost( host string ) bool {

    // Empty aliases
    if len( spc.aliases ) == 0 {
        return true
    }

    // Split
    split := strings.SplitN( host, ":", 2 )

    // Compare
    for _, val := range spc.aliases {
        if val == split[0] {
            return true
        }
    }

    return false

}

// With funcs

func( spc *Space ) WithHandler( hnd Handler ) {
    spc.handlers = append( spc.handlers, hnd )
}

func( spc *Space ) WithReverseProxy( target string ) {

    urlObj, err := url.Parse( target )
    proxy       := httputil.NewSingleHostReverseProxy( urlObj )

    spc.WithHandler( func( req *Request ) bool {

        if err != nil {
            req.PopError( 502, errors.ErrReverseProxy.Append( target ) )
        }

        // Log
        environ.Logger.OKln(
            "ID",
            req.ID(),
            req.String(),
            "Reverse Proxy",
        )

        // Update headers to allow SSL connection
        req.body.URL.Host = urlObj.Host
        req.body.URL.Scheme = urlObj.Scheme
        req.body.Header.Set( "X-Forwarded-Host", req.body.Header.Get( "Host" ) )
        req.body.Host = urlObj.Host

        proxy.ServeHTTP( req.writer, req.body )
        return true

    } )

}

func( spc *Space ) WithSymlink( targetPath, linkPath string ) {}
func( spc *Space ) WithRoute( rgxStr string, methods []string, hnd Handler ) {

    rgx, err := regexp.Compile( rgxStr )
    checker  := MakeMethodChecker( methods )

    spc.WithHandler( func( req *Request ) bool {

        if rgx != nil && err == nil && checker[req.body.Method] {

            matches := rgx.FindStringSubmatch( req.body.URL.Path )
            if len( matches ) >= 1 {
                req.vars = matches
                return hnd( req )
            }

        }
        return false

    } )

}
func( spc *Space ) WithAssetServer( prefix string ) {

    spc.WithHandler( func( req *Request ) bool {
        path := req.body.URL.Path
        if strings.HasPrefix( path, prefix ) && len( path ) > len( prefix ) {
            astPath := "asset/" + path[len( prefix ):]
            req.PopAsset( spc.volume.Asset( astPath ) )
            return true
        }
        return false
    } )

}
