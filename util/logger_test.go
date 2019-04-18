package util

import (
    "os"
    "testing"
)

func TestLogger01( t *testing.T ) {

    f, _ := os.OpenFile( "../test/TestLogger01/test.log", os.O_CREATE | os.O_WRONLY | os.O_TRUNC, 0644 )
    defer f.Close()
    lgr := NewLogger()
    lgr.AddMonoWriter( f )

    lgr.OKln( "hi", "OK" )
    lgr.Warnln( "warning!" )
    lgr.Severeln( "severe!" )

}