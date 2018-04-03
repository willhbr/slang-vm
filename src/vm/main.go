package main

import (
	"./ds"
	"./funcs"
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

type Program struct {
	Instructions []byte
	Strings      []string
}

func NewCoroutine() *Coroutine {
	return &Coroutine{Stack: MakeStack(), CurrentFrame: NewFrame()}
}

func (vm *Coroutine) Run(startIndex int) {
	program := vm.Program.Instructions
	strings := vm.Program.Strings
	size := len(program)
	index := startIndex
	currentFrame := vm.CurrentFrame
	for index < size {
		operation := program[index]
		fmt.Println(op.ToString(operation))
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
		case op.CALL_LOCAL:
			value := vm.Stack.Pop()
			// This will be the ID of the func to call, but not yet
			index++
			fmt.Println(value)
		case op.CALL_METHOD:
			value := vm.Stack.Pop()
			module := program[index]
			index++
			method := program[index]
			index++
			fun := funcs.Modules[module][method]
			vm.Stack.Push(fun(value))
		case op.APPLY:
			panic("Can't do APPLY yet")
		case op.CONST_I:
			value := program[index]
			index++
			vm.Stack.Push(ds.Value(int(value)))
		case op.CONST_S:
			idx := program[index]
			index++
			str := strings[idx]
			vm.Stack.Push(ds.Value(str))
		case op.CONST_TRUE:
			vm.Stack.Push(ds.Value(true))
		case op.CONST_FALSE:
			vm.Stack.Push(ds.Value(false))
		case op.CONST_NIL:
			vm.Stack.Push(ds.Nil)
		case op.JUMP:
			increase := int(program[index])
			index++
			index += increase
		case op.AND:
			increase := program[index]
			index++
			if vm.Stack.Pop() != 0 {
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
		default:
			panic(fmt.Errorf("Unknown instruction at %d: %d", index, program[index]))
		}
	}
}

func ParseStrings(instructions []byte) ([]string, int) {
	count := int(instructions[0])
	strings := make([]string, count, count)
	startPosition := 1

	for index := 0; index < count; index++ {
		length := int(instructions[startPosition])
		startPosition++
		endPosition := startPosition + length
		strings[index] = string(instructions[startPosition:endPosition])
		startPosition = endPosition
	}
	return strings, startPosition
}

func main() {
	instructions, err := ioutil.ReadFile(os.Args[1])
	fmt.Println(instructions)
	if err != nil {
		panic(err)
	}
	coroutine := NewCoroutine()
	strings, startIndex := ParseStrings(instructions)
	prog := Program{Instructions: instructions, Strings: strings}
	coroutine.Program = &prog
	coroutine.Run(startIndex)
}
