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
	Function func(ds.Value) ds.Value
}

func (g GoClosure) IsBuiltin() bool    { return true }
func (g SlangClosure) IsBuiltin() bool { return false }

// This is where the stdlib lives
// The stdlib should always be in the start of the array, so it can be expanded
//go:generate ruby ../../compiler/builtins.rb ./generated_funcs.go

func IO__puts(argument ds.Value) ds.Value {
	fmt.Println(argument)
	return ds.Nil
}
