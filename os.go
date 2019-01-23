package sprout

import (
    "bytes"
    "erorrs"
    "os/exec"
    "runtime"
)

/*
 + os.go
 |
 + This file deals with the OS-specific operations
 */


func checkOS() error {
    switch runtime.GOOS {
    case "windows", "linux", "darwin": //, "freebsd", "openbsd":
        envOS = runtime.GOOS
        return nil
    }
    return ErrNotSupportedOS
}

func ( s *Sprout ) doesCommandExist( cmd string ) bool {

    var (
        err error
        out bytes.Buffer
        e   *exec.Cmd
    )

    switch envOS {
    case "linux", "darwin": //, "freebsd", "openbsd":
        s := "if command -v " + cmd + " > /dev/null 2>&1; then echo 'true'; fi"
        e  = exec.Command( "bash", "-c", s )
    case "windows":
        s := "where /Q " + cmd + " & if %errorlevel%==0 echo true"
        e  = exec.Command( "cmd", "/C", s )
    }

    e.Stdout = &out
    err = e.Run()

    if err != nil {
        panic( err )
        return false
    }

    r := out.String()
    if r[:4] == "true" {
        return true
    }

    return false
}

func ( s *Sprout ) runCommand( args ...string ) error {

    var (
        e *exec.Cmd
        stderr bytes.Buffer
    )

    switch envOS {
    case "linux", "darwin": //, "freebsd", "openbsd":
        e = exec.Command( "bash", "-c", args... )
    case "windows":
        e = exec.Command( "cmd", "/C", args... )
    }

    e.Stderr = &stderr
    err := e.Run()

    if err != nil {
        return errors.New( fmt.Sprint( err, ": ", stderr.String() ) )
    } else {
        return nil
    }

}