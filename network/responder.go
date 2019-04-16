package network

import (
    "bytes"
    "encoding/json"
    "html/template"
    "io"
    "net/http"
    "time"

    "../environ"
    "../util"
    "../volume"
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

func( rsp *Responder ) SetStatus( status int ) {
    rsp.writer.WriteHeader( status )
}

func( rsp *Responder ) Header() http.Header {
    return rsp.writer.Header()
}

func( rsp *Responder ) Content( ctype, text string ) {
    rsp.Header().Set( "content-type", ctype )
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
        err error
    )

    // Marshal
    if pretty {
        p, err = json.MarshalIndent( obj, "", "  " )
    } else {
        p, err = json.Marshal( obj )
    }

    // Error
    if err != nil {
        rsp.Error( 500, ErrMalformedJson.Append( err ) )
    }

    rsp.Content( "text/json;charset=utf-8", string( p ) )

}

func( rsp *Responder ) Json( obj interface{} ) {
    rsp.json( obj, false )
}

func( rsp *Responder ) PrettyJson( obj interface{} ) {
    rsp.json( obj, true )
}
/*
func( rsp *Responder ) xml( obj interface{}, pretty bool ) {}

func( rsp *Responder ) Xml( obj interface{} ) {}

func( rsp *Responder ) PrettyXml( obj interface{} ) {}
*/
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

func( rsp *Responder ) Error( status int, err error ) {

    var (
        tmpl *template.Template
        msg = util.HttpStatusMessages[rsp.Status()]
    )

    // Change
    rsp.SetStatus( status )

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

    rsp.Template( tmpl, data )

}

func( rsp *Responder ) Blank() {
    rsp.Error( rsp.Status(), nil )
}

func( rsp *Responder ) Redirect( url string ) {}

func( rsp *Responder ) Template( tmpl *template.Template, data interface{} ) {

    if tmpl == nil {
        rsp.Error( 404, nil )
        return
    }

    // Exec
    buf := bytes.NewBuffer( nil )
    err := tmpl.Execute( buf, data )
    if err != nil {
        rsp.Error( 500, nil )
        return
    }

    final := buf.String()

    // Localize
    if rsp.req.localizer != nil {
        final = rsp.req.localizer.L( final )
    }

    // Serve
    rsp.Html( final )
    return

}

func( rsp *Responder ) Asset( ast *volume.Asset ) {

    if ast == nil {
        rsp.Error( 404, nil )
        return
    }

    final := string( ast.Bytes() )

    // Localize
    if rsp.req.localizer != nil {
        final = rsp.req.localizer.L( final )
    }

    // Serve
    rdskr := bytes.NewReader( []byte( final ) )
    rsp.File( ast.Name(), ast.ModTime(), rdskr )
    return

}