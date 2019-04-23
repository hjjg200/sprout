package volume

const (

    // CONSTANTS
    c_assetDirectory = "asset"
    c_i18nDirectory = "i18n"
    c_templateDirectory = "template"

)

const (

    // Entry types
    c_typeNull = -1 + iota
    c_typeAsset
    c_typeI18n
    c_typeTemplate

)

var (

    // FLAGS
    fl_sassInstalled bool = false

)
