package system

import (
    "../util"
)

var (
    
    // ERROR
    ErrOSNotSupported = util.NewError( 500, "the OS is not supported" )
    
)