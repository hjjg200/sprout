package sprout

/*
 + HANDLER FACTORY
 *
 * A handler factory can create various general handlers
 * A handler tends to handle a single instance unlike the With functions of Space
 */

type handler_factory struct {}
var  static_handlerbb_factory = &handler_factory{}

func HandlerFactory() *handler_factory {
    return static_handler_factory
}
func( _hdlrgen *handler_factory ) Status( _code int ) Handler {
    return func( _request *Request ) bool {
        _request.WriteStatus( _code )
        return true
    }
}
func( _hdlrgen *handler_factory ) Authenticator( _auther *Authenticator, _realm string ) Handler {
    return func( _request *Request ) bool {
        if _auther( _request ) == true {
            // Returns false so that the following handlers can handle the request
            return false
        }
        // Set the authentication realm
        _request.writer.Header().Set( "WWW-Authenticate", "Basic realm=\"" + realm + "\"" )
        _request.WriteStatus( 401 )
        // Returns true to notify that the handling should stop here
        return true
    }
}
func( _hdlrgen *handler_factory ) Template( _template *template.Template, _data_func func() interface{} ) Handler {}
func( _hdlrgen *handler_factory ) Asset( _asset *Asset ) Handler {}
func( _hdlrgen *handler_factory ) File( _path string ) Handler {}
func( _hdlrgen *handler_factory ) Directory( _path string ) Handler {}
func( _hdlrgen *handler_factory ) Text( _text, _mime_type string ) Handler {
    return func( _request *Request ) bool {
        _request.writer.Header().Set( "Content-Type", _mime_type + "; charset=utf-8" )
        _request.writer.Write( []byte( _text ) )
        return true
    }
}