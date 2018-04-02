#!/bin/bash

set -e
mkdir -p bin

if ! which peds > /dev/null; then
  go get github.com/tobgu/peds/cmd/peds
fi

back="$PWD"
cd src/vm/ds
go generate
cd "$back"
go build -o bin/slgc src/vm/main.go
