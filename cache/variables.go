package cache

import (
    "github.com/hjjg200/sprout/util"
)

const (
    switchRead = 0 + iota
    switchWrite
)

var (

    // ERRORS
    ErrEntryAccessFailed = util.NewError( 500, "failed to access the entry" )
    ErrEntryNotFound = util.NewError( 500, "the entry was not found" )
    ErrEntryWriteFail = util.NewError( 500, "failed to write the entry" )

)
