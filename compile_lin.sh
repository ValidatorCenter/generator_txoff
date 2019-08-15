#!/bin/sh

GOOS=linux GOARCH=amd64 go build -ldflags "-s" -o bin/seed_info t0.go
GOOS=linux GOARCH=amd64 go build -ldflags "-s" -o bin/gen_txoff t1.go