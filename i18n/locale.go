package i18n

import (
    "encoding/json"
    "reflect"
    "strings"
)

type Locale struct {
    name string
    set map[string] string
}

func NewLocale() ( *Locale ) {
    return &Locale{
        name: "",
        set: make( map[string] string ),
    }
}

func( lc *Locale ) Name() string {
    return lc.name
}
//
// {
//    "locale-name": {
//      "ab": "AC",
//      "de": "DE"
//    }
// }
//
func( lc *Locale ) ParseJson( data []byte ) error {

    // Unmarshal JSON Data
    var ifc interface{}
    err := json.Unmarshal( data, &ifc )
    if err != nil {
        return err
    }

    return lc.ParseMap( ifc )

}
func( lc *Locale ) ParseMap( data interface{} ) error {

    // Check the Language
    val := reflect.ValueOf( data )
        // There must be one key under the root
    if len( val.MapKeys() ) > 1 {
        return ErrMalformedJson
    }
        // Language is the key name of the key
    lang := val.MapKeys()[0].String()
        // And the key must be a map
    child := val.MapIndex( val.MapKeys()[0] )
    var rChild reflect.Value
    switch cast := child.Interface().( type ) {
    default:
        rChild = reflect.ValueOf( cast )
        if rChild.Kind() != reflect.Map {
            return ErrMalformedJson
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
            var currKey string
            if baseKey != "" {
                currKey = baseKey + "." + key.String()
            } else {
                currKey = key.String()
            }
            // Assign value
            switch cast := elem.Interface().( type ) {
            case string:
                lc.set[currKey] = cast
            default:
                rCast := reflect.ValueOf( cast )
                if rCast.Kind() == reflect.Map {
                    do( currKey, rCast )
                }
            }
        }
    }
    // Call the recursive function with an empty key
    do( "", rChild )
    lc.name = strings.ToLower( lang )

    return nil

}