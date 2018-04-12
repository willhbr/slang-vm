package types

func Channel__new(arguments ...Value) Value {
	size := 0
	if len(arguments) == 1 {
		size = arguments[0].(int)
	}
	return make(chan Value, size)
}

func Channel__receive(arguments ...Value) Value {
	ch := arguments[0].(chan Value)
	return <-ch
}

func Channel__send(arguments ...Value) Value {
	ch := arguments[0].(chan Value)
	value := arguments[1]
	ch <- value
	return Nil
}
