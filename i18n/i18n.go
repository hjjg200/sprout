package i18n

import (
    "bytes"
    "path/filepath"
    "io"
    "io/ioutil"
    "net/http"
    "net/url"
    "os"
    "sort"
    "strconv"
    "strings"

    "../util"
)

type I18n struct {
    locales map[string] *Locale
    localesMx *util.MapMutex
    localizers map[string] *Localizer
    localizersMx *util.MapMutex
    defaultLocale string
    cookie string
    queryParameter string
    threshold int
    leftDelimiter string
    rightDelimiter string
}

func New() *I18n {
    return &I18n{
        locales: make( map[string] *Locale ),
        localesMx: util.NewMapMutex(),
        localizers: make( map[string] *Localizer ),
        localizersMx: util.NewMapMutex(),
        defaultLocale: "",
        cookie: c_defaultCookie,
        queryParameter: c_defaultQueryParameter,
        threshold: c_defaultThreshold,
        leftDelimiter: c_defaultLeftDelimiter,
        rightDelimiter: c_defaultRightDelimiter,
    }
}

func( i1 *I18n ) ImportDirectory( path string ) error {

    // Read the dir
    fis, err := ioutil.ReadDir( path )
    if err != nil {
        return err
    }

    // Foreach
    for _, fi := range fis {

        // Check ext
        ext := filepath.Ext( fi.Name() )
        if ext != ".json" {
            break
        }

        // Open
        f, err := os.Open( path + "/" + fi.Name() )
        if err != nil {
            return err
        }

        // Read
        buf := bytes.NewBuffer( nil )
        io.Copy( buf, f )
        f.Close()

        // Parse
        lc  := NewLocale()
        err  = lc.ParseJson( buf.Bytes() )
        if err != nil {
            return err
        }

        // Assign
        i1.PutLocale( lc )

    }

    return nil

}
func( i1 *I18n ) L( lcName, src string ) string {
    return i1.Localize( lcName, src )
}
func( i1 *I18n ) Localize( lcName, src string ) string {

    // If 0 locale
    i1.localesMx.BeginRead()
    length := len( i1.locales )
    i1.localesMx.EndRead()

    if length == 0 {
        return src
    }

    // Do
    for i := 0; i < i1.threshold; i++ {

        // Check if there is any left delimiter
        if strings.Contains( src, i1.leftDelimiter ) == false {
            break
        }

        var (
            key       = ""
            read      = false
            srcLen    = len( src )
            lastBegin = 0
            foundKeys = make( map[string] struct{} )
            left      = i1.leftDelimiter
            right     = i1.rightDelimiter
        )

        for x := 0; x < srcLen; {
            switch {
            case strContainsAt( src, left, x ):
                key       = ""
                read      = true
                lastBegin = x
                x += len( left )
            case strContainsAt( src, right, x ):

                // Trim whitespace
                key = strings.TrimSpace( key )
                foundKeys[key] = struct{}{}

                // Remove whitespaces near delimiters
                newBlock := left + key + right
                blockLen := x + len( right ) - 1 - lastBegin + 1
                rest     := ""
                if x + len( right ) < srcLen {
                    rest = src[x + len( right ):]
                }
                src     = src[:lastBegin] + newBlock + rest
                diff   := blockLen - len( newBlock )
                srcLen -= diff
                x      -= diff

                // Reset
                key  = ""
                read = false
                x += len( right )

            default:
                if read {
                    key += string( src[x] )
                }
                x++
            }
        }

        // Replace the found keys
        for key := range foundKeys {
            lcName, err := i1.ParseSingleLocale( lcName )
            if err != nil {
                lcName = i1.defaultLocale
            }

            i1.localesMx.BeginRead()
            plc := i1.locales[lcName]
            i1.localesMx.EndRead()

            if val, ok := plc.set[key]; ok {
                src = strings.Replace(
                    src,
                    left + key + right,
                    val,
                    -1,
                )
            }
        }

    }

    return src

}

func( i1 *I18n ) HasLocale( lcName string ) bool {
    i1.localesMx.BeginRead()
    defer i1.localesMx.EndRead()
    for i := range i1.locales {
        if i == lcName {
            return true
        }
    }
    return false
}

func( i1 *I18n ) NumLocale() int {
    i1.localesMx.BeginRead()
    length := len( i1.locales )
    i1.localesMx.EndRead()
    return length
}

func( i1 *I18n ) PutLocale( locale *Locale ) {

    // Set it as the default locale if it is the first locale
    i1.localesMx.BeginRead()
    length := len( i1.locales )
    i1.localesMx.EndRead()
    if length == 0 {
        i1.defaultLocale = locale.name
    }

    // Put Locale
    i1.localesMx.BeginWrite()
    i1.locales[locale.name] = locale
    i1.localesMx.EndWrite()

    // Put Localizer
    i1.localizersMx.BeginWrite()
    i1.localizers[locale.name], _ = NewLocalizer( i1, locale.name )
    i1.localizersMx.EndWrite()

}

func( i1 *I18n ) Locale( lcName string ) ( *Locale, bool ) {
    i1.localesMx.BeginRead()
    lc, ok := i1.locales[lcName]
    i1.localesMx.EndRead()
    return lc, ok
}

