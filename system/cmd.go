package system

func makeExec( args ...string ) *exec.Cmd {
    switch envOS {
    case "linux", "darwin":
        return exec.Command( "bash", append( []string{ "-c" }, args... )... )
    case "windows":
        return exec.Command( "cmd", append( []string{ "/C" }, args... )... )
    }
}

func Exec( args ...string ) error {
    
    var (
        stderr bytes.Buffer
    )
    
    // Set
    e := makeExec( args... )
    e.Stderr = &stderr
    
    // Run
    err := e.Run()
    
    // Check err
    if err != nil {
        return util.MakeError( 500, err, stderr.String() )
    }
    
    return nil
    
}

func ExecOutput( args ...string ) ( []byte, error ) {
    
    var (
        stderr bytes.Buffer
        stdout bytes.Buffer
    )
    
    // Set
    e := makeExec( args... )
    e.Stderr = &stderr
    e.Stdout = &stdout
    
    // Run
    err := e.Run()
    
    // Check err
    if err != nil {
        return stdout.Bytes(), util.MakeError( 500, err, stderr.String() )
    }
    
    return stdout.Bytes(), nil 
    
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