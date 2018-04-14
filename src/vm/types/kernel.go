package types

import (
	"math/big"
)

func Kernel__type(_ *Program, arguments ...Value) (Value, error) {
	if len(arguments) != 1 {
		panic("Too many arguments to Kernel.type")
	}
	return GetType(arguments[0]), nil
}

func Kernel__minus(_ *Program, arguments ...Value) (Value, error) {
	result := big.Int{}
	result.Sub(arguments[0].(*big.Int), arguments[1].(*big.Int))
	return &result, nil
}

func Kernel__lessThan(_ *Program, arguments ...Value) (Value, error) {
	return arguments[0].(*big.Int).Cmp(arguments[1].(*big.Int)) == -1, nil
}

func Kernel__times(_ *Program, arguments ...Value) (Value, error) {
	result := big.Int{}
	result.Mul(arguments[0].(*big.Int), arguments[1].(*big.Int))
	return &result, nil
}
