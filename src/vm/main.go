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

func (vm *Coroutine) Run(startIndex int) {
	program := vm.Program.Instructions
	strings := vm.Program.Strings
	size := len(program)
	index := startIndex
	currentFrame := vm.CurrentFrame
	for index < size {
		operation := program[index]
		// fmt.Println(op.ToString(operation))
		index++
		switch operation {
		case op.LOAD_LOCAL:
			vm.Stack.Push(currentFrame.Registers[program[index]])
			index++
		case op.LOAD_DEF:
			vm.Stack.Push(funcs.Defs[program[index]])
			index++
		case op.STORE:
			value := vm.Stack.Pop()
			register := program[index]
			currentFrame.Registers[register] = value
			index++
		case op.INVOKE:
			fun := vm.Stack.Pop()
			arg_count := int(program[index])
			index++
			arguments := make([]ds.Value, arg_count, arg_count)
			for i := range arguments {
				arguments[i] = vm.Stack.Pop()
			}
			switch fun.(type) {
			case funcs.GoClosure:
				result := fun.(funcs.GoClosure).Function(arguments...)
				vm.Stack.Push(result)
			case funcs.SlangClosure:
				closure := fun.(funcs.SlangClosure)
				currentFrame.ContinueIndex = index
				index = int(closure.ProgramPosition)
				currentFrame = NewFrameFrom(currentFrame)
				// TODO Pass arguments and whatnot
			default:
				panic("Can't call a non-function")
			}
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
			value := vm.Stack.Pop()
			currentFrame = currentFrame.CallingFrame
			index = currentFrame.ContinueIndex
			vm.Stack.Push(value)
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
		case op.DEFINE:
			id := program[index]
			index++
			funcs.Defs[int(id)] = vm.Stack.Pop()
		case op.CLOSURE:
			start := uint(program[index])
			index++
			vm.Stack.Push(funcs.SlangClosure{ProgramPosition: start, IsProtocolMethod: false})
		default:
			panic(fmt.Errorf("Unknown instruction at %d: %d", index, program[index]))
		}
	}
}

func ParseStrings(instructions []byte, position *int) []string {
	startPosition := *position
	count := int(instructions[startPosition])
	startPosition++
	strings := make([]string, count, count)

	for index := 0; index < count; index++ {
		length := int(instructions[startPosition])
		startPosition++
		endPosition := startPosition + length
		strings[index] = string(instructions[startPosition:endPosition])
		startPosition = endPosition
	}
	*position = startPosition
	return strings
}

func ExpandDefsSlice(instructions []byte, position *int) {
	size := instructions[*position]
	*position++
	defs := make([]ds.Value, int(size), int(size))
	for i := range funcs.Defs {
		defs[i] = funcs.Defs[i]
	}
	funcs.Defs = defs
}

func main() {
	instructions, err := ioutil.ReadFile(os.Args[1])
	fmt.Println(instructions)
	if err != nil {
		panic(err)
	}
	coroutine := NewCoroutine()
	startIndex := 0
	strings := ParseStrings(instructions, &startIndex)
	ExpandDefsSlice(instructions, &startIndex)
	fmt.Println(funcs.Defs)
	prog := Program{Instructions: instructions[startIndex:], Strings: strings}
	coroutine.Program = &prog
	coroutine.Run(0)
}
