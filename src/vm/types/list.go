package types

func List__new(arguments ...Value) Value {
	list := NewList()
	for i := len(arguments) - 1; i >= 0; i-- {
		list = list.Cons(arguments[i])
	}
	return list
}

func List__conj(arguments ...Value) Value {
	list := arguments[0].(*List)
	for i := len(arguments) - 1; i > 0; i-- {
		list = list.Cons(arguments[i])
	}
	return list
}

func List__head(arguments ...Value) Value {
	list := arguments[0].(*List)
	return list.Head()
}

func List__tail(arguments ...Value) Value {
	list := arguments[0].(*List)
	return list.Tail()
}
