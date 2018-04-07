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

func Kernel__minus(co *vm.Coroutine, arguments ...ds.Value) ds.Value {
	return arguments[0].(int) - arguments[1].(int)
}

func Kernel__lessThan(co *vm.Coroutine, arguments ...ds.Value) ds.Value {
	return arguments[0].(int) < arguments[1].(int)
}

func Kernel__times(co *vm.Coroutine, arguments ...ds.Value) ds.Value {
	return arguments[0].(int) * arguments[1].(int)
}
