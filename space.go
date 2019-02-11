package sprout

/*
 || SPACE
 *
 * A space handles a single website and it can have a volume and own routing rules
 *
 */

type Space struct {
    name    string // domain
    handler Handler
}