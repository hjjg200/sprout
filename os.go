package sprout

/*
 + OS
 */

type sprout_os struct{}
var  static_os = &sprout_os{}

func OS() *sprout_os {
    return static_os
}
func( _os *sprout_os ) CommandExists( _command string ) bool {

    var (
        _error error
        _out   bytes.Buffer
        _exec  *exec.Cmd
    )

    switch SproutVariables().OS() {
    case "linux", "darwin": //, "freebsd", "openbsd":
        _exec  = exec.Command(
            "bash",
            "-c",
            "if command -v " + _command + " > /dev/null 2>&1; then echo 'true'; fi",
        )
    case "windows":
        _exec  = exec.Command(
            "cmd",
            "/C",
            "where /Q " + _command + " & if %errorlevel%==0 echo true",
        )
    }

    _exec.Stdout = &_out
    _error = e.Run()

    if _error != nil {
        return false
    }

    return _out.String()[:4] == "true"

}
func( _os *sprout_os ) RunCommand( _args ...string ) error {

    var (
        _exec      *exec.Cmd
        _error_out bytes.Buffer
    )

    switch SproutVariables().OS() {
    case "linux", "darwin": //, "freebsd", "openbsd":
        _exec = exec.Command( "bash", append( []string{ "-c" }, args... )... )
    case "windows":
        _exec = exec.Command( "cmd", append( []string{ "/C" }, args... )... )
    }

    _exec.Stderr  = &_error_out
    _error       := e.Run()

    if _error != nil {
        return ErrorFactory().New( 500, _error, _error_out.String() )
    } else {
        return nil
    }

}