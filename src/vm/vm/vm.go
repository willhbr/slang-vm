package vm

import (
	"../ds"
)

type Frame struct {
	Registers     []ds.Value
	CallingFrame  *Frame
	ContinueIndex int
}

func NewFrame() *Frame {
	return &Frame{Registers: make([]ds.Value, 100, 100), CallingFrame: nil, ContinueIndex: 0}
}

func NewFrameFrom(calling *Frame) *Frame {
	return &Frame{Registers: make([]ds.Value, 100, 100), CallingFrame: calling, ContinueIndex: 0}
}

type Stack struct {
	values *[]ds.Value
}

func (s *Stack) Pop() ds.Value {
	stack := *s.values
	value, values := stack[len(stack)-1], stack[:len(stack)-1]
	*s.values = values
	return value
}

func (s Stack) Push(val ds.Value) {
	*s.values = append(*s.values, val)
}

func (s Stack) Peek() ds.Value {
	return (*s.values)[len(*s.values)-1]
}

func MakeStack() Stack {
	vals := make([]ds.Value, 0, 100)
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
