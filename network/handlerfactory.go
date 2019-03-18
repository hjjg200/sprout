package network

type handlerFactory struct{}
var HandlerFactory = &handlerFactory{}

func( hf *HandlerFactory ) Asset( ast *Asset ) Handler {

    buf := bytes.NewBuffer()
    io.Copy( buf, ast )

    return func( req *Request ) bool {

        final := buf.String()

        // Localize
        if req.localizer != nil {
            final = req.space.volume.i18n.L( req.locale, final )
        }

        // Serve
        rdskr := bytes.NewReader( []byte( final ) )
        http.ServeContent( req.writer, req.body, ast.Name(), ast.ModTime(), rdskr )
        return true

    }
}

func( hf *HandlerFactory ) Template( tmpl *template.Template, dataFunc func( *Request ) interface{} ) Handler {
    return func( req *Request ) bool {

        // Exec
        buf := bytes.NewBuffer()
        tmpl.Execute( buf, dataFunc( req ) )
        final := buf.String()

        // Localize
        if req.localizer != nil {
            final = req.space.volume.i18n.L( req.locale, final )
        }

        // Serve
        req.writer.Header().Set( "content-type", "text/html;charset=utf-8" )
        req.writer.Write( []byte( final ) )
        return true

    }
}

func( hf *HandlerFactory ) Status( code int ) Handler {
    return func( req *Reqeust ) bool {
        req.WriteStatus( code )
        return true
    }
}

func( hf *HandlerFactory ) Authenticator( auther func( *Request ) bool, realm string ) Handler {
    return func( req *Request ) bool {
        if auther( req ) == true {
            // Returns false so that the following handlers can handle the request
            return false
        }
        // Set the authentication realm
        req.writer.Header().Set( "WWW-Authenticate", "Basic realm=\"" + realm + "\"" )

        // Returns the 401 handler
        return HandlerFactory.Status( 401 )( req )
    }
}

func( hf *HandlerFactory ) Text( text, name string ) Handler {

    mimeType := mime.TypeByExtension( filepath.Ext( name ) )

    return func( req *Request ) bool {

        // Set
        req.writer.Header().Set( "content-type", mimeType + ";charset=utf-8" )

        // Serve
        rdskr := bytes.NewReader( []byte( text ) )
        http.ServeContent( req.writer, req.body, name, time.Now(), rdskr )

        return true

    }

}

func( hf *HandlerFactory ) File( osPath string ) Handler {

}

func( hf *HandlerFactory ) Directory( osPath string ) Handler {

}