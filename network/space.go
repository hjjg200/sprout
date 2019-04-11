package network

import (
    "net/http"
    "strings"
    "regexp"

    "../volume"
)

type Space struct {
    name     string // domain
    aliases  []string
    handlers []Handler
    volume   volume.Volume
}

func NewSpace( name string ) *Space {
    return &Space{
        name: name,
        aliases: make( []string, 0 ),
    }
}

func( spc *Space ) Name() string {
    return spc.name
}

func( spc *Space ) SetName( name string ) {
    spc.name = name
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
    return HandlerFactory.Asset( spc.volume.Asset( path ) )
}

func( spc *Space ) TemplateHandler( path string, dataFunc func( *Request ) interface{} ) Handler {
    return HandlerFactory.Template( spc.volume.Template( path ), dataFunc )
}

// GENERAL

func( spc *Space ) ServeRequest( req *Request ) {

    if !spc.ContainsAlias( req.body.Host ) {
        // Bad Request
        HandlerFactory.Status( 400 )( req )
        return
    }

    spc.serveRequest( req )

}

func( spc *Space ) serveRequest( req *Request ) {

    // Check locale with the space's i18n
    if spc.volume != nil {
        if spc.volume.I18n() != nil {
            req.PopulateLocalizer( spc.volume.I18n() )
        }
    }

    // Handle
    for _, handler := range spc.handlers {
        if handler( req ) {
            break
        }
    }

}

func( spc *Space ) ServeHTTP( w http.ResponseWriter, r *http.Request ) {
    spc.ServeRequest( NewRequest( w, r ) )
}

func( spc *Space ) ContainsAlias( alias string ) bool {

    // Empty name
    if spc.name == "" && len( spc.aliases ) == 0 {
        return true
    }

    // Split
    split := strings.SplitN( alias, ":", 2 )

    // Compare
    if spc.name == split[0] {
        return true
    }
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

func( spc *Space ) WithReverseProxy( url string ) {}
func( spc *Space ) WithSymlink( targetPath, linkPath string ) {}
func( spc *Space ) WithRoute( rgxStr string, hnd Handler ) {

    rgx, err := regexp.Compile( rgxStr )
    spc.WithHandler( func( req *Request ) bool {
        if rgx != nil && err == nil {
            matches := rgx.FindStringSubmatch( req.body.URL.Path )
            if len( matches ) >= 1 {
                req.vars = matches
                hnd( req )
                return true
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
            return HandlerFactory.Asset( spc.volume.Asset( astPath ) )( req )
        }
        return false
    } )

}
