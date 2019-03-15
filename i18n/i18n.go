package i18n

import (

)

type I18n struct {
    locales map[string] *Locale
    defaultLocale string
    cookie string
    queryParameter string
    leftDelimiter string
    rightDelimiter string
}

func New() *I18n {
    return &I18n{
        locales: make( map[string] *Locale ),
        defaultLocale: "",
        cookie: c_defaultCookie,
        queryParameter: c_defaultQueryParameter,
        leftDelimiter: c_defaultLeftDelimiter,
        rightDelimiter: c_defaultRightDelimiter,
    }
}

func( i1 *I18n ) SetDefaultLocale( locale string ) {

}
func( i1 *I18n ) SetQueryParameter( param string ) {
    
}
func( i1 *I18n ) SetCookie( cookie string ) {
    if len( cookie ) == 0 {
        panic( ErrInvalidParamter )
        return
    }
    i1.cookie = cookie
}
func( i1 *I18n ) SetDelimiters( left, right string ) {
    if len( left ) == 0 || len( right ) == 0 {
        panic( ErrInvalidDelimiters )
        return
    }
    i1.leftDelimiter = left
    i1.rightDelimiter = right
}