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

func( vol *Volume ) LocalizedAsset( )
func( vol *Volume ) Localize( lcName, src string ) string {
    return vol.i18n.Localize( lcName, src )
}
func( vol *Volume ) ExecuteTemplate( path string, data interface{} ) ( string, error ) {
    
    // 
    tmpl, ok := vol.Template( path )
    if !ok {
        return "", ErrTemplateNonExistent.Append( path )
    }
    
    //
    buf := bytes.NewBuffer()
    err := tmpl.Execute( buf, data )
    if err != nil {
        return "", ErrTemplateExecError.Append( err )
    }
    
    return buf.String(), nil
    
}

// Importers

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

func( vol *Volume ) ImportLocaleDirectory( path string ) error {
    return vol.i18n.ImportDirectory( path )
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