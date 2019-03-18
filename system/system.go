package system

import (
    "runtime"
)

func init() {
    
    // Check OS
    switch runtime.GOOS {
    case "windows", "linux", "darwin":
    default:
        panic( ErrOSNotSupported.Append( runtime.GOOS ) )
    }
        
}