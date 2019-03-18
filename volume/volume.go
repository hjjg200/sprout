package volume

type Volume struct {
    assets map[string] *Asset
    i18n *i18n.I18n
    templates []*template.Template
}

func NewVolume() *Volume {
    return &Volume{
        assets: make( map[string] *Asset ),
        i18n: i18n.New(),
        templates: template.New( "" ),
    }
}

func PathHasPrefix( path, prefix string ) bool {
    _, err := filepath.Rel( prefix, path )
    return err == nil
}

// Getters

func( vol *Volume ) Asset( path string ) ( *Asset, bool ) {
    return vol.assets[path]
}
func( vol *Volume ) Localizer( lcName string ) ( *i18n.Localizer, bool ) {
    lczr, err := i18n.NewLocalizer( vol.i18n, lcName )
    if err != nil {
        return nil, false
    }
    return lczr, true
}
func( vol *Volume ) Template( path string ) ( *template.Template, bool ) {
    tmpl := vol.templates.Lookup( path )
    if tmpl == nil {
        return nil, false
    }
    return tmpl, true
}

// General

func( vol *Volume ) Localize( lcName, src string ) string {
    return vol.i18n.Localize( lcName, src )
}
func( vol *Volume ) ExecuteTemplate( path string, data interface{} ) ( string, error ) {

    // Get
    tmpl, ok := vol.Template( path )
    if !ok {
        return "", ErrTemplateNonExistent.Append( path )
    }

    // Read
    buf := bytes.NewBuffer()
    err := tmpl.Execute( buf, data )
    if err != nil {
        return "", ErrTemplateExecError.Append( err )
    }

    return buf.String(), nil

}

func( vol *Volume ) Reset() {
    vol.assets = make( map[string] *Asset )
    vol.i18n = i18n.New()
    vol.templates = template.New( "" )
}

// Importers

func( vol *Volume ) PutItem( path string, rd io.Reader, modTime time.Time ) error {

    switch {
    case PathHasPrefix( path, c_assetDirectory ):
        ast := NewAsset( filepath.Base( path ), rd, modTime )
        return vol.PutAsset( path, ast )
    case PathHasPrefix( path, c_localeDirectory ):

        // Read
        buf := bytes.NewBuffer()
        io.Copy( buf, rd )

        // Parse
        lc := i18n.NewLocale()
        err := lc.ParseJson( buf.Bytes() )
        if err != nil {
            return err
        }

        // Assign
        return vol.PutLocale( lc )

    case PathHasPrefix( path, c_templateDirectory ):

        // Read
        buf := bytes.NewBuffer()
        io.Copy( buf, rd )

        // Add
        return vol.PutTemplate( path, buf.String() )

    default:
        return ErrInvalidPath.Append( path )
    }
    return nil

}

func( vol *Volume ) PutAsset( path string, ast *Asset ) error {
    if _, ok := vol.assets[path]; ok {
        return ErrOccupiedPath.Append( path )
    }
    if !PathHasPrefix( path, c_assetDirectory ) {
        return ErrInvalidPath.Append( path )
    }
    vol.assets[path] = ast
    return nil
}

func( vol *Volume ) PutLocale( lc *i18n.Locale ) error {
    return vol.i18n.AddLocale( lc )
}

func( vol *Volume ) PutTemplate( path string, text string ) error {
    if _, err := vol.templates.Lookup( path ); err == nil {
        return ErrOccupiedPath.Append( path )
    }
    if !PathHasPrefix( path, c_templateDirectory ) {
        return ErrInvalidPath.Append( path )
    }
    _, err := vol.templates.New( path ).Parse( text )
    if err != nil {
        ErrInvalidTemplate.Append( path, err )
    }
    return nil
}

// Filesystem Related
//
// * Notice how the methods use osPath instaed of path to differentiate paths in volume and paths in operating system.

func( vol *Volume ) ImportLocaleDirectory( osPath string ) error {
    return vol.i18n.ImportDirectory( osPath )
}

/*
 + WalkFuncBasedOn
 *
 * filepath.WalkFunc
 * relpath must be volume-path-ready e.g.) templates/
 */

func( vol *Volume ) WalkFuncBasedOn( basePath string ) filepath.WalkFunc {

    return func( osPath string, fi os.FileInfo, err error ) error {

        // Rel
        relPath, relErr := filepath.Rel( basePath, osPath )
        if relErr != nil {
            return ErrInvalidPath.Append( relErr, "basePath:", basePath, "osPath:", osPath )
        }

        // Open
        f, err := os.Open( relPath )
        if err != nil {
            return ErrFileError.Append( relPath )
        }

        // Add
        err = vol.PutItem( relPath, f, fi.ModTime() )
        if err != nil {
            return err
        }

        // Close
        err = f.Close()
        if err != nil {
            return ErrFileError.Append( err )
        }

        return nil

    }

}

func( vol *Volume ) ImportZip( zr *zip.Reader ) error {

    for fh := range zr.File {

        // Open
        f, err := fh.Open()
        if err != nil {
            return ErrZipImport.Append( fh.Name, err )
        }

        // Put
        err = vol.PutItem( fh.Name, f, fh.Modified )
        if err != nil {
            return ErrZipImport.Append( err )
        }

    }

}

func( vol *Volume ) ExportZip( wr io.Writer ) error {

    zwr := zip.NewWriter( wr )
    exportFunc := func( exps []Exporter ) error {
        for _, exp := range exps {

            // Export
            err := exp.Export( zwr )
            if err != nil {
                return ErrZipExport.Append( err )
            }

            return nil

        }
    }

    // Assets
    err := exportFunc( vol.assets )
    if err != nil {
        return err
    }

    // Locale
    i1exp := i18nExporter{ vol.i18n }
    err := exportFunc( []Exporter{ i1exp } )
    if err != nil {
        return err
    }

    // Templats
    err := exportFunc( vol.templates )
    if err != nil {
        return err
    }

    zwr.Close()
    return nil

}