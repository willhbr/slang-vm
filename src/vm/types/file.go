package types

import (
	"io/ioutil"
)

func File__read(_ *Program, arguments ...Value) (Value, error) {
	name := arguments[0].(string)
	result, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, err
	}
	return string(result), nil
}
