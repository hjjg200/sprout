package i18n

import (

)

type acceptLanguageEntry struct{
    locale  string
    qFactor float64
}
type acceptLanguageEntries []acceptLanguageEntry

func( entries acceptLanguageEntries ) Len() int { return len( entries ) }
func( entries acceptLanguageEntries ) Swap( i, j int ) {
    entries[i], entries[j] = entries[j], entries[i]
}
func( entries acceptLanguageEntries ) Less( i, j int ) bool {
    return entries[i].qFactor < entries[j].qFactor
}

func IsValidLocaleName( name string ) bool {

    // https://tools.ietf.org/html/rfc3066#page-2

    var (
        alphanum = false
        priLen = 0 // Length of primary tag
        subLen = 0 // Length of sub tag
    )
    for _, c := range name {
        switch {
        case c >= 'A' && c <= 'Z',
             c >= 'a' && c <= 'z',
             c >= '0' && c <= '9':
            if !alphanum {
                if c >= '0' && c <= '9' {
                    return false
                }
                subLen++
            } else {
                priLen++
            }
        case c == '-':
            if alphanum {
                return false
            }
            alphanum = true
        default:
            return false
        }
    }
    if priLen > 8 || subLen > 8 {
        return false
    }
    return true
}

func strContainsAt( src, sub string, i int ) bool {
    
    // If the rest string is shorter than the substring
    if i + len( sub ) > len( src ) {
        return false
    }
    // Check the first letter
    if src[i] != sub[0] {
        return false
    }
    // Check the whole substring
    return src[i:i + len( sub )] == sub
    
}