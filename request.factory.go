package sprout

/*
 + REQUEST FACTORY
 */

type request_factory struct {}
var  static_request_factory = &request_factory{}

func RequestFactory() *request_factory {
    return *static_request_factory
}
func( _reqfac *request_factory ) New( _w http.ResponseWrtier, _r *http.Request ) *Request {
    return &Request{
        writer: _w,
        body: _r,
    }
}