# go-dyn-server

It is an extended example web server written in Go with support of dynamic libs.

## Building plugins

```bash

# http.Handle
GOOS=linux GOARCH=amd64 go build -buildmode=plugin -o include/handlers/ip-${GOOS}-${GOARCH}.so ./plugins/handlers/ip
 # Middleware
GOOS=linux GOARCH=amd64 go build -buildmode=plugin -o include/middlewares/cors-${GOOS}-${GOARCH}.so ./plugins/middlewares/cors
```
or

```bash
 go build -buildmode=plugin -o include/handlers/ip-${GOOS}-${GOARCH}.so -ldflags "-w -s" ./plugins/handlers/ip.go
```

### build all plugins for a specific combination of GOOS-GOARCH:

```bash
GOOS=linux GOARCH=amd64 bash ./build-plugins.sh
GOOS=linux GOARCH=arm64 bash ./build-plugins.sh
```
or

```bash
GOOS=linux GOARCH=amd64 ./build-plugins.sh
GOOS=linux GOARCH=arm64 ./build-plugins.sh
```

## Building server
```bash
cd cmd/dyn-server
go build -o dyn-server .
# static server for comparison with dyn-server
go build -tags static -o static-server .
```

## Running server
SERVER_PLUGIN_ROOT must contain `handlers` directory.
```bash
SERVER_PLUGIN_ROOT=./include ./cmd/dyn-server/dyn-server
```


## Benchmark

Benchmark parameters:
- Threads: 8
- Connections: 64
- Duration: 1 minute


WIN 11 WSL2 Ubuntu 22.04

(ASUS ROG Strix G17)

AMD Ryzen 9 6900HX

Kingston Fury Impact 32GB DDR5 4800MHz SO-DIMM

```bash
wrk -t8 -c64 -d1m http://127.0.0.1:8080/ip
```

Static

 | Type | Lat. Avg | Lat. Stdev | Lat. Max | Lat. +/- Stdev | Req. Avg | Req. Stdev | Req. Max | Req. +/- Stdev |
 |:-----|----------|------------|----------|----------------|----------|------------|----------|----------------|
 |Static|435.24us|601.95us|18.95ms|87.53%|29.06k|2.57k|38.19k|68.44%|
 |Dynamic|441.52us|580.51us|12.41ms|86.93%|28.19k|2.44k|38.81k|67.71%|

