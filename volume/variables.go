package volume

const (

    // CONSTANTS
    c_assetDirectory = "asset"
    c_i18nDirectory = "i18n"
    c_templateDirectory = "template"
    
)

var (

    // FLAGS
    fl_sassInstalled bool = false

    // ERRORS
    ErrUnableToProcessAsset = util.NewError( 500, "unable to process the asset" )
    ErrOccupiedPath = util.NewError( 500, "the asset path is already occupied" )
    ErrInvalidPath = util.NewError( 500, "the given path is invalid" )
    ErrInvalidTemplate = util.NewError( 500, "the given template is invalid" )
    ErrTemplateNonExistent = util.NewError( 500, "there is no template by the given name" )
    ErrTemplateExecError = util.NewError( 500, "the template execution resulted in an error" )
    
)