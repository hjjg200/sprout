package sprout

import (
    "bytes"
    "errors"
    "runtime"
    "os/exec"
)

/*
 + os.go
 |
 + This file deals with the OS-specific operations
 */


func checkOS() error {
    var (
        ErrNotSupportedOS = errors.New( "sprout: the OS is not supported" )
    )
    switch runtime.GOOS {
    case "windows", "linux", "darwin", "freebsd", "openbsd":
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
    case "linux", "darwin", "freebsd", "openbsd":
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

func ( s *Sprout ) runCommand( cmd string ) error {
    var e *exec.Cmd
    switch envOS {
    case "linux", "darwin", "freebsd", "openbsd":
        e = exec.Command( "bash", "-c", cmd )
    case "windows":
        e = exec.Command( "cmd", "/C", cmd )
    }
    return e.Run()
}