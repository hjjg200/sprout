package volume

import (
    "archive/zip"
    "bytes"
    "html/template"
    "io"
    "os"
    "path/filepath"
    "strings"
    "time"

    "../i18n"
    "../util"
)

type Volume struct {
    assets map[string] *Asset
    i18n *i18n.I18n
    rootTemplate *template.Template
    templates map[string] *Template
}

func NewVolume() *Volume {
    vol := &Volume{}
    vol.Reset()
    return vol
}

// Getters

func( vol *Volume ) Asset( path string ) ( *Asset, bool ) {
    ast, ok := vol.assets[path]
    return ast, ok
}
func( vol *Volume ) Localizer( lcName string ) ( *i18n.Localizer, bool ) {
    lczr, err := i18n.NewLocalizer( vol.i18n, lcName )
    if err != nil {
        return nil, false
    }
    return lczr, true
}
func( vol *Volume ) Template( path string ) ( *Template, bool ) {
    tmpl, ok := vol.templates[path]
    return tmpl, ok
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
    buf := bytes.NewBuffer( nil )
    err := tmpl.Execute( buf, data )
    if err != nil {
        return "", ErrTemplateExecError.Append( err )
    }

    return buf.String(), nil

}

func( vol *Volume ) Reset() {
    rootTemplate := template.New( "" )

    vol.assets = make( map[string] *Asset )
    vol.i18n = i18n.New()
    vol.rootTemplate = rootTemplate
    vol.templates = map[string] *Template{
        "": &Template{
            Template: rootTemplate,
            text: "",
        },
    }
}

// Importers

func( vol *Volume ) PutItem( path string, rd io.Reader, modTime time.Time ) error {

    switch {
    case strings.HasPrefix( path, c_assetDirectory ):
        ast := NewAsset( filepath.Base( path ), rd, modTime )
        return vol.PutAsset( path, ast )
    case strings.HasPrefix( path, c_i18nDirectory ):

        // Read
        buf := bytes.NewBuffer( nil )
        io.Copy( buf, rd )

        // Parse
        lc := i18n.NewLocale()
        err := lc.ParseJson( buf.Bytes() )
        if err != nil {
            return err
        }

        // Assign
        return vol.PutLocale( lc )

    case strings.HasPrefix( path, c_templateDirectory ):

        // Read
        buf := bytes.NewBuffer( nil )
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
    if !strings.HasPrefix( path, c_assetDirectory ) {
        return ErrInvalidPath.Append( path )
    }
    vol.assets[path] = ast
    return nil
}

func( vol *Volume ) PutLocale( lc *i18n.Locale ) error {
    return vol.i18n.AddLocale( lc )
}

func( vol *Volume ) PutTemplate( path string, text string ) error {

    // Check if exists
    if tmpl := vol.rootTemplate.Lookup( path ); tmpl != nil {
        return ErrOccupiedPath.Append( path )
    }

    // Check path
    if !strings.HasPrefix( path, c_templateDirectory ) {
        return ErrInvalidPath.Append( path )
    }

    // Parse
    tmpl, err := vol.rootTemplate.New( path ).Parse( text )
    vol.templates[path] = &Template{
        Template: tmpl,
        text: text,
    }
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
 + walkFuncBasedOn
 *
 * filepath.WalkFunc
 * relpath must be volume-path-ready e.g.) templates/
 */

func( vol *Volume ) ImportDirectory( path string ) error {
    return filepath.Walk( path, vol.WalkFuncBasedOn( path ) )
}

func( vol *Volume ) walkFuncBasedOn( basePath string ) filepath.WalkFunc {

    return func( osPath string, fi os.FileInfo, err error ) error {

        // Ignore dir
        if fi.IsDir() {
            return nil
        }

        // Rel
        relPath, relErr := filepath.Rel( basePath, osPath )
        if relErr != nil {
            return ErrInvalidPath.Append( relErr, "basePath:", basePath, "osPath:", osPath )
        }
        relPath = filepath.ToSlash( relPath )

        // Open
        f, err := os.Open( osPath )
        if err != nil {
            return ErrFileError.Append( osPath )
        }

        // Add and ignore error
        err = vol.PutItem( relPath, f, fi.ModTime() )
        if !util.ErrorHasPrefix( err, ErrInvalidPath ) {
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

    for _, fh := range zr.File {

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

    return nil

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
        }
        return nil
    }

    // Assets
    astexp := &assetExporter{ vol.assets }
    err := exportFunc( []Exporter{ astexp } )
    if err != nil {
        return err
    }

    // Locale
    i1exp := &i18nExporter{ vol.i18n }
    err = exportFunc( []Exporter{ i1exp } )
    if err != nil {
        return err
    }

    // Templats
    tmplexp := &templateExporter{ vol.templates }
    err = exportFunc( []Exporter{ tmplexp } )
    if err != nil {
        return err
    }

    zwr.Close()
    return nil

}