package errors

var (

    // VOLUME

    ErrUnableToProcessAsset = newType( "ErrUnableToProcessAsset" )
    ErrOccupiedPath = newType( "ErrOccupiedPath" )
    ErrInvalidPath = newType( "ErrInvalidPath" )
    ErrInvalidTemplate = newType( "ErrInvalidTemplate" )
    ErrTemplateNonExistent = newType( "ErrTemplateNonExistent" )
    ErrTemplateExecError = newType( "ErrTemplateExecError" )
    ErrFileError = newType( "ErrFileError" )
    ErrZipImport = newType( "ErrZipImport" )
    ErrZipExport = newType( "ErrZipExport" )
    ErrAssetExport = newType( "ErrAssetExport" )
    ErrI18nExport = newType( "ErrI18nExport" )
    ErrTemplateExport = newType( "ErrTemplateExport" )
    ErrCompileFailure = newType( "ErrCompileFailure" )
    ErrDirectoryError = newType( "ErrDirectoryError" )
    ErrItemNonExistent = newType( "ErrItemNonExistent" )
    ErrPathNonExistent = newType( "ErrPathNonExistent" )

    // I18N

    ErrInvalidLocale = newType( "ErrInvalidLocale" )
    ErrInvalidDelimiters = newType( "ErrInvalidDelimiters" )
    ErrInvalidThreshold = newType( "ErrInvalidThreshold" )
    ErrInvalidParameter = newType( "ErrInvalidParameter" )
    ErrMalformedJson = newType( "ErrMalformedJson" )
    ErrLocaleNonExistent = newType( "ErrLocaleNonExistent" )
    ErrLocaleExists = newType( "ErrLocaleExists" )
    ErrCookieNonExistent = newType( "ErrCookieNonExistent" )
    ErrQueryNonExistent = newType( "ErrQueryNonExistent" )
    ErrUrlHasNoLocale = newType( "ErrUrlHasNoLocale" )

    // CACHE

    ErrEntryAccessFailed = newType( "ErrEntryAccessFailed" )
    ErrEntryNotFound = newType( "ErrEntryNotFound" )
    ErrEntryWriteFail = newType( "ErrEntryWriteFail" )

    // NETWORK

    ErrStartingServer = newType( "ErrStartingServer" )
    ErrStoppingServer = newType( "ErrStoppingServer" )
    ErrServerExited = newType( "ErrServerExited" )
    ErrRequestClosed = newType( "ErrRequestClosed" )
    ErrDifferentStatusCode = newType( "ErrDifferentStatusCode" )

    // SYSTEM

    ErrOSNotSupported = newType( "ErrOSNotSupported" )

)