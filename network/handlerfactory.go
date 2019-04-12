package network

import (
    "bytes"
    "html/template"
    "mime"
    "net/http"
    "path/filepath"
    "time"

    "../volume"
)

type handlerFactory struct{}
var HandlerFactory = &handlerFactory{}

func( hf *handlerFactory ) Asset( ast *volume.Asset ) Handler {

    return func( req *Request ) int {

        if ast == nil {
            return HandlerFactory.Status( 404 )( req )
        }

        final := string( ast.Bytes() )

        // Localize
        if req.localizer != nil {
            final = req.localizer.L( final )
        }

        // Serve
        rdskr := bytes.NewReader( []byte( final ) )
        http.ServeContent( req.writer, req.body, ast.Name(), ast.ModTime(), rdskr )
        return 200

    }

}

func( hf *handlerFactory ) Template( tmpl *template.Template, dataFunc func( *Request ) interface{} ) Handler {
    return func( req *Request ) int {

        if tmpl == nil {
            return HandlerFactory.Status( 404 )( req )
        }

        // Exec
        buf := bytes.NewBuffer( nil )
        err := tmpl.Execute( buf, dataFunc( req ) )
        if err != nil {
            return HandlerFactory.Status( 500 )( req )
        }

        final := buf.String()

        // Localize
        if req.localizer != nil {
            final = req.localizer.L( final )
        }

        // Serve
        req.writer.Header().Set( "content-type", "text/html;charset=utf-8" )
        req.writer.Write( []byte( final ) )
        return 200

    }
}

func( hf *handlerFactory ) Status( code int ) Handler {
    return func( req *Request ) int {
        req.WriteStatus( code )
        return code
    }
}

func( hf *handlerFactory ) BasicAuth( auther func( string, string ) bool, realm string ) Handler {
    return func( req *Request ) int {
        // Id and pass
        id, pw, ok := req.body.BasicAuth()

        if ok && auther( id, pw ) == true {
            // Returns false so that the following handlers can handle the request
            return 100
        }
        // Set the authentication realm
        req.writer.Header().Set( "WWW-Authenticate", "Basic realm=\"" + realm + "\"" )

        // Returns the 401 handler
        return HandlerFactory.Status( 401 )( req )
    }
}

func( hf *handlerFactory ) Text( text, name string ) Handler {

    mimeType := mime.TypeByExtension( filepath.Ext( name ) )

    return func( req *Request ) int {

        // Set
        req.writer.Header().Set( "content-type", mimeType + ";charset=utf-8" )

        // Serve
        rdskr := bytes.NewReader( []byte( text ) )
        http.ServeContent( req.writer, req.body, name, time.Now(), rdskr )

        return 200

    }

}
/*
func( hf *handlerFactory ) File( osPath string ) Handler {

}

func( hf *handlerFactory ) Directory( osPath string ) Handler {

}*/