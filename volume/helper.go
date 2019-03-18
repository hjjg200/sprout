package volume

func BaseWithoutExt( path string ) string {
    base := filepath.Base( path )
    ext  := filepath.Ext( path )
    if len( ext ) == 0 {
        return base
    }
    return base[:len( base ) - len( ext )]
}

type i18nExporter struct{
    i18n *i18n.I18n
}

func( i1exp *i18nExporter ) Export( zwr *zip.Writer ) error {

    modTime := time.Now()

    // range locales
    for _, lc := range i1exp.i18n.Locales() {

        // Map
        lcMap := make( map[string] map[string] string )
        lcMap[lc.Name()] = lc.Set()

        // json
        data, err := json.Marshal( lcMap )
        if err != nil {
            return ErrI18nExport.Append( lc.Name(), err )
        }

        // Create
        path := c_localeDirectory + "/" + lc.name + ".json"
        fh   := zip.FileHeader{
            Name: path,
            Modified: modTime,
        }
        wr, err := zwr.CreateHeader( fh )
        if err != nil {
            return ErrI18nExport.Append( lc.Name(), err )
        }

        // Write
        wr.Write( data )

    }

    return nil

}