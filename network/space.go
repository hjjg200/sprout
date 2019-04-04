package network

import (
    "net/http"
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

func( spc *Space ) ServeRequest( req *Request ) {

    if !spc.ContainsAlias( req.body.URL.Host ) {
        // Bad Request
        HandlerFactory.Status( 400 )( req )
        return
    }

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

    // Compare
    if spc.name == alias {
        return true
    }
    for _, val := range spc.aliases {
        if val == alias {
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
            if rgx.MatchString( req.body.URL.Path ) {
                hnd( req )
                return true
            }
        }
        return false
    } )

}
func( spc *Space ) WithAssetServer() {}

func( spc *Space ) WithAsset( path string ) {
    spc.WithHandler( func( req* Request ) bool {
        ast, ok := spc.vol.Asset( path )
        if !ok {
            HandlerFactory.Status( 404 )( req )
            return true
        }
        HandlerFactory.Asset( ast )( req )
        return true
    } )
}
func( spc *Space ) WithTemplate( path string, dataFunc func( *Request ) interface{} ) {
    spc.WithHandler( func( req *Request ) bool {
        tmpl, ok := spc.vol.Template( path )
        if !ok {
            HandlerFactory.Status( 404 )( req )
            return true
        }
        HandlerFactory.Template( tmpl, dataFunc )( req )
        return true
    } )
}

func( spc *Space ) WithAuthenticator( auther func( *Request ) bool ) {}