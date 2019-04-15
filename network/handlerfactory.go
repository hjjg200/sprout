package network

import (
    "bytes"
    "fmt"
    "html/template"
    "runtime"

    "../environ"
    "../util"
    "../volume"
)

type handlerFactory struct{}
var HandlerFactory = &handlerFactory{}

func( hf *handlerFactory ) Asset( ast *volume.Asset ) Handler {

    return func( req *Request ) bool {

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
        req.RespondFile( ast.Name(), ast.ModTime(), rdskr )
        return true

    }

}

func( hf *handlerFactory ) Template( tmpl *template.Template, dataFunc func( *Request ) interface{} ) Handler {
    return func( req *Request ) bool {

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
        req.Respond( 200, final )
        return true

    }
}

func( hf *handlerFactory ) Status( code int ) Handler {
    return func( req *Request ) bool {

        // Content
        c   := fmt.Sprint( code )
        msg := util.HttpStatusMessages[code]
        t := `<!doctype html>
<html>
    <head>
        <title>` + c + " " + msg + `</title>
        <style>
            html {
                font-family: sans-serif;
                line-height: 1.0;
                padding: 0;
            }
            body {
                color: hsl( 220, 5%, 45% );
                text-align: center;
                padding: 10px;
                margin: 0;
            }
            div {
                border: 1px dashed hsl( 220, 5%, 88% );
                padding: 20px;
                margin: 0 auto;
                max-width: 300px;
                text-align: left;
            }
            h1, h2, h3 {
                display: block;
                margin: 0 0 5px 0;
            }
            footer {
                color: hsl( 220, 5%, 68% );
                font-family: monospace;
                font-size: 1em;
                text-align: right;
                line-height: 1.3;
            }
        </style>
    </head>
    <body>
        <div>
            <h1>` + c + `</h1>
            <h3>` + msg + `</h3>
            <footer>` + environ.AppName + " " + environ.AppVersion + `<br />on ` + runtime.GOOS + `</footer>
        </div>
    </body>
</html>`

        req.Respond( code, t )
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