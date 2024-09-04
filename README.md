# go-dyn-server

It is an extended example web server written in Go with support of dynamic libs.

## Important
To use go plugins there are some requirements:
- GOOS and GOARCH of the plugin must match the GOOS and GOARCH of the runtime of the program that loads the plugin.
- The plugin must be compiled with the same version of Go as the runtime (GOVERSION).
- I assume that gcc must be also the same version for plugins and application. If you using alpine image, plugins must be compiled also with alpine image.
- Plugins built with go1.22.1 does not work within go1.23.0 runtime.

So, not only you need to keep plugins matching GOOS, GOARCH, but GOVERSION (or toolchain) but also gcc version. Of course it can be automated, but it can be messy.

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

### Benchmark parameters:
- Threads: 8
- Connections: 64
- Duration: 1 minute

### Testing machine:

- OS: Linux Mint 21.2 (Victoria)
- Kernel: Linux version 6.8.0-40-generic
- ROG Strix G713RW (AMD Ryzen 9 6900HX, Kingston Fury Impact 32GB DDR5 4800MHz SO-DIMM)

```bash
wrk -t8 -c64 -d1m http://127.0.0.1:8080/ip
```

 | Type | Lat. Avg | Lat. Stdev | Lat. Max | Lat. +/- Stdev | Req. Avg | Req. Stdev | Req. Max | Req. +/- Stdev |
 |:-----|----------|------------|----------|----------------|----------|------------|----------|----------------|
 |Static|272.28us|278.22us|4.73ms|83.35%|35.90k|1.99k|58.97k|72.30%|
 |Dynamic|303.64us|334.08us|6.80ms|84.74%|34.22k|3.31k|43.43k|69.35%|

Right now it is not pretty convenient.