package errors

var (

    ErrAssetProcessingFailure = newType( "ErrAssetProcessingFailure" )
    ErrOccupiedPath = newType( "ErrOccupiedPath" )
    ErrInvalidPath = newType( "ErrInvalidPath" )
    ErrNotFound = newType( "ErrNotFound" )
    ErrInvalidObject = newType( "ErrInvalidObject" )
    ErrTemplateExecFailure = newType( "ErrTemplateExecFailure" )
    ErrIOError = newType( "ErrIOError" )
    ErrImportFailure = newType( "ErrImportFailure" )
    ErrExportFailure = newType( "ErrExportFailure" )
    ErrCompileFailure = newType( "ErrCompileFailure" )
    ErrInvalidParameter = newType( "ErrInvalidParameter" )
    ErrMalformedJson = newType( "ErrMalformedJson" )
    ErrUrlHasNoLocale = newType( "ErrUrlHasNoLocale" )
    ErrServerOperation = newType( "ErrServerOperation" )
    ErrServerExited = newType( "ErrServerExited" )
    ErrDifferentStatusCode = newType( "ErrDifferentStatusCode" )
    ErrOSNotSupported = newType( "ErrOSNotSupported" )
    ErrMalformedAcceptLang = newType( "ErrMalformedAcceptLang" )
    ErrReverseProxy = newType( "ErrReverseProxy" )
    ErrCmdExecFailure = newType( "ErrCmdExecFailure" )

)
