package cache

type Cache struct {
    entries map[string] *Entry
    hs *together.HoldSwitch
}

func NewCache() *Cache {

}

