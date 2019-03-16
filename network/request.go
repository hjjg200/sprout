package network

type Request struct {
    space  *Space
    body   *http.Request
    writer http.ResponseWriter
    locale string
}

func NewRequest( parent *Space, w http.ResponseWriter, r *http.Request ) *Request {
    
    // New
    req := &Request{
        space: parent,
        body: r,
        writer: w,
    }
    
    // Check locale
    i1 := parent.volume.i18n // shorthand
    switch i1.NumLocale() {
    case 0:
    default:
        lcName, err := i1.ParseUrl( r.URL )
        if err != nil {
            lcName, err = i1.ParseCookies( r.Cookies )
            if err != nil {
                lcName, err = i1.ParseAcceptLanguage( r.Header.Get( "accept-language" ) )
                if err != nil {
                    lcName = i1.DefaultLocale()
                }
            }
        }
        req.locale = lcName
    }
    
    // return
    return req
    
}
