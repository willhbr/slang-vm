package ds

import (
	"bytes"
	"fmt"
)

func (v Vector) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("[")
	end := v.Len() - 1
	current := 0
	v.Range(func(value Value) bool {
		buffer.WriteString(fmt.Sprintf("%v", value))
		if current != end {
			buffer.WriteString(" ")
			current++
		}
		return true
	})
	buffer.WriteString("]")
	return buffer.String()
}

func (v Map) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("{")
	end := v.Len() - 1
	current := 0
	v.Range(func(key, value Value) bool {
		buffer.WriteString(fmt.Sprintf("%v %v", key, value))
		if current != end {
			buffer.WriteString(", ")
			current++
		}
		return true
	})
	buffer.WriteString("}")
	return buffer.String()
}
