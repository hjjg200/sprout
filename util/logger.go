package util

import (
    "fmt"
    "io"
    "os"
    "strings"
    "time"
)

type Logger struct{
    colors []io.Writer
    monos  []io.Writer
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
        colors: []io.Writer{ os.Stdout },
        monos: []io.Writer{},
    }
}

func( lgr *Logger ) print( prefix string, args ...interface{} ) {

    if len( lgr.colors ) + len( lgr.monos ) > 0 {

        out := prefix + " "
        out += secondsFromStart()
        out += " - "

        for i := range args {
            if i > 0 {
                out += " "
            }
            out += fmt.Sprint( args[i] )
        }

        for _, w := range lgr.colors {
            fmt.Fprint( w, out )
        }
        for _, w := range lgr.monos {
            fmt.Fprint( w, stripAnsiColor( out ) )
        }

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

func( lgr *Logger ) AddColorWriter( w io.Writer ) {
    lgr.colors = append( lgr.colors, w )
}

func( lgr *Logger ) AddMonoWriter( w io.Writer ) {
    lgr.monos = append( lgr.monos, w )
}

func stripAnsiColor( str string ) string {

    for i := 0; i < 5; i++ {
        pos := strings.Index( str, "\033" )
        if pos == -1 {
            break
        }
        pos2 := strings.Index( str[pos:], "m" ) + pos
        if pos2 == -1 {
            break
        }
        if pos2 == len( str ) - 1 {
            str = str[:pos]
        } else {
            str = str[:pos] + str[pos2 + 1:]
        }
    }
    return str

}