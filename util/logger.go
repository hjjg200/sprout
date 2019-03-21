package util

import (
    "fmt"
    "io"
    "os"
    "time"
)

type Logger struct{
    w io.Writer
}

var (
    LogTimeFormat = "15:04:05"
)

func formattedTime() string {
    return time.Now().Format( "15:04:05" )
}

func NewLogger() *Logger {
    return &Logger{
        w: os.Stdout,
    }
}

func( lgr *Logger ) print( prefix string, args ...interface{} ) {
    if lgr.w != nil {
        fmt.Fprint( lgr.w, prefix, " ", formattedTime(), " - " )
        for i := range args {
            if i > 0 {
                fmt.Fprint( lgr.w, " " )
            }
            fmt.Fprint( lgr.w, args[i] )
        }
    }
}

func( lgr *Logger ) println( prefix string, args ...interface{} ) {
    lgr.print( prefix, args... )
    fmt.Fprint( lgr.w, "\n" )
}

func( lgr *Logger ) OKln( args ...interface{} ) {
    lgr.println( "[  OK  ]", args... )
}

func( lgr *Logger ) Warnln( args ...interface{} ) {
    lgr.println( "[ WARN ]", args... )
}

func( lgr *Logger ) Severeln( args ...interface{} ) {
    lgr.println( "[SEVERE]", args... )
    os.Exit( 1 )
}

func( lgr *Logger ) Panicln( args ...interface{} ) {
    lgr.print( "[PANIC!]", args... )
    panic( "" )
}

func( lgr *Logger ) SetOutput( w io.Writer ) {
    lgr.w = w
}