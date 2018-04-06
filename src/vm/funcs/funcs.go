package funcs

import (
	"../ds"
	"../vm"
	"fmt"
)

type Closure interface {
	IsBuiltin() bool
}

type SlangClosure struct {
	ProgramPosition  uint
	IsProtocolMethod bool
}

type GoClosure struct {
	Function func(*vm.Coroutine, ...ds.Value) ds.Value
}

func (g GoClosure) IsBuiltin() bool    { return true }
func (g SlangClosure) IsBuiltin() bool { return false }

// This is where the stdlib lives
// The stdlib should always be in the start of the array, so it can be expanded
//go:generate ruby ../../compiler/builtins.rb ./generated_funcs.go

func IO__puts(co *vm.Coroutine, arguments ...ds.Value) ds.Value {
	for i := range arguments {
		arg := arguments[i]
		// TODO Call a real method to turn things into a string
		switch arg.(type) {
		case ds.Atom:
			fmt.Print(co.Program.Strings[int(arg.(ds.Atom))])
		default:
			fmt.Print(arguments[i])
		}
	}
	fmt.Println()
	return ds.Nil
}

func Kernel__type(co *vm.Coroutine, arguments ...ds.Value) ds.Value {
	if len(arguments) != 1 {
		panic("Too many arguments to Kernel.type")
	}
	return ds.GetType(arguments[0])
}
