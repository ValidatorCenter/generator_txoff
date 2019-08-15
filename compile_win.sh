#!/bin/sh

GOOS=windows GOARCH=amd64 go build -ldflags "-s" -o bin/seed_info.exe t0.go
GOOS=windows GOARCH=amd64 go build -ldflags "-s" -o bin/gen_txoff.exe t1.go