func( i1 *I18n ) LocaleNames() []string {

    i1.localesMx.BeginRead()
    names := make( []string, 0 )
    for lcName := range i1.locales {
        names = append( names, lcName )
    }
    i1.localesMx.EndRead()

    return names

}

func( i1 *I18n ) Localizer( lcName string ) ( *Localizer ) {
    i1.localizersMx.BeginRead()
    lczr := i1.localizers[lcName]
    i1.localizersMx.EndRead()
    return lczr
}

func( i1 *I18n ) DefaultLocale() string {
    return i1.defaultLocale
}

func( i1 *I18n ) SetDefaultLocale( lcName string ) error {

    //
    if len( lcName ) == 0 {
        return ErrInvalidParameter.Append( lcName )
    }

    //
    i1.localesMx.BeginRead()
    _, ok := i1.locales[lcName]
    i1.localesMx.EndRead()
    if !ok {
        return ErrLocaleNonExistent.Append( lcName )
    }

    //
    i1.defaultLocale = lcName
    return nil
}

func( i1 *I18n ) SetQueryParameter( param string ) error {
    if len( param ) == 0 {
        return ErrInvalidParameter.Append( param )
    }
    i1.queryParameter = param
    return nil
}

func( i1 *I18n ) SetCookie( cookie string ) error {
    if len( cookie ) == 0 {
        return ErrInvalidParameter.Append( cookie )
    }
    i1.cookie = cookie
    return nil
}

func( i1 *I18n ) SetDelimiters( left, right string ) error {
    if len( left ) == 0 || len( right ) == 0 {
        return ErrInvalidDelimiters.Append( left, right )
    }
    i1.leftDelimiter = left
    i1.rightDelimiter = right
    return nil
}

func( i1 *I18n ) MakeCookie( lcName string ) *http.Cookie {
    return &http.Cookie{
        Name: i1.cookie,
        Value: lcName,
        Path: "/", // for every page
        MaxAge: 0, // persistent cookie
    }
}

/*
 + ParseAcceptLangauge
 *
 * Retreives a suitable locale for localization from http accept language string
 */

func( i1 *I18n ) ParseAcceptLanguage( acptLng string ) ( string, error ) {

    // Split the header
    entries := make( acceptLanguageEntries, 0 )
    split   := strings.Split( acptLng, "," )
    for i := range split {
        // Remove whitespaces and to lowercase
        split[i] = strings.TrimSpace( split[i] )
        split[i] = strings.ToLower( split[i] )

        if semicolon := strings.Index( split[i], ";" ); semicolon != -1 {
            // If there is the q-factor
            lcName       := split[i][:semicolon]
            qFactor, err := strconv.ParseFloat( split[i][semicolon + 3:], 64 )
            if err != nil {
                //panic( err )
                return "", err // Malformed accept-language
            }
            entries = append( entries, acceptLanguageEntry{
                locale: lcName,
                qFactor: qFactor,
            } )
        } else {
            // Since its q-factor is default which is 1.0, it is prepended
            entries = append( entries, acceptLanguageEntry{
                locale: split[i],
                qFactor: 1.0,
            } )
        }
    }

    // Sort Entries
    sort.Sort( sort.Reverse( entries ) )

    // Check one by one
    for i := range entries {
        switch entries[i].locale {
        case "*":
            // If wildcard, return default
            return i1.defaultLocale, nil
        default:
            // Check if localizer has the language
            lcName, err := i1.ParseSingleLocale( entries[i].locale )
            if err == nil {
                return lcName, nil
            }
        }
    }

    // If not found
    return "", ErrLocaleNonExistent.Append( acptLng )

}

func( i1 *I18n ) ParseCookies( cks []*http.Cookie ) ( string, error ) {

    // Range cookies
    for i := range cks {
        if cks[i].Name == i1.cookie {
            lcName, err := i1.ParseSingleLocale( cks[i].Value )
            if err != nil {
                return "", ErrLocaleNonExistent.Append( cks[i].Value )
            }
            return lcName, nil
        }
    }
    return "", ErrCookieNonExistent

}

func( i1 *I18n ) ParseUrlPath( u *url.URL ) ( string, error ) {

    // If too short
    if len( u.Path ) == 1 {
        return "", ErrUrlHasNoLocale
    }

    // Parse
    split := strings.SplitN( u.Path[1:], "/", 2 )
    return i1.ParseSingleLocale( split[0] )

}

func( i1 *I18n ) ParseUrlQuery( u *url.URL ) ( string, error ) {

    // Query
    vals   := u.Query()
    lcName := vals.Get( i1.queryParameter )

    //
    if lcName == "" {
        return "", ErrQueryNonExistent
    } else {
        lcName, err := i1.ParseSingleLocale( lcName )
        if err != nil {
            return "", err
        }
        return lcName, nil
    }

}

func( i1 *I18n ) ParseSingleLocale( lcName string ) ( string, error ) {

    // Check if exists
    i1.localesMx.BeginRead()
    _, ok := i1.locales[lcName]
    i1.localesMx.EndRead()
    if ok {
        return lcName, nil
    }

    // Check if the lang exists
    split := strings.SplitN( lcName, "-", 2 )
    if len( split ) >= 1 {
        i1.localesMx.BeginRead()
        defer i1.localesMx.EndRead()
        for i := range i1.locales {
            if strings.HasPrefix( i, split[0] ) {
                return i, nil
            }
        }
    }

    // Return default
    return "", ErrLocaleNonExistent

}