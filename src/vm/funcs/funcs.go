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

func Channel__new(co *vm.Coroutine, arguments ...ds.Value) ds.Value {
	size := 0
	if len(arguments) == 1 {
		size = arguments[0].(int)
	}
	return make(chan ds.Value, size)
}

func Channel__receive(co *vm.Coroutine, arguments ...ds.Value) ds.Value {
	ch := arguments[0].(chan ds.Value)
	return <-ch
}

func Channel__send(co *vm.Coroutine, arguments ...ds.Value) ds.Value {
	ch := arguments[0].(chan ds.Value)
	value := arguments[1]
	ch <- value
	return ds.Nil
}
