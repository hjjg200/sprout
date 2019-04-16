package network

import (
    "encoding/json"
    "html/template"
    "io"
    "net/http"
    "time"

    "../environ"
    "../util"
)

type Responder struct {
    req    *Request
    writer *responseWriter
}

/*
func( req *Request ) Responder( code int ) *Responder
    this can be found in request.go
*/

func( rsp *Responder ) Status() int {
    return rsp.writer.status
}

func( rsp *Responder ) Header() http.Header {
    return rsp.writer.Header()
}

func( rsp *Responder ) Content( contentType, text string ) {
    rsp.Header().Set( "content-type", contentType )
    rsp.writer.Write( []byte( text ) )
}

func( rsp *Responder ) Html( html string ) {
    rsp.Content( "text/html;charset=utf-8", html )
}

func( rsp *Responder ) Text( text string ) {
    rsp.Content( "text/plain;charset=utf-8", text )
}

func( rsp *Responder ) json( obj interface{}, pretty bool ) {

    // Json
    var (
        p []byte
    )

    // Marshal
    if pretty {
        p, _ = json.MarshalIndent( obj, "", "  " )
    } else {
        p, _ = json.Marshal( obj )
    }

    rsp.Content( "text/json;charset=utf-8", string( p ) )

}

func( rsp *Responder ) Json( obj interface{} ) {
    rsp.json( obj, false )
}

func( rsp *Responder ) PrettyJson( obj interface{} ) {
    rsp.json( obj, true )
}

func( rsp *Responder ) xml( obj interface{}, pretty bool ) {}

func( rsp *Responder ) Xml( obj interface{} ) {}

func( rsp *Responder ) PrettyXml( obj interface{} ) {}

func( rsp *Responder ) File( name string, modTime time.Time, rdskr io.ReadSeeker ) {

    // Serve
    http.ServeContent( rsp.writer, rsp.req.body, name, modTime, rdskr )

}

func( rsp *Responder ) Attachment( name string, modTime time.Time, rdskr io.ReadSeeker ) {

    // Set it to octet stream so that it won't be executed or compiled
    rsp.Header().Set( "Content-Type", "application/octet-stream" )
    rsp.Header().Set( "Content-Transfer-Encoding", "Binary" )
    rsp.Header().Set( "Content-Disposition", "attachment; filename=\"" + name + "\"" )
    rsp.File( name, modTime, rdskr )

}

func( rsp *Responder ) Error( err error ) {

    var (
        tmpl *template.Template
        msg = util.HttpStatusMessages[rsp.Status()]
    )

    if rsp.req.space.volume != nil {
        tmpl = rsp.req.space.volume.Template( environ.ErrorPageTemplatePath )
    }
    if tmpl == nil {
        tmpl = environ.DefaultErrorPageTemplate
    }

    // Map
    var data map[string] interface{}

    if err == nil {
        data = map[string] interface{} {
            "status": rsp.Status(),
            "message": msg,
        }
    } else {
        data = map[string] interface{} {
            "status": rsp.Status(),
            "message": msg,
            "error": err.Error(),
        }
    }

    HandlerFactory.Template( tmpl, func( req *Request ) interface{} {
        return data
    } )( rsp.req )

}

func( rsp *Responder ) Blank() {
    rsp.Error( nil )
}

func( rsp *Responder ) Redirect( url string ) {}

func( rsp *Responder ) Template( *template.Template, data interface{} ) {}

func( rsp *Responder ) Asset( *volume.Asset ) {}

