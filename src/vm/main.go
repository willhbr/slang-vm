package main

import (
	"fmt"
	"io/ioutil"
)

const (
	LOAD       = 1
	STORE      = 2
	DISPATCH   = 3
	APPLY      = 4
	CONST_I    = 5
	CONST_S    = 6
	JUMP       = 7
	AND        = 8
	OR         = 9
	RETURN     = 10
	NEW_MAP    = 11
	NEW_VECTOR = 12
	NEW_LIST   = 13
	CONS       = 14
	INSERT     = 15
)

type Value int

type Frame struct {
	Registers []Value
}

func NewFrame() *Frame {
	return &Frame{Registers: make([]Value, 100, 100)}
}

type Stack struct {
	values *[]Value
}

func (s *Stack) Pop() Value {
	stack := *s.values
	value, values := stack[len(stack)-1], stack[:len(stack)-1]
	*s.values = values
	return value
}

func (s Stack) Push(val Value) {
	*s.values = append(*s.values, val)
}

func (s Stack) Peek() Value {
	return (*s.values)[len(*s.values)-1]
}

func MakeStack() Stack {
	vals := make([]Value, 0, 100)
	return Stack{values: &vals}
}

type Coroutine struct {
	Stack        Stack
	CurrentFrame *Frame
	Program      *Program
}

type Program []byte

func NewCoroutine() *Coroutine {
	return &Coroutine{Stack: MakeStack(), CurrentFrame: NewFrame()}
}

func (vm *Coroutine) Run(startIndex int) {
	program := *vm.Program
	size := len(program)
	index := 0
	currentFrame := vm.CurrentFrame
	for index < size {
		operation := program[index]
		index++
		switch operation {
		case LOAD:
			vm.Stack.Push(currentFrame.Registers[program[index]])
			index++
		case STORE:
			value := vm.Stack.Pop()
			register := program[index]
			currentFrame.Registers[register] = value
			index++
		case DISPATCH:
			value := vm.Stack.Pop()
			// This will be the ID of the func to call, but not yet
			index++
			fmt.Println(value)
		case APPLY:
			panic("Can't do APPLY yet")
		case CONST_I:
			value := program[index]
			index++
			vm.Stack.Push(Value(value))
		case CONST_S:
			panic("Can't do CONST_S yet")
		case JUMP:
			panic("Can't do JUMP yet")
		case AND:
			increase := program[index]
			index++
			if vm.Stack.Peek() != 0 {
				index += int(increase)
				if index >= size {
					break
				}
			}
		case OR:
			panic("Can't do OR yet")
		case RETURN:
			panic("Can't do RETURN yet")
		case NEW_MAP:
			panic("Can't do NEW_MAP yet")
		case NEW_VECTOR:
			panic("Can't do NEW_VECTOR yet")
		case NEW_LIST:
			panic("Can't do NEW_LIST yet")
		case CONS:
			panic("Can't do CONS yet")
		case INSERT:
			panic("Can't do INSERT yet")
		}
	}
}

func main() {
	program, err := ioutil.ReadFile("output.bin")
	if err != nil {
		panic(err)
	}
	coroutine := NewCoroutine()
	prog := Program(program)
	coroutine.Program = &prog
	coroutine.Run(0)
}
