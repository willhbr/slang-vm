package funcs

import (
	"../ds"
	"../vm"
)

type Closure interface {
	IsBuiltin() bool
}

type SlangClosure struct {
	ProgramPosition  uint
	IsProtocolMethod bool
	// TODO Replace this with a more efficient mapping
	Registers []ds.Value
}

func NewSlangClosure(position uint) SlangClosure {
	return SlangClosure{
		ProgramPosition:  position,
		IsProtocolMethod: false,
		Registers:        make([]ds.Value, 100, 100),
	}
}

type GoClosure struct {
	Function func(*vm.Coroutine, ...ds.Value) ds.Value
}

func (g GoClosure) IsBuiltin() bool    { return true }
func (g SlangClosure) IsBuiltin() bool { return false }

// This is where the stdlib lives
// The stdlib should always be in the start of the array, so it can be expanded
//go:generate ruby ../../compiler/builtins.rb ./generated_funcs.go
