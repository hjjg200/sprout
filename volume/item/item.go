package item

type Item interface {
    Bytes() []byte
    String() string
}