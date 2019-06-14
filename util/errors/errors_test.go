package errors

import (
    "fmt"
    "testing"
)

func TestErrors01( t *testing.T ) {
    fmt.Println( testError01_a() )
}
func testError01_a() error {
    return Stack( testError01_b() )
}
func testError01_b() error {
    return Stack( testError01_c() )
}
func testError01_c() error {
    return Stack( testError01_d() )
}
func testError01_d() error {
    return fmt.Errorf( "IO error" )
}