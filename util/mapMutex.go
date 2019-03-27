package util

import (
    "github.com/hjjg200/together"
)

type MapMutex struct {
    hs *together.HoldSwitch
}

const (
    c_mapMutexRead = 0
    c_mapMutexWrite = 1
)

func NewMapMutex() *MapMutex {
    return &MapMutex{
        hs: together.NewHoldSwitch(),
    }
}

func( mmx *MapMutex ) BeginRead() {
    mmx.hs.Add( c_mapMutexRead, 1 )
}

func( mmx *MapMutex ) EndRead() {
    mmx.hs.Done( c_mapMutexRead )
}

func( mmx *MapMutex ) BeginWrite() {
    mmx.hs.Add( c_mapMutexWrite, 1 )
}

func( mmx *MapMutex ) EndWrite() {
    mmx.hs.Done( c_mapMutexWrite )
}