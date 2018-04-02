#!/bin/bash

set -e

./build.sh

ruby src/compiler/main.rb "$1" "$1"c

bin/slgc "$1"c
