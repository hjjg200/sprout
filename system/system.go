package system

import (
    "runtime"

    "../environ"
)

func init() {

    // Check OS
    switch runtime.GOOS {
    case "windows", "linux", "darwin":
    default:
        environ.Logger.Panicln( ErrOSNotSupported.Append( runtime.GOOS ) )
    }

}