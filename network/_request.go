package network

type Request struct {
    id          int64
    body        *http.Request
    writer      http.ResponseWriter
    status      int
    wroteHeader bool
    localizer   *i18n.Localizer
    space       *Space
    rsp         *Responder
    vars        []string
}

var (
    lastRequestID int64 = -1
)

func NewRequest( w http.ResponseWriter, r *http.Request ) *Request {
    return &Request{
        id: nextRequestID(),
        body: r,
        writer: w,
    }
}

func nextRequestID() int64 {
    lastRequestID++
    return lastRequestID
}

func( req *Request ) ID() int64 {
    return req.id
}

func( req *Request ) Body() *http.Request {
    return req.body
}

func( req *Request ) Localizer() *i18n.Localizer {
    return req.localizer
}

func( req *Request ) Status() int {
    return req.status
}

func( req *Request ) Vars() []string {
    return req.vars
}

func( req *Request ) Write( p []byte ) ( int, error ) {
    return req.writer.Write( p )
}

func( req *Request ) SetStatus( status int ) {

    if req.wroteHeader {
        environ.Logger.Panicln( ErrDifferentStatusCode.Append( "ID", req.ID() ) )
    }

    req.wroteHeader = true
    req.status      = status
    req.writer.WriteHeader( status )

    // Args
    args := []interface{}{
        "ID",
        req.ID(),
        req.String(),
        status,
    }

    // Log
    switch {
    case status >= 500:
        environ.Logger.Warnln( args... )
    default:
        environ.Logger.OKln( args... )
    }

}

func( req *Request ) WriteHeader( status int ) {
    req.SetStatus( status )
}

func( req *Request ) Header() http.Header {
    return req.writer.Header()
}

func( req *Request ) Responder() *Responder {

    if req.rsp != nil {
        return req.rsp
    }

    rsp := &Responder{ req }
    req.rsp = rsp
    return rsp

}

func( req *Request ) String() string {
    return fmt.Sprintf(
        "%s %s %s <= %s",
        req.body.Method,
        req.body.Host + req.body.URL.Path,
        req.body.Proto,
        req.body.RemoteAddr,
    )
}

func( req *Request ) PopulateLocalizer( i1 *i18n.I18n ) {

    // Check locale
    switch i1.NumLocale() {
    case 0:
    default:
        lcName, err := i1.ParseUrlPath( req.body.URL )
        if err != nil {
            lcName, err = i1.ParseUrlQuery( req.body.URL )
            if err != nil {
                lcName, err = i1.ParseCookies( req.body.Cookies() )
                if err != nil {
                    lcName, err = i1.ParseAcceptLanguage( req.body.Header.Get( "accept-language" ) )
                    if err != nil {
                        lcName = i1.DefaultLocale()
                    }
                }
            }
        } else { // When the URL path contains the locale at its beginning
            // Remove the locale from the URL for further processing
            path  := req.body.URL.Path[1:]
            split := strings.SplitN( path, "/", 2 )

            // Set path
            req.body.URL.Path = "/"
            if len( split ) == 2 {
                req.body.URL.Path += split[1]
            }
        }
        req.localizer = i1.Localizer( lcName )
        return

    }

    req.localizer = nil

}
