package cache

type Porter interface {
    Export() ( *Cache, error )
    Import( *Cache ) error
}