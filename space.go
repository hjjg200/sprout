package sprout

/*
 + SPACE
 *
 * A space handles a single website and it can have a volume and own routing rules
 * The With functions of Space tend to handle and manage multiple instances, compared to general handlers
 */

type Space struct {
    name     string // domain
    aliases  []string
    handlers []Handler
    volume   *Volume
}

func( _space *Space ) Name() string {
    return _space.name
}
func( _space *Space ) SetName( _name string ) {
    _space.name = _name
}
func( _space *Space ) Aliases() []string {
    _copy := make( []string, len( _space.aliases ) )
    copy( _copy, _space.aliases )
    return _copy
}
func( _space *Space ) SetAliases( _aliases []string ) {
    _copy := make( []string, len( _aliases ) )
    copy( _copy, _aliases )
    _space.aliases = _copy
}
func( _space *Space ) AddAlias( _alias string ) {
    _space.aliases = append( _space.aliases, _alias )
}
func( _space *Space ) Handlers() []Handler {
    _copy := make( []Handler, len( _space.handlers ) )
    copy( _copy, _space.handlers )
    return _copy
}
func( _space *Space ) SetHandlers( _handlers []Handler ) {
    _copy := make( []Handler, len( _handlers ) )
    copy( _copy, _hanlders )
    _space.handlers = _copy
}
func( _space *Space ) WithHandler( _handler Handler ) {
    _space.handlers = append( _space.handlers, _handler )
}
func( _space *Space ) WithReverseProxy( _url string ) {}
func( _space *Space ) WithSymlink( _link, _path string ) {}
func( _space *Space ) WithRoute( _regexp string, _handler Handler ) {}
func( _space *Space ) WithAssetServer( _base string ) {}
func( _space *Space ) WithAuthenticator( _authenticator Authenticator ) {}