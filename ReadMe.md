to test the krakend plugin
```
docker build -t plugin .
docker run -p 8080:8080 plugin
```


```
go mod init krakend-debugger.go
go mod tidy

go build  -buildmode=plugin -o plugins/krakend-debugger.so ./plugins/krakend-debugger.go
krakend run -c krakend.plugin.json
krakend check-plugin  -s plugins/go.sum
```