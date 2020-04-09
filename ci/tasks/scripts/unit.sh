#!/bin/bash

set -eux

export GOPATH="$PWD/gopath"

cd ansi
go test ./...