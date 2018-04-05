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

//go:generate ruby ../../compiler/builtins.rb ./generated_funcs.go

func IO__puts(argument ds.Value) ds.Value {
	fmt.Println(argument)
	return ds.Nil
}
