package sprout

import (
    "encoding/json"
    "io"
    "io/ioutil"
    "reflect"
    "strings"
)

// localizer is private type

type localizer struct {
    locales map[string] map[string] string
}

func ( sp *Sprout ) newLocalizer() ( *localizer ) {

    lc         := &localizer{}
    lc.locales  = make( map[string] map[string] string )
/*
    _locale_dir       := sp.cwd + "/" + envDirLocale
    _file_infos, _err := ioutil.ReadDir( _locale_dir )
    if _err != nil {
        return nil, _err
    }

    for _, _fi := range _file_infos {

        _base := _fi.Name()
        _ext  := path.Ext( _base )
        if _ext == ".json" {

            _bytes, _err := ioutil.ReadFile( _locale_dir + "/" + _base )
            if _err != nil {
                continue
            }

            var _json_interface interface{}
            _err         = json.Unmarshal( _bytes, &_json_interface )
            if _err != nil {
                continue
            }

            var _map_func func ( string, reflect.Value )
            _locale_map := make( map[string] string )
            _map_func    = func ( __base_key string, __map reflect.Value ) {
                for _, __k := range __map.MapKeys() {

                    __it       := __map.MapIndex( __k )
                    __next_key := __base_key
                    if __next_key != "" {
                        __next_key = __next_key + "." + __k.String()
                    } else {
                        __next_key = __k.String()
                    }

                    switch __v := __it.Interface().(type) {
                    case string:
                        _locale_map[__next_key] = __v
                    default:
                        __rv := reflect.ValueOf( __v )
                        if __rv.Kind() == reflect.Map {
                            _map_func( __next_key, __rv )
                        }
                    }
                }
            }

            _map_func( "", reflect.ValueOf( _json_interface ) )
            _locale_str := _base[:len( _base ) - len( _ext )]

            lc.locales[_locale_str] = make( map[string] string )
            for _k, _v := range _locale_map {
                lc.locales[_locale_str][_k] = _v
            }

        }

    }
*/
    return lc

}

func ( lc *localizer ) removeLocale( _locale string ) error {
    if _, _ok := lc.locales[_locale]; _ok {
        delete( lc.locales, _locale )
        return nil
    } else {
        return ErrInvalidLocale
    }
}

func ( lc *localizer ) appendLocale( _locale string, _json []byte ) error {

    var _json_interface interface{}
    _err := json.Unmarshal( _json, &_json_interface )
    if _err != nil {
        return _err
    }

    var _map_func func ( string, reflect.Value )
    _locale_map := make( map[string] string )
    _map_func    = func ( __base_key string, __map reflect.Value ) {
        for _, __k := range __map.MapKeys() {

            __it       := __map.MapIndex( __k )
            __next_key := __base_key
            if __next_key != "" {
                __next_key = __next_key + "." + __k.String()
            } else {
                __next_key = __k.String()
            }

            switch __v := __it.Interface().(type) {
            case string:
                _locale_map[__next_key] = __v
            default:
                __rv := reflect.ValueOf( __v )
                if __rv.Kind() == reflect.Map {
                    _map_func( __next_key, __rv )
                }
            }
        }
    }

    _map_func( "", reflect.ValueOf( _json_interface ) )

    if _, _ok := lc.locales[_locale]; !_ok {
        lc.locales[_locale] = make( map[string] string )
    }

    for _k, _v := range _locale_map {
        lc.locales[_locale][_k] = _v
    }

    return nil

}

func ( lc *localizer ) hasLocale( locale string ) bool {
    _, ok := lc.locales[locale]
    return ok
}

func ( lc *localizer ) localize_reader( _r io.Reader, locale string, threshold int ) ( io.Reader, error ) {

    _bytes, _err := ioutil.ReadAll( _r )
    if _err != nil {
        return _r, _err
    }

    _string, _err := lc.localize( string( _bytes ), locale, threshold )
    if _err != nil {
        return _r, _err
    }

    return strings.NewReader( _string ), nil

}

func ( lc *localizer ) localize( src string, locale string, threshold int ) ( string, error ) {

    _do := func() {
        var (
            _last_char  = rune( 0 )
            _last_key   = ""
            _read       = false
            _found_keys = make( []string, 0 )
        )

        for _i, _c := range src {
            switch _c {
            case '%':
                if _last_char == '{' {
                    _read = true
                } else {
                    if _i + 1 <= len( src ) - 1 {
                        if src[_i + 1] == '}' {
                            _found_keys = append( _found_keys, _last_key )
                            _last_key   = ""
                            _read       = false
                        }
                    }
                }
            default:
                if _read {
                    _last_key += string( _c )
                }
            }
            _last_char = _c
        }

        for _, _key := range _found_keys {
            _v, _ok := lc.locales[locale][_key]
            if _ok {
                src = strings.Replace( src, "{%" + _key + "%}", _v, -1 )
            }
        }
    }

    // Threshold
    if threshold < 0 {
        threshold = default_localizer_threshold
    }

    for _cycles := 0; _cycles < threshold; _cycles++ {
        if strings.Contains( src, "{%" ) /* && strings.Contains( src, "%}" ) */ {
            _do()
        } else {
            break
        }
    }

    return src, nil
}