#!/bin/bash

set -e

./build.sh

ruby src/compiler/main.rb "$1" "$1"c

if [ "$2" = nc ]; then
  exit
fi
bin/slgc "$1"c
