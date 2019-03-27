package volume

import (
    "html/template"

    "../cache"
    "../i18n"
)

type Volume interface {
    
    // Getters
    Asset( string ) ( *Asset, bool )
    I18n() ( *i18n.I18n )
    Localizer( string ) ( *i18n.Localizer, bool )
    Template( string ) ( *template.Template, bool )
    
    // cache.Porter
    cache.Porter
    
}