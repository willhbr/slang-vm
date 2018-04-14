package types

import "fmt"

func Int___rArr_string(_ *Program, args ...Value) (Value, error) {
	return fmt.Sprintf("%d", args[0]), nil
}

func String___rArr_string(_ *Program, args ...Value) (Value, error) {
	return fmt.Sprintf("%d", args[0]), nil
}

func Atom___rArr_string(prog *Program, args ...Value) (Value, error) {
	return ":" + prog.Strings[args[0].(Atom)], nil
}

func Atom__value(prog *Program, args ...Value) (Value, error) {
	return prog.Strings[args[0].(Atom)], nil
}
