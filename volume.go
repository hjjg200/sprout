package sprout

/*
 + VOLUME
 *
 * A volume contains assets, locale files, template files, etc.
 */

type Volume struct {
    assets     map[string] *Asset
    localizer  *Localizer
    templates  *template.Template
    hold_group *together.HoldGroup
}

func( _volume *Volume ) Asset( _key string ) *Asset {}
func( _volume *Volume ) PutAsset( _key string, _asset *Asset ) error {}
func( _volume *Volume ) Template( _key string ) *template.Template {}
func( _volume *Volume ) PutTemplate( _key string, _content string ) error {}
func( _volume *Volume ) Locale( _key string ) *Locale {}
func( _volume *Volume ) Locales() []string {}
func( _volume *Volume ) PutLocale( _key string, _locale *Locale ) {}
func( _volume *Volume ) Localize( _source, _locale string ) ( string, error ) {}
func( _volume *Volume ) Update() error {} // Update the volume from the source
func( _volume *Volume ) WhitelistExtension( _ext string ) error {}
func( _volume *Volume ) WhitelistedExtensions() []string {}
func( _volume *Volume ) ToArchive( _path string ) error {}