package volume

import (
    "bytes"
    
    "../system"
)

func CompileScss( src string ) ( string, error ) {
    
    // Cmd
    sti := bytes.NewReader( []byte( src ) )
    sto := bytes.NewBuffer( nil )
    // sass --stdin
    // : this makes sass to accept input from the stdin
    // : e.g. echo 'body {color:white;}' | sass --stdin
    err := system.Exec( sti, sto, nil, "sass", "--stdin", "--style=compressed" )

    if err != nil {
        return "", ErrCompileFailure.Append( "sass", err )
    }
    
    return sto.String(), nil
    
}