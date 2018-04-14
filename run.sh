#!/bin/bash

set -e

./build.sh

dest="$1"
shift 1

ruby src/compiler/main.rb "$dest" *.slg "$@"

bin/slang "$dest"
