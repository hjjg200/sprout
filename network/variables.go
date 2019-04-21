package network

import (
    "github.com/hjjg200/sprout/util"
)

var (

    // ERROR DEFINITIONS
    ErrStartingServer = util.NewError( 500, "failed to start the server" )
    ErrStoppingServer = util.NewError( 500, "failed to stop the server" )
    ErrServerExited = util.NewError( 500, "the server exited" )
    ErrRequestClosed = util.NewError( 500, "the request is already closed" )
    ErrDifferentStatusCode = util.NewError( 500, "attempted to write different status code" )
    ErrMalformedJson = util.NewError( 500, "the given json is malformed. Defaulting to an empty JSON object" )

)
