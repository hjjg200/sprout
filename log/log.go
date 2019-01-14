package log

import (
    "fmt"
    "io"
    "os"
    "time"
)

const (
    stateInfo   = "INFO"
    stateWarn   = "WARN"
    stateSevere = "SEVERE"
    timeFormat  = "15:04:05"
)

var (
    output io.Writer = os.Stderr
)



func SetOutput( w io.Writer ) { output = w }

func formatState( state string ) string {
    t := time.Now().Format( timeFormat )
    return fmt.Sprintf( "[%s %s]", t, state )
}
func printLog( state string, v ...interface{} ) { 
    v = append( []interface{}{ formatState( state ) + " " }, v... )
    fmt.Fprint( output, v... )
}
func printlnLog( state string, v ...interface{} ) {
    v = append( []interface{}{ formatState( state ) }, v... )
    fmt.Fprintln( output, v... )
}

func Info( v ...interface{} ) { printLog( stateInfo, v... ) }
func Warn( v ...interface{} ) { printLog( stateWarn, v... ) }
func Severe( v ...interface{} ) {
    printLog( stateSevere, v... )
    os.Exit( 1 )
}

func Infof( format string, v ...interface{} ) { printLog( stateInfo, fmt.Sprintf( format, v... ) ) }
func Warnf( format string, v ...interface{} ) { printLog( stateWarn, fmt.Sprintf( format, v... ) ) }
func Severef( format string, v ...interface{} ) {
    printLog( stateSevere, fmt.Sprintf( format, v... ) )
    os.Exit( 1 )
}

func Infoln( v ...interface{} ) { printlnLog( stateInfo, v... ) }
func Warnln( v ...interface{} ) { printlnLog( stateWarn, v... ) }
func Severeln( v ...interface{} ) {
    printlnLog( stateSevere, v... )
    os.Exit( 1 )
}