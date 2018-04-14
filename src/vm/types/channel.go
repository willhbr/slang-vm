package types

func Channel__new(arguments ...Value) (Value, error) {
	size := 0
	if len(arguments) == 1 {
		size = arguments[0].(int)
	}
	return make(chan Value, size), nil
}

func Channel__receive(arguments ...Value) (Value, error) {
	ch := arguments[0].(chan Value)
	return <-ch, nil
}

func Channel__send(arguments ...Value) (Value, error) {
	ch := arguments[0].(chan Value)
	value := arguments[1]
	ch <- value
	return Nil, nil
}
