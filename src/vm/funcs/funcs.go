package funcs

import (
	"../ds"
	"../vm"
)

type SlangClosure struct {
	ProgramPosition uint
	// TODO Replace this with a more efficient mapping
	Registers []ds.Value
}

func NewSlangClosure(position uint) SlangClosure {
	return SlangClosure{
		ProgramPosition: position,
		Registers:       make([]ds.Value, 100, 100),
	}
}

type GoClosure struct {
	Function func(*vm.Coroutine, ...ds.Value) ds.Value
}

type ProtocolClosure struct {
	ID int
}

func (g GoClosure) IsBuiltin() bool       { return true }
func (g SlangClosure) IsBuiltin() bool    { return false }
func (g ProtocolClosure) IsBuiltin() bool { return false }

// This is where the stdlib lives
// The stdlib should always be in the start of the array, so it can be expanded
//go:generate ruby ../../compiler/builtins.rb ./generated_funcs.go
