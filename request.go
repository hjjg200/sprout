package sprout

/*
 + REQUEST
 *
 * A request is the context of a http request, containing locale, writer, request, etc.
 * Requests have functions that are for responding to requests
 */

type Request struct {
    writer http.ResponseWriter
    body   *http.Request
    locale string
}