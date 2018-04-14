#!/bin/bash

FLAGS='--show_instructions --show_modules --show_order --show_time'

set -e

./build.sh

failure=''

for suite in test/*; do
  if ruby src/compiler/main.rb "$suite.slgc" "$suite" $FLAGS; then
    bin/slang "$suite.slgc"
  else
    failure='failed'
    echo "$suite failed to compile"
  fi
done

if [ "$failure" = failed ]; then
  exit 1
fi

