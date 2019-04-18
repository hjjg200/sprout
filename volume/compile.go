package volume

import (
    "bytes"

    "../system"
)

type Compiler struct {
    outExt string
    fn func( *Asset ) ( *Asset, error )
}

var Compilers = map[string] Compiler {
    ".scss": Compiler{ ".css", CompileScss },
}

func CompileScss( ast *Asset ) ( *Asset, error ) {

    // Vars
    bwoext := BaseWithoutExt( ast.Name() )

    // Cmd
    sti := bytes.NewReader( ast.Bytes() )
    sto := bytes.NewBuffer( nil )
    // sass --stdin
    // : this makes sass to accept input from the stdin
    // : e.g. echo 'body {color:white;}' | sass --stdin
    err := system.Exec( sti, sto, nil, "sass", "--stdin", "--style=compressed" )

    if err != nil {
        return nil, ErrCompileFailure.Append( "sass", err )
    }

    // Make
    rd := bytes.NewReader( sto.Bytes() )

    return NewAsset( bwoext + ".css", rd, ast.ModTime() ), nil

}