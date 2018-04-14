package types

func Vector__conj(arguments ...Value) (Value, error) {
	vec := arguments[0].(Vector)
	for i := 1; i < len(arguments); i++ {
		vec = vec.PushBack(arguments[i])
	}
	return vec, nil
}

func Vector__head(arguments ...Value) (Value, error) {
	vec := arguments[0].(Vector)
	if vec.Len() != 0 {
		return vec.Get(0), nil
	} else {
		return Nil, nil
	}
}

func Vector__tail(arguments ...Value) (Value, error) {
	vec := arguments[0].(Vector)
	if vec.Len() == 0 {
		return NewVectorImpl(), nil
	} else {
		return vec.GetSlice(1, vec.Len()), nil
	}
}
