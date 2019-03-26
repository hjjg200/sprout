package util

type String string

func( s String ) IsIn( args ...string ) bool {
    for i := range args {
        if string( s ) == args[i] {
            return true
        }
    }
    return false
}