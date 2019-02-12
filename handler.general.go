package sprout

/*
 + HANDLER GENERAL
 *
 * This file includes general handlers
 * A handler tends to handle a single instance unlike the With functions of Space
 */

type handler_general struct {}
var  static_handler_general = &handler_general{}

func HandlerGeneral() *handler_general {}
func( _hdlrgen *handler_general ) Status( _code int ) Handler {}
func( _hdlrgen *handler_general ) Template( _template *template.Template, _data_func func() interface{} ) Handler {}
func( _hdlrgen *handler_general ) Asset( _asset *Asset ) Handler {}
func( _hdlrgen *handler_general ) File( _path string ) Handler {}
func( _hdlrgen *handler_general ) Directory( _path string ) Handler {}
func( _hdlrgen *handler_general ) ReverseProxy( _url string ) Handler {}
func( _hdlrgen *handler_general ) Text( _mime_type, _text string ) Handler {}