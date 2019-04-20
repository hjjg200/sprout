package system

import (
    "bytes"
    "io"
    "os/exec"
    "strings"
    "runtime"

    "../util"
)

func NewCmd( args ...string ) *exec.Cmd {
    switch runtime.GOOS {
    case "linux", "darwin":
        return exec.Command( "bash", []string{ "-c",  strings.Join( args, " " ) }... )
    case "windows":
        return exec.Command( "cmd", append( []string{ "/C" }, args... )... )
    }
    return nil
}

func Exec( stdin io.Reader, stdout, stderr io.Writer, args ...string ) error {

    // Set
    e := NewCmd( args... )
    e.Stdin  = stdin
    e.Stdout = stdout

    // Error writer
    errbuf  := bytes.NewBuffer( nil )
    writers := []io.Writer{ errbuf }
    if stderr != nil {
        writers = append( writers, stderr )
    }
    e.Stderr = io.MultiWriter( writers... )

    // Run
    err := e.Run()

    // Check err
    if err != nil {
        return util.NewError( 500, err, errbuf.String() )
    }

    return nil

}

func DoesCommandExist( cmd string ) ( bool, error ) {

    var (
        err error
        out bytes.Buffer
        e   *exec.Cmd
    )

    switch runtime.GOOS {
    case "linux", "darwin":
        s := "if command -v " + cmd + " > /dev/null 2>&1; then echo 'true'; fi"
        e  = exec.Command( "bash", "-c", s )
    case "windows":
        s := "where /Q " + cmd + " & if %errorlevel%==0 echo true"
        e  = exec.Command( "cmd", "/C", s )
    }

    e.Stdout = &out
    err = e.Run()

    if err != nil {
        return false, err
    }

    r := out.String()
    if r[:4] == "true" {
        return true, nil
    }

    return false, nil

}