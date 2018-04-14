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
	CatchIndexes  CatchIndexStack
}

type CatchIndex struct {
	Index     int
	StackSize int
}

func NewFrame() *Frame {
	return &Frame{
		Registers:     make([]types.Value, 100, 100),
		CallingFrame:  nil,
		ContinueIndex: 0,
		CatchIndexes:  CatchIndexStack{values: new([]CatchIndex)},
	}
}

func NewFrameFrom(calling *Frame) *Frame {
	return &Frame{
		Registers:     make([]types.Value, 100, 100),
		CallingFrame:  calling,
		ContinueIndex: 0,
		CatchIndexes:  CatchIndexStack{values: new([]CatchIndex)},
	}
}

// lol no generics
type CatchIndexStack struct {
	values *[]CatchIndex
}

func NewCatchIndexStack() CatchIndexStack {
	return CatchIndexStack{values: new([]CatchIndex)}
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

func (s Stack) IsEmpty() bool {
	return len(*s.values) == 0
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
