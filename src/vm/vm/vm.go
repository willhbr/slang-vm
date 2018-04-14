package vm

import (
	"../types"
	"fmt"
)

type Frame struct {
	Registers     []types.Value
	CallingFrame  *Frame
	ContinueIndex int
	StackPosition int
	CatchIndexes  IntStack
}

func NewFrame() *Frame {
	return &Frame{
		Registers:     make([]types.Value, 100, 100),
		CallingFrame:  nil,
		ContinueIndex: 0,
		CatchIndexes:  IntStack{values: new([]int)},
	}
}

func NewFrameFrom(calling *Frame) *Frame {
	return &Frame{
		Registers:     make([]types.Value, 100, 100),
		CallingFrame:  calling,
		ContinueIndex: 0,
		CatchIndexes:  IntStack{values: new([]int)},
	}
}

// lol no generics
type IntStack struct {
	values *[]int
}

func (s IntStack) String() string {
	return fmt.Sprintf("%+v", *s.values)
}

func (s *IntStack) Pop() int {
	stack := *s.values
	value, values := stack[len(stack)-1], stack[:len(stack)-1]
	*s.values = values
	return value
}

func (s *IntStack) Trim(newSize int) {
	if newSize > len(*s.values) {
		return
	}
	*s.values = (*s.values)[0:newSize]
}

func (s IntStack) IsEmpty() bool {
	return len(*s.values) == 0
}

func (s IntStack) Push(val int) {
	*s.values = append(*s.values, val)
}

func (s IntStack) Peek() int {
	return (*s.values)[len(*s.values)-1]
}

func (s IntStack) Len() int {
	return len(*s.values)
}

type Stack struct {
	values *[]types.Value
}

func (s Stack) String() string {
	return fmt.Sprintf("%+v", *s.values)
}

func (s *Stack) Pop() types.Value {
	stack := *s.values
	value, values := stack[len(stack)-1], stack[:len(stack)-1]
	*s.values = values
	return value
}

func (s *Stack) Trim(newSize int) {
	if newSize > len(*s.values) {
		return
	}
	*s.values = (*s.values)[0:newSize]
}

func (s Stack) Push(val types.Value) {
	*s.values = append(*s.values, val)
}

func (s Stack) Peek() types.Value {
	return (*s.values)[len(*s.values)-1]
}

func (s Stack) Len() int {
	return len(*s.values)
}

// Get something a certain offset down, or false if out of range
func (s Stack) PeekFromTopMinus(distance int) (types.Value, bool) {
	idx := len(*s.values) - distance
	if idx < 0 {
		return nil, false
	}
	return (*s.values)[idx], true
}

func MakeStack() Stack {
	vals := make([]types.Value, 0, 100)
	return Stack{values: &vals}
}

type Coroutine struct {
	Stack        Stack
	CurrentFrame *Frame
	Program      *Program
}

type Program struct {
	Instructions []byte
	Strings      []string
}

func NewCoroutine() *Coroutine {
	return &Coroutine{Stack: MakeStack(), CurrentFrame: NewFrame()}
}
