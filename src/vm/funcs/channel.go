package funcs

import (
	"../ds"
	"../vm"
)

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
