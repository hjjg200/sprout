package network

import (
    "../util"
)

var (

    ErrStartingServer = util.NewError( 500, "failed to start the server" )
    ErrStoppingServer = util.NewError( 500, "failed to stop the server" )
    ErrServerExited = util.NewError( 500, "the server exited" )
    
)