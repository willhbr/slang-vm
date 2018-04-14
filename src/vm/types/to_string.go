package types

import "fmt"

func Int___rArr_string(args ...Value) (Value, error) {
	return fmt.Sprintf("%d", args[0]), nil
}

func String___rArr_string(args ...Value) (Value, error) {
	return fmt.Sprintf("%d", args[0]), nil
}
