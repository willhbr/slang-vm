package main

import (
	"./ds"
	op "./op_codes"
	"fmt"
	"io/ioutil"
	"os"
)

type Frame struct {
	Registers []ds.Value
}

func NewFrame() *Frame {
	return &Frame{Registers: make([]ds.Value, 100, 100)}
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
		fmt.Printf("Running: %s\n", op.ToString(operation))
		index++
		switch operation {
		case op.LOAD:
			vm.Stack.Push(currentFrame.Registers[program[index]])
			index++
		case op.STORE:
			value := vm.Stack.Pop()
			register := program[index]
			currentFrame.Registers[register] = value
			index++
		case op.DISPATCH:
			value := vm.Stack.Pop()
			// This will be the ID of the func to call, but not yet
			index++
			fmt.Println(value)
		case op.APPLY:
			panic("Can't do APPLY yet")
		case op.CONST_I:
			value := program[index]
			index++
			vm.Stack.Push(ds.Value(value))
		case op.CONST_S:
			panic("Can't do CONST_S yet")
		case op.JUMP:
			panic("Can't do JUMP yet")
		case op.AND:
			increase := program[index]
			index++
			if vm.Stack.Peek() != 0 {
				index += int(increase)
				if index >= size {
					break
				}
			}
		case op.OR:
			panic("Can't do OR yet")
		case op.RETURN:
			panic("Can't do RETURN yet")
		case op.NEW_MAP:
			m := ds.NewMap()
			vm.Stack.Push(m)
		case op.NEW_VECTOR:
			v := ds.NewVector()
			vm.Stack.Push(v)
		case op.NEW_LIST:
			l := ds.NewList()
			vm.Stack.Push(l)
		case op.CONS:
			panic("Can't do CONS yet")
		case op.INSERT:
			panic("Can't do INSERT yet")
		}
	}
}

func main() {
	program, err := ioutil.ReadFile(os.Args[1])
	fmt.Println(program)
	if err != nil {
		panic(err)
	}
	coroutine := NewCoroutine()
	prog := Program(program)
	coroutine.Program = &prog
	coroutine.Run(0)
}
