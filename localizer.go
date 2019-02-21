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

func( _localizer *Localizer ) EnsureLocale( _request *Request ) {
    // EnsureLocale sets the locale value of Request and puts locale cookie if there isn't
    _locale := LocaleOf( _request )
    if _locale == "" {
        // If there is no locale, default to the default locale
        _locale = _localizer.default_locale
    }
    // Set
    _request.locale = _locale
    _cookie = &http.Cookie{
        Name: LocalizerVariables().CookieName(),
        Value: _locale,
        Path: "/", // for every page
        MaxAge: 0, // persistent cookie
    }
    http.SetCookie( _request.writer, _cookie )
}
func( _localizer *Localizer ) LocaleOf( _request *Request ) string {
    if _locale := _localizer.LocaleOfURL( _request.URL.Path ); _locale != "" { return _locale }
    if _locale := _localizer.LocaleOfCookie( _request.body.Cookies() ); _locale != "" { return _locale }
    if _locale := _localizer.LocaleOfAcceptLanguage( _request.body.Header.Get( "accept-language" ) ); _locale != "" { return _locale }
    return ""
}
func( _localizer *Localizer ) LocaleOfURL( _url string ) string {
    if len( _url ) == 1 {
        return ""
    }
    _split := strings.SplitN( _url[1:], "/", 2 )
    return _localizer.LocaleOfAcceptLanguage( _split[0] )
}
func( _localizer *Localizer ) LocaleOfCookie( _cookies []*http.Cookie ) string {
    for i := range _cookies {
        if _cookies[i].Name == LocalizerVariables().CookieName() {
            if _, _ok := _localizer.locales[_cookies[i].Value]; _ok {
                return _cookies[i].Value
            }
        }
    }
    return ""
}
func( _localizer *Localizer ) LocaleOfAcceptLanguage( _accept_language string ) string {

    // Split the header
    _entries := make( accept_language_entries, 0 )
    _split   := strings.Split( _accept_language, "," )
    for i := range _split {
        // Remove whitespaces and to lowercase
        _split[i] = strings.TrimSpace( _split[i] )
        _split[i] = strings.ToLower( _split[i] )

        if _semicolon := strings.Index( _split[i], ";" ); _semicolon != -1 {
            // If there is the q-factor
            _language         := _split[i][:_semicolon]
            _q_factor, _error := strconv.ParseFloat( _split[i][_semicolon + 3:], 64 )
            if _error != nil {
                panic( _error ) // Malformed accept-language
            }
            _entries = append( _entries, accept_language_entry{
                language: _language,
                q_factor: _q_factor,
            } )
        } else {
            // Since its q-factor is default which is 1.0, it is prepended
            _entries = append( _entries, accept_language_entry{
                language: _split[i],
                q_factor: 1.0,
            } )
        }
    }

    // Sort Entries
    sort.Sort( sort.Reverse( _entries ) )

    // Check one by one
    for i := range _entries {
        switch _entries[i].language {
        case "*":
            return _localizer.default_locale
        default:
            // Check if localizer has the language
            _, _ok := _localizer.locales[_entries[i].language]
            if _ok {
                return _entries[i].language
            }
            // If not, try matching langauge only, without the region
            _split := strings.SplitN( _entries[i].language, "-", 2 )
            for key := range _localizer.locales {
                if strings.HasPrefix( key, _split[0] ) {
                    return key
                }
            }
        }
    }

    // If not found
    return ""

}
func( _localizer *Localizer ) Locales() map[string] *Locale {
    _copy := make( map[string] *Locale )
    copy( _copy, _localizer.locales )
    return _copy
}
func( _localizer *Localizer ) SetLocales( _locales map[string] *Locale ) {
    _copy := make( map[string] *Locale )
    copy( _copy, _locales )
    _localizer.locales = _copy
}
func( _localizer *Localizer ) DefaultLocale() string {
    return _localizer.default_locale
}
func( _localizer *Localizer ) SetDefaultLocale( _locale string ) {
    for key := range _localizer.locales {
        if key == _locale {
            _localizer.default_locale = _locale
            return
        }
    }
    // The locale was not found
    panic( LocalizerVariables().ErrorLocaleNotFound() )
}
func( _localizer *Localizer ) SetDelimiters( _left, _right string ) {
    // Check if either is empty
    if len( _left ) == 0 || len( _right ) == 0 {
        panic( LocalizerVariables().ErrorInvalidDelimiter() )
    }
    _localizer.left_delimiter  = _left
    _localizer.right_delimiter = _right
}
func( _localizer *Localizer ) SetThreshold( _threshold int ) {
    // Check if it is below 1
    if _threshold < 1 {
        panic( LocalizerVariables().ErrorInvalidThreshold() )
    }
    _localizer.threshold = _threshold
}