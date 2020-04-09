#!/bin/bash
set -eux

### Taken from https://github.com/fuzzitdev/example-go/blob/master/fuzzit.sh

export GOPATH="$PWD/gopath"
export PATH="$GOPATH/bin:$PATH"

cd ansi

## Install go-fuzz
go get -u github.com/dvyukov/go-fuzz/go-fuzz github.com/dvyukov/go-fuzz/go-fuzz-build

## Install clang
apt-get update
apt-get install -y clang

go-fuzz-build -libfuzzer -o ansi.a .
clang -fsanitize=fuzzer ansi.a -o ansi

## Install fuzzit latest version:
wget -O fuzzit https://github.com/fuzzitdev/fuzzit/releases/latest/download/fuzzit_Linux_x86_64
chmod a+x fuzzit

## upload fuzz target for long fuzz testing on fuzzit.dev server or run locally for regression
./fuzzit create job --type "${ACTION}" aoldershaw/ansi ansi