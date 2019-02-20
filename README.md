# sprout

A simple web framework for Golang.

## Overview

### Naming Convention

- Public variables -- PascalCase
- Private variables -- snake_case
- Local variables -- _snake_case ( snake case with leading an underscore )
- Variables inside function in function -- _2snake_case ( snake case with an underscore and the depth )
    - 2 here means its depth is 2

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