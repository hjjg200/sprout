package volume

import (
    "../util"
)

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

    // ERRORS
    ErrUnableToProcessAsset = util.NewError( 500, "unable to process the asset" )
    ErrOccupiedPath = util.NewError( 500, "the asset path is already occupied" )
    ErrInvalidPath = util.NewError( 500, "the given path is invalid" )
    ErrInvalidTemplate = util.NewError( 500, "the given template is invalid" )
    ErrTemplateNonExistent = util.NewError( 500, "there is no template by the given name" )
    ErrTemplateExecError = util.NewError( 500, "the template execution resulted in an error" )
    ErrFileError = util.NewError( 500, "was not able to access the file" )
    ErrZipImport = util.NewError( 500, "was not able to import the zip file" )
    ErrZipExport = util.NewError( 500, "unable to export the zip" )
    ErrAssetExport = util.NewError( 500, "unable to export the asset" )
    ErrI18nExport = util.NewError( 500, "unable to export the i18n instance" )
    ErrTemplateExport = util.NewError( 500, "unable to export the template" )
    ErrCompileFailure = util.NewError( 500, "unable to compile the given source" )
    ErrDirectoryError = util.NewError( 500, "unable to access the directory" )

)