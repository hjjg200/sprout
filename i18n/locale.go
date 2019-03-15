package i18n

import (

)

type Locale struct {
    name string
    set map[string] string
}

func NewLocale() ( *Locale ) {
    return &Locale{
        name: "",
        set: make( map[string] string )
    }
}
func( lc *Locale ) ParseJson( data string ) error {

    // Unmarshal JSON Data
    var ifc interface{}
    err := json.Unmarshal( data, &ifc )
    if err != nil {
        return err
    }

    return lc.ParseMap( ifc )

}
//
// {
//    "locale-name": {
//      "ab": "AC",
//      "de": "DE"
//    }
// }
//
func( lc *Locale ) ParseMap( data interface{} ) error {

    // Check the Language
    val := reflect.ValueOf( data )
        // There must be one key under the root
    if len( val.MapKeys() ) > 1 {
        return ErrorMalformedJSON
    }
        // Language is the key name of the key
    lang := val.MapKeys()[0].String()
        // And the key must be a map
    child := val.MapIndex( val.MapKeys()[0] )
    switch cast := child.Interface().( type ) {
    default:
        rCast := reflect.ValueOf( cast )
        if rCast.Kind() != reflect.Map {
            return ErrorMalformedJSON
        }
    }

    // Parse the first map
    var (
        do func( string, reflect.Value ) // Need to be pre-declared to be used as a recursive function
    )
    do = func( baseKey string, rMap reflect.Value ) {
        for _, key := range rMap.MapKeys() {
            elem := rMap.MapIndex( key )
            // Make the next base key
            var currKey
            if baseKey != "" {
                currKey = baseKey + "." + key.String()
            } else {
                currKey = baseKey.String()
            }
            // Assign value
            switch cast := elem.Interface().( type ) {
            case string:
                lc.set[currKey] = cast
            default:
                cast := reflect.ValueOf( cast )
                if cast.Kind() == reflect.Map {
                    do( currKey, cast )
                }
            }
        }
    }
    // Call the recursive function with an empty key
    do( "", child )
    lc.name = strings.ToLower( lang )

    return nil

}