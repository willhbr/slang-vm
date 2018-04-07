package funcs

import (
	"../ds"
	"../vm"
)

func Kernel__type(co *vm.Coroutine, arguments ...ds.Value) ds.Value {
	if len(arguments) != 1 {
		panic("Too many arguments to Kernel.type")
	}
	return ds.GetType(arguments[0])
}
