package system

import (
    "github.com/hjjg200/sprout/util"
)

var (

    // ERROR
    ErrOSNotSupported = util.NewError( 500, "the OS is not supported" )

)
