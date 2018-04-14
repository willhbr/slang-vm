package vm

import (
	"../types"
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
		CatchIndexes:  MakeCatchIndexStack(),
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

type Coroutine struct {
	Stack        ValueStack
	CurrentFrame *Frame
	Program      *Program
}

type Program struct {
	Instructions []byte
	Strings      []string
}

func NewCoroutine() *Coroutine {
	return &Coroutine{Stack: MakeValueStack(), CurrentFrame: NewFrame()}
}
