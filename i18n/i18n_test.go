package i18n

import (
    "fmt"
    "testing"
)

func TestI18n01( t *testing.T ) {
    i := New()
    i.ImportDirectory( "../test/TestI18n01" )
    i.defaultLocale = "en"
    fmt.Println( i.ParseAcceptLanguage( "en;q=0.9, ko;q=0.95, ja;q=1, *;q=0.5" ) )
    fmt.Println( i.Locales() )
    fmt.Println( i.locales["en-us"].set )
    fmt.Println( i.L( "en", `
{% abc %}
{% b.person %}    {% b.def.ghi %}
` ) )
}