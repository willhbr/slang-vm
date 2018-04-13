package types

func Vector__conj(arguments ...Value) Value {
	vec := arguments[0].(Vector)
	for i := 1; i < len(arguments); i++ {
		vec = vec.PushBack(arguments[i])
	}
	return vec
}

func Vector__head(arguments ...Value) Value {
	vec := arguments[0].(Vector)
	if vec.Len() != 0 {
		return vec.Get(0)
	} else {
		return Nil
	}
}

func Vector__tail(arguments ...Value) Value {
	vec := arguments[0].(Vector)
	if vec.Len() == 0 {
		return NewVectorImpl()
	} else {
		return vec.GetSlice(1, vec.Len())
	}
}
