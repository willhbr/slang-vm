#!/bin/bash

set -e

./build.sh

dest="$1"

ruby src/compiler/main.rb "$dest" *.slg

if [ "$2" = nc ]; then
  exit
fi
bin/slgc "$dest"
