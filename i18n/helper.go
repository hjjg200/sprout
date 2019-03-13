package i18n

import (

)

type acceptLangEntry struct{
    language string
    qFactor  float64
}
type acceptLangEntries []acceptLangEntry

func( entries acceptLangEntries ) Len() int { return len( entries ) }
func( entries acceptLangEntries ) Swap( i, j int ) {
    entries[i], entries[j] = entries[j], entries[i]
}
func( entries acceptLangEntries ) Less( i, j int ) bool {
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