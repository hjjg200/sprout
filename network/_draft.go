package network

type Request chan interface{}
type requestInfo struct {
    id          int64
    body        *http.Request
    writer      http.ResponseWriter
    status      int
    wroteHeader bool
    localizer   *i18n.Localizer
    space       *Space
    vars        []string
}

var requestInfos map[chan interface{}] *requestInfo

func main() {

    req := NewRequest( w, r )
    req <- 200
    req <- Header{ "Content-Type", "text/plain" }
    req <- reader
    req <- text
    req <- someStruct // if xml tag is set as xml,
                      // if json as json
                      // if both, first one
                      // if none, DefaultObjectEncoding
    req <- []interface{}{ a, b, c } // print in order without spaces in between
    req <- someMap // as DefaultObjectEncoding

    {
        req <- Header{ "Content-Type",  }
        req <- 200
        req <- tmpl

    }

}