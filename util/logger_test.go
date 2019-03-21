package util

import (
    "testing"
)

func TestLog1( t *testing.T ) {
    LogEnabled = true
    Logger.OKln( "hi", "OK" )
    Logger.Warnln( "warning!" )
    Logger.Severeln( "severe!" )
}