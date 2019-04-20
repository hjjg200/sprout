package volume

import (
    "time"
    "runtime"

    "../cache"
    "../environ"
)

var DefaultVolume Volume

func init() {

    chc := cache.NewCache()
    create := func( path, content string ) {
        w, _ := chc.Create( path, time.Now() )
        w.Write( []byte( content ) )
        w.Close()
    }

    create( environ.ErrorPageTemplatePath, `<!doctype html>
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
            <footer>` + environ.AppName + " " + environ.AppVersion + `<br />on ` + runtime.GOOS + `</footer>
        </div>
    </body>
</html>` )

    create( environ.IndexPageTemplatePath, `<!doctype html>
<html>
    <head>
        <title>{{ .title }}</title>
        <meta charset="utf-8">
    </head>
    <body>
        <ul class="breadcrumb">
            {{ range .breadcrumb }}
                <li>
                    <a href="{{ .href }}">
                        <span class="name">{{ .name }}</span>
                        <span class="slash">/</span>
                    </a>
                </li>
            {{ end }}
        </ol>
        <ol class="entries">
            {{ range .entries }}
                <li>

                </li>
            {{ end }}
        </ol>
    </body>
</html>` )

    DefaultVolume = NewBasicVolume()
    DefaultVolume.Import( chc )

}