package sprout

/*
 + LOCALIZER
 */

type Localizer struct{
    locales map[string] *Locale
    default_locale string
    threshold int
    left_delimiter string
    right_delimiter string
}

func( )