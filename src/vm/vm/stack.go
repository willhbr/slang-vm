package vm

import (
  "../types"
  "fmt"
)

// There is no God

type ValueStack struct {
	values *[]types.Value
}

func (s ValueStack) String() string {
	return fmt.Sprintf("%+v", *s.values)
}

func (s *ValueStack) Pop() types.Value {
	stack := *s.values
	value, values := stack[len(stack)-1], stack[:len(stack)-1]
	*s.values = values
	return value
}

func (s *ValueStack) Trim(newSize int) {
	if newSize > len(*s.values) {
		return
	}
	*s.values = (*s.values)[0:newSize]
}

func (s ValueStack) IsEmpty() bool {
	return len(*s.values) == 0
}

func (s ValueStack) Push(val types.Value) {
	*s.values = append(*s.values, val)
}

func (s ValueStack) Peek() types.Value {
	return (*s.values)[len(*s.values)-1]
}

func (s ValueStack) Len() int {
	return len(*s.values)
}

// Get something a certain offset down, or false if out of range
func (s ValueStack) PeekFromTopMinus(distance int) (types.Value, bool) {
	idx := len(*s.values) - distance
	if idx < 0 {
		return nil, false
	}
	return (*s.values)[idx], true
}

func MakeValueStack() ValueStack {
	vals := make([]types.Value, 0, 100)
	return ValueStack{values: &vals}
}

type CatchIndexStack struct {
	values *[]CatchIndex
}

func (s CatchIndexStack) String() string {
	return fmt.Sprintf("%+v", *s.values)
}

func (s *CatchIndexStack) Pop() CatchIndex {
	stack := *s.values
	value, values := stack[len(stack)-1], stack[:len(stack)-1]
	*s.values = values
	return value
}

func (s *CatchIndexStack) Trim(newSize int) {
	if newSize > len(*s.values) {
		return
	}
	*s.values = (*s.values)[0:newSize]
}

func (s CatchIndexStack) IsEmpty() bool {
	return len(*s.values) == 0
}

func (s CatchIndexStack) Push(val CatchIndex) {
	*s.values = append(*s.values, val)
}

func (s CatchIndexStack) Peek() CatchIndex {
	return (*s.values)[len(*s.values)-1]
}

func (s CatchIndexStack) Len() int {
	return len(*s.values)
}

// Get something a certain offset down, or false if out of range
func (s CatchIndexStack) PeekFromTopMinus(distance int) (*CatchIndex, bool) {
	idx := len(*s.values) - distance
	if idx < 0 {
		return nil, false
	}
	return &(*s.values)[idx], true
}

func MakeCatchIndexStack() CatchIndexStack {
	vals := make([]CatchIndex, 0, 100)
	return CatchIndexStack{values: &vals}
}
