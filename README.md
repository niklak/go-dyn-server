# go-dyn-server

It is an extended example web server written in Go with support of dynamic libs.

## Building plugins

```bash
# http.Handle
 go build -buildmode=plugin -o include/handlers/ip.so ./plugins/handlers/ip
 # Pre-Middleware
 go build -buildmode=plugin -o include/middlewares/cors.so ./plugins/middlewares/cors
```
or

```bash
 go build -buildmode=plugin -o include/handlers/ip.so -ldflags "-w -s" ./plugins/handlers/ip.go
```

## Building server
```bash
go build -o dyn-server cmd/dyn-server/*
```

## Running server
SERVER_PLUGIN_ROOT must contain `handlers` directory.
```bash
SERVER_PLUGIN_ROOT=./include ./dyn-server
```