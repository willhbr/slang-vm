package funcs

import (
	"../ds"
	"fmt"
)

//go:generate ruby ../../compiler/builtins.rb ./generated_funcs.go

func IO__puts(argument ds.Value) ds.Value {
	fmt.Println(argument)
	return ds.Nil
}
