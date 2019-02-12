package sprout

/*
 + HANDLER
 *
 * A handler is a type that handles http requests and respond to it
 * A handler returns true if the handler handled the request and no further action is needed
 */

type Handler func( *Request ) bool

func( _handler *Handler ) ServeHTTP( _w http.ResponseWriter, _r *http.Request ) {} // interface http.Handler