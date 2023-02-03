#!/usr/bin/env bash
export PATH=$PATH:/usr/local/go/bin
GOOS=linux GOARCH=arm64 go build -ldflags "-s -w" -buildvcs=false -o app .