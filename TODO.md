# TODO

1. Use filepath.Walk for archiving and processing assets
1. CachedRoute RealtimeRoute, template
1. Option to view directory in WithSymlink
1. Cache localized output of cached assets
1. Realtime locale in realtime methods
1. Put default locale in localizer not sprout instance
1. Put type Session into type Request
1. Crate which includes asset, locale,
1. Domain
    - put server aliases into mux
    - a server can have several muxes
    - direct requests to matching aliases (regexp)

## Structure

|- Sprout
    |- Servers
        |- Spaces
            |- Volume
                |- Localizer
                |- Template
                |- Asset
                |- Cache