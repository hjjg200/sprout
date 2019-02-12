package sprout

/*
 + LOCALE
 *
 * A locale is a set of strings that is capable of localizing text on its own
 */

type Locale struct {
    lang   string
    locale map[string] string
}

func( _locale *Locale ) Language() string {}
func( _locale *Locale ) Locale() map[string] string {}
func( _locale *Locale ) Localize( _source string ) string {}
func( _locale *Locale ) localize( _source string, _threshold int ) string {}
func( _locale *Locale ) SetThreshold( _threshold int ) {}