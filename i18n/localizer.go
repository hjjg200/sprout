package i18n

/*
 + Localizer
 *
 * Localizer is a shorthand method for locale-specific localizing.
 */

type Localizer struct {
    parent *I18n
    lcName string
}

func NewLocalizer( i1 *I18n, lcName string ) ( *Localizer, error ) {
    if !i1.HasLocale( lcName ) {
        return nil, ErrLocaleNonExistent.Append( lcName )
    }
    return &Localizer{
        parent: i1,
        lcName: lcName,
    }, nil
}

func( lczr *Localizer ) L( src string ) string {
    return lczr.Localize( src )
}
func( lczr *Localizer ) Localize( src string ) string {
    return lczr.parent.Localize( lczr.lcName, src )
}