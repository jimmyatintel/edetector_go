#!/bin/bash

cd ..
go build -o build/server_$1.exe -ldflags "-X main.version=$1" ./cmd/server/working_server.go
go build -o build/connector_$2.exe -ldflags "-X main.version=$2" ./cmd/connector/connector.go
go build -o build/parser_$3.exe -ldflags "-X main.version=$3" ./cmd/parser/parser.go
go build -o build/treebuilder_$4.exe -ldflags "-X main.version=$4" ./cmd/builder/builder.go