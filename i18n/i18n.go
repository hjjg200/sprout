package i18n

import (
    "bytes"
    "io"
    "io/ioutil"
    "os"
    "sort"
    "strconv"
    "strings"
)

type I18n struct {
    locales map[string] *Locale
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
        i1.AddLocale( lc )
        
    }
    
    return nil
    
}
func( i1 *I18n ) L( locale, src string ) string {
    return i1.Localize( locale, src )
}
func( i1 *I18n ) Localize( locale, src string ) string {

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
            plc, localeOk := i1.locales[locale]
            if !localeOk {
                plc = i1.locales[i1.defaultLocale]
            }
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
func( i1 *I18n ) AddLocale( locale *Locale ) {
    // Set it as the default locale if it is the first locale
    if len( i1.locales ) == 0 {
        i1.defaultLocale = locale.name
    }
    i1.locales[locale.name] = locale
}
func( i1 *I18n ) SetDefaultLocale( lcName string ) error {
    if len( lcName ) == 0 {
        return ErrInvalidParameter
    }
    if _, ok := i1.locales[lcName]; !ok {
        return ErrLocaleNonExistent
    }
    i1.defaultLocale = lcName
    return nil
}
func( i1 *I18n ) SetQueryParameter( param string ) error {
    if len( param ) == 0 {
        return ErrInvalidParameter
    }
    i1.queryParameter = param
    return nil
}
func( i1 *I18n ) SetCookie( cookie string ) error {
    if len( cookie ) == 0 {
        return ErrInvalidParameter
    }
    i1.cookie = cookie
    return nil
}
func( i1 *I18n ) SetDelimiters( left, right string ) error {
    if len( left ) == 0 || len( right ) == 0 {
        return ErrInvalidDelimiters
    }
    i1.leftDelimiter = left
    i1.rightDelimiter = right
    return nil
}

/*
 + ParseAcceptLangauge
 *
 * Retreives a suitable locale for localization from http accept language string
 */

func( i1 *I18n ) ParseAcceptLanguage( acptLng string ) string {

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
                panic( err ) // Malformed accept-language
                return i1.defaultLocale
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
            return i1.defaultLocale
        default:
            // Check if localizer has the language
            _, ok := i1.locales[entries[i].locale]
            if ok {
                return entries[i].locale
            }
            // If not, try matching langauge only, without the region
            split := strings.SplitN( entries[i].locale, "-", 2 )
            for key := range i1.locales {
                if strings.HasPrefix( key, split[0] ) {
                    return key
                }
            }
        }
    }

    // If not found
    return i1.defaultLocale

}