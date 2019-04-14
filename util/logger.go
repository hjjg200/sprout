package util

import (
    "fmt"
    "io"
    "os"
    "time"
)

type Logger struct{
    w io.Writer
    timeType int
}

var (
    logStartTime time.Time
)

func init() {
    logStartTime = time.Now()
}

func secondsFromStart() string {
    return fmt.Sprintf( "%10.3f", float64( time.Now().Sub( logStartTime ) ) / float64( time.Second ) )
}

func NewLogger() *Logger {
    return &Logger{
        w: os.Stdout,
    }
}

func( lgr *Logger ) print( prefix string, args ...interface{} ) {

    if lgr.w != nil {

        out := prefix + " "
        out += secondsFromStart()
        out += " - "

        for i := range args {
            if i > 0 {
                out += " "
            }
            out += fmt.Sprint( args[i] )
        }

        fmt.Fprint( lgr.w, out )

    }

}

func( lgr *Logger ) println( prefix string, args ...interface{} ) {
    args = append( args, "\n" )
    lgr.print( prefix, args... )
}

func( lgr *Logger ) OKln( args ...interface{} ) {
    lgr.println( "[  \033[32;1mOK\033[0m  ]", args... )
}

func( lgr *Logger ) Warnln( args ...interface{} ) {
    lgr.println( "[ \033[33;1mWARN\033[0m ]", args... )
}

func( lgr *Logger ) Severeln( args ...interface{} ) {
    lgr.println( "[\033[31;1mSEVERE\033[0m]", args... )
    os.Exit( 1 )
}

func( lgr *Logger ) Panicln( args ...interface{} ) {
    lgr.print( "\033[41;37;1m[PANIC!]\033[0m", args... )
    panic( "" )
}

func( lgr *Logger ) SetOutput( w io.Writer ) {
    lgr.w = w
}

func( lgr *Logger ) SetTimeType( t int ) {
    lgr.timeType = t
}