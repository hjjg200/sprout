package network

import (
    "bytes"
    "encoding/json"
    "fmt"
    "html/template"
    "io"
    "net/http"
    "strings"
    "time"

    "github.com/hjjg200/sprout/environ"
    "github.com/hjjg200/sprout/i18n"
    "github.com/hjjg200/sprout/util"
    "github.com/hjjg200/sprout/util/errors"
    "github.com/hjjg200/sprout/volume"
)

type Request struct {
    id          int64
    body        *http.Request
    writer      http.ResponseWriter
    status      int
    wroteHeader bool
    localizer   *i18n.Localizer
    space       *Space
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

func( req *Request ) Space() *Space {
    return req.space
}

func( req *Request ) Volume() volume.Volume {
    // Returns DefaultVolume if space is nil
    if req.space == nil {
        return volume.DefaultVolume
    }
    return req.space.Volume()
}

func( req *Request ) Vars() []string {
    return req.vars
}

func( req *Request ) Write( p []byte ) ( int, error ) {
    return req.writer.Write( p )
}

func( req *Request ) SetStatus( status int ) {

    if req.wroteHeader {
        environ.Logger.Panicln( errors.ErrDifferentStatusCode.Append( "ID", req.ID() ) )
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
        http.SetCookie( req, i1.MakeCookie( lcName ) )
        return

    }

    req.localizer = nil

}

// POP METHODS

func( req *Request ) Pop( status int, text string, ctype string ) {
    req.Header().Set( "content-type", ctype )
    req.SetStatus( status )
    req.Write( []byte( text ) )
}

func( req *Request ) PopError( status int, err error ) {

    var (
        tmpl *template.Template
        msg = util.HttpStatusMessages[status]
    )

    tmpl = req.Volume().Template( environ.ErrorPageTemplatePath )

    if tmpl == nil {
        environ.Logger.Warnln( "Missing template", environ.ErrorPageTemplatePath )
        req.PopText( status, fmt.Sprintf( "%d %s", status, msg ) )
        return
    }

    // Map
    data := map[string] interface{} {
        "status": status,
        "message": msg,
    }

    // Raise
    if err != nil {
        environ.Logger.Warnln( "ID", req.ID(), err )
    }

    req.PopTemplate( status, tmpl, data )

}

func( req *Request ) PopBlank( status int ) {
    req.PopError( status, nil )
}

func( req *Request ) PopText( status int, text string ) {
    req.Pop( status, text, "text/plain;charset=utf-8" )
}

func( req *Request ) PopHtml( status int, html string ) {
    req.Pop( status, html, "text/html;charset=utf-8" )
}

func( req *Request ) PopTemplate( status int, tmpl *template.Template, data interface{} ) {

    if tmpl == nil {
        req.PopError( 404, nil )
        return
    }

    // Exec
    buf := bytes.NewBuffer( nil )
    err := tmpl.Execute( buf, data )
    if err != nil {
        req.PopError( 500, nil )
        return
    }

    final := buf.String()

    // Localize
    if req.localizer != nil {
        final = req.localizer.L( final )
    }

    // Serve
    req.PopHtml( status, final )
    return

}

func( req *Request ) popJson( status int, obj interface{}, pretty bool ) {

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
        req.PopError( 500, errors.ErrMalformedJson.Append( err ) )
    }

    req.Pop( status, string( p ), "text/json;charset=utf-8" )

}

func( req *Request ) PopJson( status int, obj interface{} ) {
    req.popJson( status, obj, false )
}

func( req *Request ) PopPrettyJson( status int, obj interface{} ) {
    req.popJson( status, obj, true )
}

func( req *Request ) PopAsset( ast *volume.Asset ) {

    if ast == nil {
        req.PopBlank( 404 )
        return
    }

    // Check version
    v, ok  := req.body.URL.Query()[c_queryAssetVersion]
    astVer := ast.Version()

    switch {
    case !ok,
        len( v ) != 1,
        v[0] != astVer:
        params := req.body.URL.Query()
        params.Set( c_queryAssetVersion, astVer )
        req.body.URL.RawQuery = params.Encode()
        req.PopRedirect( 307, req.body.URL.String() )
        return
    }

    final := string( ast.Bytes() )

    // Localize
    if req.localizer != nil {
        final = req.localizer.L( final )
    }

    // Serve
    rdskr := bytes.NewReader( []byte( final ) )
    req.PopFile( ast.Name(), ast.ModTime(), rdskr )
    return

}

func( req *Request ) PopFile( name string, modTime time.Time, rdskr io.ReadSeeker ) {
    http.ServeContent( req, req.body, name, modTime, rdskr )
}

func( req *Request ) PopAttachment( name string, modTime time.Time, rdskr io.ReadSeeker ) {

    req.Header().Set( "Content-Type", "application/octet-stream" )
    req.Header().Set( "Content-Transfer-Encoding", "Binary" )
    req.Header().Set( "Content-Disposition", "attachment; filename=\"" + name + "\"" )
    req.PopFile( name, modTime, rdskr )

}

func( req *Request ) PopRedirect( code int, url string ) {
    http.Redirect( req.writer, req.body, url, code )
}