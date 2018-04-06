package funcs

import (
	"../ds"
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
	Function func(...ds.Value) ds.Value
}

func (g GoClosure) IsBuiltin() bool    { return true }
func (g SlangClosure) IsBuiltin() bool { return false }

// This is where the stdlib lives
// The stdlib should always be in the start of the array, so it can be expanded
//go:generate ruby ../../compiler/builtins.rb ./generated_funcs.go

func IO__puts(arguments ...ds.Value) ds.Value {
	for i := range arguments {
		fmt.Print(arguments[i])
	}
	fmt.Println()
	return ds.Nil
}

func Kernel__type(arguments ...ds.Value) ds.Value {
	if len(arguments) != 1 {
		panic("Too many arguments to Kernel.type")
	}
	return ds.GetType(arguments[0])
}
