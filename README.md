# sprout

A simple web framework for Golang.

## Examples

### Creating a Server

```go
import (
    "log"
    "github.com/hjjg200/go/sprout"
)

func main() {
    s := sprout.New()
    // AddRoute returns error if the regexp is invalid
    s.AddRoute( "^/(index.html?)?$", indexPage )

    // Production server caches assets and prints them when requested
    go s.StartServer( ":80" )
    // Dev server always shows the latest version
    go s.StartDevServer( ":8080" )
    // Listen to commands from Stdin
    go s.ListenCommands()

    // Listen to signals in case you want to gracefully shutdown the servers
    s.ListenSignal()
}

func indexPage( w http.ResponseWriter, r *http.Request ) {
    w.Write( []byte( "hello" ) )
}
```