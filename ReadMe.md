to compile go files use:

```
go build  -buildmode=plugin -o headerModPlugin.so headerModPlugin.go

go build  -buildmode=plugin -o plugins/headerModPlugin.so plugins/headerModPlugin.go
```

to run karkend

```
krakend run --config krakend.plugin.json  --debug

```

```
docker run -p 8080:8080 -v "${PWD}:/etc/krakend/" devopsfaith/krakend run  -c /etc/krakend/krakend.json
```