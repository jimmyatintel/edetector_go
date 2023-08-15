#!/bin/bash

cd ..
go build -o build/server.exe -ldflags "-X main.version=$1" ./cmd/server/working_server.go
go build -o build/connector.exe -ldflags "-X main.version=$2" ./cmd/connector/connector.go
go build -o build/parser.exe -ldflags "-X main.version=$3" ./cmd/parser/parser.go

# bash -ac 'source ./build/.env && ./build/connector.exe'