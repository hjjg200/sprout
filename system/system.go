package system

import (
    "fmt"
    "runtime"

    "github.com/hjjg200/sprout/environ"
)

func init() {

    // Check OS
    switch runtime.GOOS {
    case "windows", "linux", "darwin":
    default:
        environ.Logger.Panicln( fmt.Errorf( "The OS, %s, is not supported", runtime.GOOS ) )
    }

}
