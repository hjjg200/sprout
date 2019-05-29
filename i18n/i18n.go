package i18n

import (
    "bytes"
    "fmt"
    "io"
    "io/ioutil"
    "net/http"
    "net/url"
    "os"
    "path/filepath"
    "sort"
    "strconv"
    "strings"
)

type I18n struct {
    // Maps in I18n
    //  writes: seldom
    //  reads: frequent
    //  size: small
    //  => create buffer and replace
    locales map[string] *Locale
    localizers map[string] *Localizer
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
        localizers: make( map[string] *Localizer ),
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
    length := len( i1.locales )

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

            plc := i1.locales[lcName]

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
    _, ok := i1.locales[lcName]
    return ok
}

func( i1 *I18n ) NumLocale() int {
    length := len( i1.locales )
    return length
}

func( i1 *I18n ) PutLocale( locale *Locale ) {

    // Set it as the default locale if it is the first locale
    length := len( i1.locales )
    if length == 0 {
        i1.defaultLocale = locale.name
    }

    // Put Locale
    buf := make( map[string] *Locale )
    for k, v := range i1.locales {
        buf[k] = v
    }
    buf[locale.name] = locale
    i1.locales = buf

    // Put Localizer
    buf2 := make( map[string] *Localizer )
    for k, v := range i1.localizers {
        buf2[k] = v
    }
    buf2[locale.name], _ = NewLocalizer( i1, locale.name )
    i1.localizers = buf2

}

func( i1 *I18n ) RemoveLocale( lcName string ) error {

    // Check
    if !i1.HasLocale( lcName ) {
        return fmt.Errorf( "Locale %s was not found", lcName )
    }

    // Buffer
    buf := make( map[string] *Locale )
    for k, v := range i1.locales {
        buf[k] = v
    }
    buf2 := make( map[string] *Localizer )
    for k, v := range i1.localizers {
        buf2[k] = v
    }

    // Delete
    delete( buf, lcName )
    delete( buf2, lcName )

    // Assign
    i1.locales = buf
    i1.localizers = buf2

    return nil
}

func( i1 *I18n ) Locale( lcName string ) ( *Locale, bool ) {
    lc, ok := i1.locales[lcName]
    return lc, ok
}

func( i1 *I18n ) Locales() map[string] *Locale {
    return i1.locales
}

func( i1 *I18n ) LocaleNames() []string {

    names := make( []string, 0 )
    for lcName := range i1.locales {
        names = append( names, lcName )
    }

    return names

}

func( i1 *I18n ) Localizer( lcName string ) ( *Localizer ) {
    lczr := i1.localizers[lcName]
    return lczr
}

func( i1 *I18n ) DefaultLocale() string {
    return i1.defaultLocale
}

func( i1 *I18n ) SetDefaultLocale( lcName string ) error {

    // Check
    if len( lcName ) == 0 {
        return fmt.Errorf( "Given parameter, %s, is invalid", lcName )
    }

    // Check
    _, ok := i1.locales[lcName]
    if !ok {
        return fmt.Errorf( "Given parameter, %s, is invalid", lcName )
    }

    // Assign
    i1.defaultLocale = lcName
    return nil
}

func( i1 *I18n ) SetQueryParameter( param string ) error {
    if len( param ) == 0 {
        return fmt.Errorf( "Given parameter, %s, is invalid", param )
    }
    i1.queryParameter = param
    return nil
}

func( i1 *I18n ) SetCookie( cookie string ) error {
    if len( cookie ) == 0 {
        return fmt.Errorf( "Given parameter, %s, is invalid", cookie )
    }
    i1.cookie = cookie
    return nil
}

func( i1 *I18n ) SetDelimiters( left, right string ) error {
    if len( left ) == 0 || len( right ) == 0 {
        return fmt.Errorf( "Given parameters -- %s and %s -- are invalid", left, right )
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
                return "", fmt.Errorf( "Given accept langauge is invalid: %s; %s", acptLng, err.Error() ) // Malformed accept-language
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
    return "", fmt.Errorf( "Locale for the accept langauge, %s, was not found", acptLng )

}

func( i1 *I18n ) ParseCookies( cks []*http.Cookie ) ( string, error ) {

    // Range cookies
    for i := range cks {
        if cks[i].Name == i1.cookie {
            lcName, err := i1.ParseSingleLocale( cks[i].Value )
            if err != nil {
                break
            }
            return lcName, nil
        }
    }
    return "", fmt.Errorf( "Locale for the cookie value, %s, was not found", cks[i].Value )

}

func( i1 *I18n ) ParseUrlPath( u *url.URL ) ( string, error ) {

    // If too short
    if len( u.Path ) == 1 {
        return "", fmt.Errorf( "Locale for the url path, %s, was not found", u.Path )
    }

    // Parse
    split := strings.SplitN( u.Path[1:], "/", 2 )
    return i1.ParseSingleLocale( split[0] )

}

func( i1 *I18n ) ParseUrlQuery( u *url.URL ) ( string, error ) {

    // Query
    vals   := u.Query()
    lcName := vals.Get( i1.queryParameter )

    // Check
    if lcName == "" {
        return "", fmt.Errorf( "Query parameter was not found" )
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
    _, ok := i1.locales[lcName]
    if ok {
        return lcName, nil
    }

    // Check if the lang exists
    split := strings.SplitN( lcName, "-", 2 )
    if len( split ) >= 1 {
        for i := range i1.locales {
            if strings.HasPrefix( i, split[0] ) {
                return i, nil
            }
        }
    }

    // Return default
    return "", fmt.Errorf( "Locale for %s was not found", lcName )

}
