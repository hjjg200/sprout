package system

import (
    "runtime"

    "github.com/hjjg200/sprout/environ"
)

func init() {

    // Check OS
    switch runtime.GOOS {
    case "windows", "linux", "darwin":
    default:
        environ.Logger.Panicln( ErrOSNotSupported.Append( runtime.GOOS ) )
    }

}
