package sprout

/*
 + REQUEST
 *
 * A request is the context of a http request, containing locale, writer, request, etc.
 * Requests have functions that are for responding to requests
 */

type Request struct {
    writer http.ResponseWriter
    body   *http.Request
    locale string
}

func( _request *Request ) RemoveLocaleFromURL() {
    if len( _request.locale ) == 0 {
        return
    }
    // Check if the url starts with locale name
    _prefix := "/" + _request.locale
    if strings.HasPrefix( _request.URL.Path, _prefix ) {
        _request.URL.Path = "/" + strings.TrimPrefix( _prefix )
    }
}
func( _request *Request ) Writer() http.ResponseWriter {
    return _request.writer
}
func( _request *Request ) Body() *http.Request {
    return _request.body
}
func( _request *Request ) Locale() string {
    return _request.locale
}
func( _request *Request ) WriteJSON( _data interface{} ) {
    _error := json.NewEncoder( _request.writer ).Encode( _data )
    if _error != nil {
        _Error := Error{
            code: 500,
            details: ErrorFactory().New( "request:", _error )
        }
        _request.WriteErorr( _Error )
    }
}
func( _request *Request ) WriteStatus( _code int ) {

    _request.Writer.Header().Set( "Content-Type", "text/html; charset=utf-8" )
    _code_string := fmt.Sprint( _code )
    _message     := SproutVariables().HTTPStatusMessages()( _code )
    _html        := `<!doctype html>
<html>
    <head>
        <title>` + _code_string + " " + _message + `</title>
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
            <h1>` + _code_string + `</h1>
            <h3>` + _message + `</h3>
            <footer>` + SproutVariables().AppName() + " " + SproutVariables().Version() + `<br />on ` + SproutVariables().OS() + `</footer>
        </div>
    </body>
</html>`
    _request.writer.Write( []byte( _html ) )

}
func( _request *Request ) WriteError( _error error ) {

    if _Error, _ok := _error.( Error ); _ok {
        _request.WriteStatus( _Error.code )
        if _Error.details != nil {
            SproutLogger().Warnln( _Error.details )
        }
        return
    } else {
        _request.WriteStatus( 500 )
        return
    }

}