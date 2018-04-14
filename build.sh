#!/bin/bash

set -e
mkdir -p bin

if ! which peds > /dev/null; then
  go get github.com/tobgu/peds/cmd/peds
fi

ruby generics.rb src/vm

back="$PWD"
cd src/vm/op_codes
go generate
cd "$back"
cd src/vm/types
go generate
cd "$back"
go build -o bin/slang src/vm/main.go
