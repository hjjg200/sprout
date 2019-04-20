package volume

import (
    "html/template"

    "../cache"
    "../i18n"
)

type Volume interface {

    // Getters
    Asset( string ) *Asset
    I18n() ( *i18n.I18n )
    Localizer( string ) *i18n.Localizer
    Template( string ) *template.Template
    
    // Setters
    SetFallback( Volume )

    // cache.Porter
    cache.Porter

}