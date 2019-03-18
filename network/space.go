package network

type Space struct {
    name     string // domain
    aliases  []string
    handlers []Handler
    volume   *Volume
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

func( spc *Space ) AddHandler( handler *Handler ) {
    spc.handlers = append( spc.handlers, handler )
}

func( spc *Space ) Volume() *Volume {
    return spc.volume
}

func( spc *Space ) SetVolume( vol *Volume ) {
    spc.volume = vol
}

func( spc *Space ) ServeRequest( req *Request ) {

    if !spc.ContainsAlias( req.body.URL.Host ) {
        // Bad Request
        HandlerFactory.Status( 400 )( req )
        return
    }

    // Check locale with the space's i18n
    req.CheckLocale( spc.volume.i18n )

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