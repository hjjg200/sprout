package sprout

/*
 + SERVER HELPER
 *
 * ServerHelper includes helper methods for Server
 */

type server_helper struct{}
var  static_server_helper = &server_helper{}

func ServerHelper() *server_helper {
    return static_server_helper
}
func( _srvhlpr *server_helper ) ParseHostFromHTTPHost( _http_host string ) string {
    _index := strings.Index( _http_host, ":" )
    if _index != -1 {
        _http_host = _http_host[:_index + 1]
    }
    return _http_host
}
func( _srvhlpr *server_helper ) SpaceContainsAlias( _space *Space, _alias string ) bool {
    if _space.name == _alias {
        return true
    }
    for _, _value := range _space.aliases {
        if _value == _alias {
            return true
        }
    }
    return false
}