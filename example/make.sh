#!/bin/bash

rm -rf ./out

go build -o ./out/mac/server ./server/main.go
go build -o ./out/mac/client ./client/main.go

export GOOS=linux GOARCH=amd64
go build -o ./out/linux/server ./server/main.go
go build -o ./out/linux/client ./client/main.go

cp -rf ./testdata ./out/ssl
