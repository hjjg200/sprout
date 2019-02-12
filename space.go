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

func( _space *Space ) Name() string {}
func( _space *Space ) SetName( _name string ) {}
func( _space *Space ) Aliases() []string {}
func( _space *Space ) SetAliases( _aliases []string ) {}
func( _space *Space ) AddAlias( _alias string ) {}
func( _space *Space ) Handlers() []Handler {}
func( _space *Space ) SetHandlers( _handlers []Handler ) {}
func( _space *Space ) WithHandler( _handler Handler ) {}
func( _space *Space ) WithSymlink( _link, _path string ) {}
func( _space *Space ) WithRoute( _regexp string, _handler Handler ) {}
func( _space *Space ) WithAssetServer( _base string ) {}
func( _space *Space ) WithAuthenticator( _authenticator Authenticator ) {}