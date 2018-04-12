package types

func Vector__conj(arguments ...Value) Value {
	vec := arguments[0].(*Vector)
	for i := 1; i < len(arguments); i++ {
		vec = vec.Append(arguments[i])
	}
	return vec
}

func Vector__head(arguments ...Value) Value {
	vec := arguments[0].(*Vector)
	return vec.Get(0)
}

func Vector__tail(arguments ...Value) Value {
	vec := arguments[0].(*Vector)
	return vec.Slice(1, vec.Len())
}
