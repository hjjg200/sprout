package volume

import (
    "bytes"
    "path/filepath"

    "github.com/hjjg200/sprout/system"
    "github.com/hjjg200/sprout/util/errors"
)

type Compilers struct {
    i2o map[string] string
    o2i map[string] []string
    funcs map[string] func( []byte ) ( []byte, error )
}

var DefaultCompilers = NewCompilers()

func NewCompilers() *Compilers {
    return &Compilers{
        i2o: make( map[string] string ),
        o2i: make( map[string] []string ),
        funcs: make( map[string] func( []byte ) ( []byte, error ) ),
    }
}

func( cps *Compilers ) InputOf( path string ) ( []string, bool ) {
    ext := filepath.Ext( path )
    in, ok := cps.o2i[ext]
    if !ok {
        return nil, ok
    }
    ret := []string{}
    for _, i := range in {
        ret = append( ret, filepath.ToSlash( filepath.Dir( path ) ) + "/" + BaseWithoutExt( path ) + i )
    }
    return ret, ok
}

func( cps *Compilers ) OutputOf( path string ) ( string, bool ) {
    ext := filepath.Ext( path )
    out, ok := cps.i2o[ext]
    if !ok {
        return "", ok
    }
    return filepath.ToSlash( filepath.Dir( path ) ) + "/" + BaseWithoutExt( path ) + out, ok
}

func( cps *Compilers ) Compile( ast *Asset ) ( *Asset, error ) {

    ext    := filepath.Ext( ast.Name() )
    fn, ok := cps.funcs[ext]

    if !ok {
        return nil, errors.ErrCompileFailure.Append( ext )
    }

    out, err := fn( ast.Bytes() )
    if err != nil {
        return nil, errors.ErrCompileFailure.Append( err )
    }

    outName, _ := cps.OutputOf( ast.Name() )
    rd         := bytes.NewReader( out )
    return NewAsset( outName, rd, ast.ModTime() ), nil

}

func( cps *Compilers ) Put( in []string, out string, fn func( []byte ) ( []byte, error ) ) {

    for _, i := range in {
        cps.funcs[i] = fn
        cps.i2o[i]   = out
    }
    cps.o2i[out] = in

}

func init() {
    DefaultCompilers.Put( []string{ ".scss" }, ".css", CompileScss )
}

func CompileScss( in []byte ) ( []byte, error ) {

    // Cmd
    sti := bytes.NewReader( in )
    sto := bytes.NewBuffer( nil )
    // sass --stdin
    // : this makes sass to accept input from the stdin
    // : e.g. echo 'body {color:white;}' | sass --stdin
    err := system.Exec( sti, sto, nil, "sass", "--stdin", "--style=compressed" )

    if err != nil {
        return nil, errors.ErrCompileFailure.Append( "sass", err )
    }

    return sto.Bytes(), nil

}
