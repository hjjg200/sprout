package util

import (
    "testing"
)

func TestLogger01( t *testing.T ) {
    LogEnabled = true
    Logger.OKln( "hi", "OK" )
    Logger.Warnln( "warning!" )
    Logger.Severeln( "severe!" )
}