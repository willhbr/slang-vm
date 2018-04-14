package types

func List__new(arguments ...Value) (Value, error) {
	l := NewList()
	list := &l
	for i := len(arguments) - 1; i >= 0; i-- {
		list = list.Cons(arguments[i])
	}
	return *list, nil
}

func List__conj(arguments ...Value) (Value, error) {
	l := arguments[0].(List)
	list := &l
	for i := len(arguments) - 1; i > 0; i-- {
		list = list.Cons(arguments[i])
	}
	return *list, nil
}

func List__head(arguments ...Value) (Value, error) {
	list := arguments[0].(*List)
	return list.Head(), nil
}

func List__tail(arguments ...Value) (Value, error) {
	list := arguments[0].(*List)
	return list.Tail(), nil
}
