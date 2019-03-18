package network

type Space struct {
    name     string // domain
    aliases  []string
    handlers []Handler
    volume   *Volume
}

func NewSpace( name string, aliases []string ) *Space {
    return &Space{
        name: name,
        aliases: aliases,
    }
}

func( spc *Space ) ServeHTTP( w http.ResponseWriter, r *http.Request ) {

    req := NewRequest( w, r )
    if !spc.ContainsAlias( req.body.URL.Host ) {
        // Bad Request
        HandlerFactory.Status( 400 )( req )
        return
    }

    // Handle
    for _, handler := range spc.handlers {
        if handler( req ) {
            break
        }
    }

}

func( spc *Space ) AddHandler( handler *Handler ) {
    spc.handlers = append( spc.handlers, handler )
}

func( spc *Space ) AddAlias( alias string ) {
    spc.aliases = append( spc.aliases, alias )
}

func( spc *Space ) AssignVolume( vol *Volume ) {
    spc.volume = vol
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