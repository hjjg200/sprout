package i18n

import (
    "encoding/json"
    "reflect"
    "strings"

    "github.com/hjjg200/sprout/util"
    "github.com/hjjg200/sprout/util/errors"
)

type Locale struct {
    name string
    set map[string] string
    setMx *util.MapMutex
}

func NewLocale() ( *Locale ) {
    return &Locale{
        name: "",
        set: make( map[string] string ),
        setMx: util.NewMapMutex(),
    }
}

func( lc *Locale ) Name() string {
    return lc.name
}

func( lc *Locale ) Set() map[string] string {

    // Copy
    lc.setMx.BeginRead()
    buf := make( map[string] string )
    for k, v := range lc.set {
        buf[k] = v
    }
    lc.setMx.EndRead()

    return buf

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
        return errors.ErrMalformedJson.Raise( err )
    }

    return lc.ParseMap( ifc )

}
func( lc *Locale ) ParseMap( data interface{} ) error {

    // Lock
    lc.setMx.BeginWrite()
    defer lc.setMx.EndWrite()

    // Check the Language
    val := reflect.ValueOf( data )
        // There must be one key under the root
    if len( val.MapKeys() ) > 1 {
        return errors.ErrMalformedJson.Raise( "no language key found" )
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
            return errors.ErrMalformedJson.Raise( "there are no sets in the json" )
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
