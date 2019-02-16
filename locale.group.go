package sprout

/*
 + LOCALE GROUP
 */

type LocaleGroup struct{
    locales map[string] *Locale
    default_locale string
    threshold int
}

func( )