package volume

import (
    "archive/zip"
    "encoding/json"
    "io"
    "time"

    "../i18n"
)

type Exporter interface{
    Export( *zip.Writer ) error
}


type assetExporter struct{
    assets map[string] *Asset
}

func( astexp *assetExporter ) Export( zwr *zip.Writer ) error {

    // range assets
    for path, ast := range astexp.assets {

        // Create
        wr, err := zwr.CreateHeader( &zip.FileHeader{
            Name: path,
            Modified: ast.modTime,
        } )
        if err != nil {
            return ErrAssetExport.Append( path )
        }

        // Write
        io.Copy( wr, ast )

    }

    return nil

}

type i18nExporter struct{
    i18n *i18n.I18n
}

func( i1exp *i18nExporter ) Export( zwr *zip.Writer ) error {

    modTime := time.Now()

    // range locales
    for lcName, lc := range i1exp.i18n.Locales() {

        // Map
        lcMap := make( map[string] map[string] string )
        lcMap[lcName] = lc.Set()

        // json
        data, err := json.Marshal( lcMap )
        if err != nil {
            return ErrI18nExport.Append( lcName, err )
        }

        // Create
        wr, err := zwr.CreateHeader( &zip.FileHeader{
            Name: c_i18nDirectory + "/" + lcName + ".json",
            Modified: modTime,
        } )
        if err != nil {
            return ErrI18nExport.Append( lcName, err )
        }

        // Write
        wr.Write( data )

    }

    return nil

}

type templateExporter struct{
    templates map[string] *Template
}

func( tmplexp *templateExporter ) Export( zwr *zip.Writer ) error {

    modTime := time.Now()

    // ranage templates
    for path, tmpl := range tmplexp.templates {

        // Create
        wr, err := zwr.CreateHeader( &zip.FileHeader{
            Name: path,
            Modified: modTime,
        } )
        if err != nil {
            return ErrTemplateExport.Append( path )
        }

        // Write
        wr.Write( []byte( tmpl.Text() ) )

    }

    return nil

}