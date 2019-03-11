package sprout

/*
 + VOLUME MANAGER
 *
 * A volume manager is a pseudo-static object used for managing Volume
 */

type volume_factory struct {}
var  static_volume_factory = &volume_factory{}

func VolumeFactory() *volume_factory {}
func( _volfac *volume_factory ) FromDirectory( _path string ) ( *Volume, error ) {

    _volume := &Volume{
        assets: make( map[string] *Asset ),
        localizer: LocaleFactory().New(),
        templates: template.New( "" ),
        hold_group: together.NewHoldGroup(),
    }

    _walk_func := func( _2path string, _2info os.FileInfo, _2error error ) error {
        switch {
        case strings.HasPrefix( _2path, VolumeVariables().ConstAssetDirectory() + "/" ):

        case strings.HasPrefix( _2path, VolumeVariables().ConstLocaleDirectory() + "/" ):
            // Check extension
            _2extension := filepath.Ext( _2path )
            if _2extension == VolumeVariables().ConstLocaleExtension() {
                // Read file
                _2file, _2locale_error := os.Open( _2path )
                if _2locale_error != nil {
                    return _2locale_error
                }
                _2buffer := bytes.NewBuffer()
                io.Copy( _2buffer, _2file )
                _2file.Close()
                // Add locale
                _2locale, _2locale_error := LocaleFactory().FromJSON( _2buffer.Bytes() )
                if _2locale_error != nil {
                    return _2locale_error
                }
                _volume.localizer.AddLocale( _2locale )
            }
        case strings.HasPrefix( _2path, VolumeVariables().ConstTemplateDirectory() + "/" ):
            // Chcek extension
            _2found := false
            _2extension = filepath.Ext( _2path )
            for _, _2value := range VolumeVariables().TemplateExtensions() {
                if _2value == _2extension {
                    _2found = true
                }
            }
            if _2found == false {
                return nil
            }
            // Read file
            _2file, _2template_error := os.Open( _2path )
            if _2template_error != nil {
                return _2template_error
            }
            _2buffer := bytes.NewBuffer()
            io.Copy( _2buffer, _2file )
            _2file.Close()
            // Add to templates
            _, _2template_error = _volume.templates.New( _2path ).Parse( _2buffer.String() )
            if _2template_error != nil {
                return _2template_error
            }
        }
        return nil
    }

}
func( _volfac *volume_factory ) New() ( *Volume ) {}