package network

type Request struct {
    body   *http.Request
    writer http.ResponseWriter
    locale string
}