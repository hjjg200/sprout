package network

import (
    "html/template"

    "../volume"
)

type handlerFactory struct{}
var HandlerFactory = &handlerFactory{}

func( hf *handlerFactory ) Asset( ast *volume.Asset ) Handler {
    return func( req *Request ) bool {
        req.PopAsset( ast )
        return true
    }
}

func( hf *handlerFactory ) Template( tmpl *template.Template, dataFunc func( *Request ) interface{} ) Handler {
    return func( req *Request ) bool {
        req.PopTemplate( 200, tmpl, dataFunc( req ) )
        return true
    }
}

func( hf *handlerFactory ) Status( status int ) Handler {
    return func( req *Request ) bool {
        req.PopBlank( status )
        return true
    }
}

func( hf *handlerFactory ) BasicAuth( auther func( string, string ) bool, realm string ) Handler {
    return func( req *Request ) bool {
        // Id and pass
        id, pw, ok := req.body.BasicAuth()

        if ok && auther( id, pw ) == true {
            // Returns false so that the following handlers can handle the request
            return false
        }
        // Set the authentication realm
        req.Header().Set( "WWW-Authenticate", "Basic realm=\"" + realm + "\"" )

        // Returns the 401 handler
        return HandlerFactory.Status( 401 )( req )
    }
}

/*
func( hf *handlerFactory ) File( osPath string ) Handler {

}

func( hf *handlerFactory ) Directory( osPath string ) Handler {

}*/