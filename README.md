# Slang

_A simple compiled Lisp with coroutines and immutable data structures_

## Usage

Install Go (I use 1.8) and Ruby (I use 2.4.1). `build.sh` should install dependencies and create the `slang` binary in `bin/`.

Scripts can be run using `run.sh` (which also rebuilds everything). It will build all `.slg` files in the current directory, writes the compiled program to a file (given as the first argument) and runs that file in the VM.

```
(module Main)

(IO.puts "Hello world!")
```

```shell
./run.sh script.slgc
```
