package volume

type Exporter interface{
    Export( *zip.Writer ) error
}
