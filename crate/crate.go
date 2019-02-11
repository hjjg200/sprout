package crate

type Crate struct {
    assets    map[string] *Asset
    localizer localizer.Localizer
    templates *template.Template
    mod_time  map[string] time.Time

    default_locale string
}

const (
    envDirAsset = "asset"
    envDirLocale = "locale"
    envDirTemplate = "template"
)

/*
funcs to create

MakeArchive
FromArchive
FromDirectory
Asset( string ) *Asset, bool
Template( string ) *Template, bool
Localize( string, string, int ) string, error

*/