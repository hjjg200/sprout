package system

import (
    "runtime"

    "github.com/hjjg200/sprout/environ"
    "github.com/hjjg200/sprout/util/errors"
)

func init() {

    // Check OS
    switch runtime.GOOS {
    case "windows", "linux", "darwin":
    default:
        environ.Logger.Panicln( errors.ErrOSNotSupported.Raise( runtime.GOOS ) )
    }

}
