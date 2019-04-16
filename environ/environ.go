package environ

import (
    "html/template"
    "runtime"

    "../util"
)

const (
    AppName = "sprout"
    AppVersion = "pre-alpha"

    ErrorPageTemplatePath = "template/error_page.html"
)

var (
    DefaultErrorPageTemplate *template.Template
)

var Logger = util.NewLogger()

func init() {
    DefaultErrorPageTemplate, _ = template.New( "" ).Parse( `<!doctype html>
<html>
    <head>
        <title>{{ .status }} {{ .message }}</title>
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
            <h1>{{ .status }}</h1>
            <h3>{{ .message }}</h3>
            <footer>` + AppName + " " + AppVersion + `<br />on ` + runtime.GOOS + `</footer>
        </div>
    </body>
</html>` )
}