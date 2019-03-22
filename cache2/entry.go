package cache

type Entry struct {
    pData    *[]byte
    pNewData *[]byte
    hs       *together.HoldSwitch
}

func NewEntry( data []byte ) *Entry {

    newData := make( []byte, len( data ) )
    copy( newData, data )

    // hs
    hs := together.NewHoldSwitch()
    hs.Handlers( c_entryRead, e.beginRead, nil )

    return &Entry{
        pData: &newData,
        pNewData: nil,
        hs: hs,
    }

}

func( e *Entry ) Data() []byte {
    e.hs.Add( c_entryRead, 1 )
    defer e.hs.Done( c_entryRead )
    return *e.data
}

func( e *Entry ) SetData( data []byte ) {

    e.hs.Add( c_entryWrite, 1 )

    newData := make( []byte, len( data ) )
    copy( newData, data )
    e.pNewData = &newData

    e.hs.Done( c_entryWrite )
    e.hs.Flush()

}

func( e *Entry ) beginRead() {
    if e.pNewData != nil {
        e.pData = e.pNewData
        e.pNewData = nil
    }
}