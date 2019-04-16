# Sprout

## Var

```text
Sprout => sprt
Space => spc
Server => srv
Request => req
Responder => rsp
Handler => hnd
Volume => vol
Cache => chc
```

## Design Plan

```text
sprout/
- sprout.go
- i18n/
    - i18n.go
...
- volume/
    - volume.go
    - realtime.go # realtime volume
    - zip.go # zip volume
    - directory.go
    ...
    - item/
        - asset.go
        - template.go
    - export/
        - export.go
```


```text
working dir example/
- asset/
- cache/
    - volume/
- i18n/
- template/

```

### To Do

- image optimizer
- the source file must stay intact
    - should not make any file in the source directory
- volumes
    - ZipVolume
        - import single time
    - DirectoryVolume
        - realtime
    - WebVolume
        - realtime
- cache
    - volume