package funcs

import (
	"../ds"
	"../vm"
	"math/big"
)

func Kernel__type(co *vm.Coroutine, arguments ...ds.Value) ds.Value {
	if len(arguments) != 1 {
		panic("Too many arguments to Kernel.type")
	}
	return ds.GetType(arguments[0])
}

func Kernel__minus(co *vm.Coroutine, arguments ...ds.Value) ds.Value {
	result := big.Int{}
	result.Sub(arguments[0].(*big.Int), arguments[1].(*big.Int))
	return &result
}

func Kernel__lessThan(co *vm.Coroutine, arguments ...ds.Value) ds.Value {
	return arguments[0].(*big.Int).Cmp(arguments[1].(*big.Int)) == -1
}

func Kernel__times(co *vm.Coroutine, arguments ...ds.Value) ds.Value {
	result := big.Int{}
	result.Mul(arguments[0].(*big.Int), arguments[1].(*big.Int))
	return &result
}
