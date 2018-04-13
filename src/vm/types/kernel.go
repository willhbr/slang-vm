package types

import (
	"math/big"
)

func Kernel__type(arguments ...Value) Value {
	if len(arguments) != 1 {
		panic("Too many arguments to Kernel.type")
	}
	return GetType(arguments[0])
}

func Kernel__minus(arguments ...Value) Value {
	result := big.Int{}
	result.Sub(arguments[0].(*big.Int), arguments[1].(*big.Int))
	return &result
}

func Kernel__lessThan(arguments ...Value) Value {
	return arguments[0].(*big.Int).Cmp(arguments[1].(*big.Int)) == -1
}

func Kernel__times(arguments ...Value) Value {
	result := big.Int{}
	result.Mul(arguments[0].(*big.Int), arguments[1].(*big.Int))
	return &result
}
