package types

import (
	"math/big"
)

func String__get(_ *Program, arguments ...Value) (Value, error) {
	index := arguments[1].(*big.Int).Int64()
	str := arguments[0].(string)
	return string(str[index]), nil
}

func String___eq_(_ *Program, arguments ...Value) (Value, error) {
	str1 := arguments[0]
	str2 := arguments[1]
	return str1 == str2, nil
}
