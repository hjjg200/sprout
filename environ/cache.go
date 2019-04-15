package environ

import (
    "runtime"
    "time"

    "../cache"
)

var DefaultCache *cache.Cache

func init() {

    zeroTime := time.Time{}

    DefaultCache = cache.NewCache()
    DefaultCache.Create( "template/", zeroTime )

    w, err := DefaultCache.Create( "template/error_page.html", zeroTime )
    w.Write( []byte( `<!doctype html>
    <html>
        <head>
            <title>{{ .code }} {{ .message }}</title>
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
                <h1>{{ .code }}</h1>
                <h3>{{ .message }}</h3>
                <footer>` + environ.AppName + " " + environ.AppVersion + `<br />on ` + runtime.GOOS + `</footer>
            </div>
        </body>
    </html>` ) )

    DefaultCache.Flush()

}