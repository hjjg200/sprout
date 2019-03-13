package i18n

import (
)

type LocalizerFactory struct{}
var  staticLocalizerFactory = &LocalizerFactory{}

func( fac *LocalizerFactory ) FromDirectory( p2d string ) ( *Localizer, error ) {
    
}
func( fac *LocalizerFactory ) FromJsons( jsons [][]byte ) ( *Localizer, error ) {
 
    // jsons might be empty
    locales := make( map[string] *Locale )
    
    // Assign each locale
    for _, json := range jsons {
        locale, err := LocaleFactory.FromJson( json )
        if err != nil {
            return nil, err
        }
        locales[locale.name] = locale
    }
    
    // Set the default locale
    var defaultLocale string
    for _, locale := range locales {
        defaultLocale = locale.name
        break
    }
    
    // Assemble
    return &Localizer{
        locales: locales,
        defaultLocale: defaultLocale,
        leftDelimiter: c_defaultLeftDelimiter,
        rightDelimiter: c_defaultRightDelimiter,
        threshold: c_defaultThreshold,
    }, nil
    
}