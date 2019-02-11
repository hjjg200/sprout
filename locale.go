package sprout

/*
 + LOCALE
 *
 * A locale is a set of strings that is capable of localizing text on its own
 */

type Locale struct {
    locale map[string] string
}