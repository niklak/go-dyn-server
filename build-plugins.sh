#! /bin/bash

for folder in ./plugins/handlers/*; do
    if [ -d "$folder" ]; then
        base_name=$(basename "$folder")
        echo "Building $folder"
        go build -buildmode=plugin -o include/handlers/${base_name}-${GOOS}-${GOARCH}.so -ldflags "-w -s" ./${folder}
    fi
done


for folder in ./plugins/middlewares/*; do
    if [ -d "$folder" ]; then
        base_name=$(basename "$folder")
        echo "Building $folder"
        go build -buildmode=plugin -o include/middlewares/${base_name}-${GOOS}-${GOARCH}.so -ldflags "-w -s" ./${folder}
    fi
done