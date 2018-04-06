package main

import (
	"./ds"
	"./funcs"
	op "./op_codes"
	"./vm"
	"fmt"
	"io/ioutil"
	"os"
)

func Run(co *vm.Coroutine, startIndex int) {
	program := co.Program.Instructions
	strings := co.Program.Strings
	size := len(program)
	index := startIndex
	currentFrame := co.CurrentFrame
	for index < size {
		operation := program[index]
		// fmt.Println(op.ToString(operation))
		index++
		switch operation {
		case op.LOAD_LOCAL:
			co.Stack.Push(currentFrame.Registers[program[index]])
			index++
		case op.LOAD_DEF:
			co.Stack.Push(funcs.Defs[program[index]])
			index++
		case op.STORE:
			value := co.Stack.Pop()
			register := program[index]
			currentFrame.Registers[register] = value
			index++
		case op.INVOKE:
			fun := co.Stack.Pop()
			arg_count := int(program[index])
			index++
			arguments := make([]ds.Value, arg_count, arg_count)
			for i := range arguments {
				arguments[i] = co.Stack.Pop()
			}
			switch fun.(type) {
			case funcs.GoClosure:
				result := fun.(funcs.GoClosure).Function(co, arguments...)
				co.Stack.Push(result)
			case funcs.SlangClosure:
				closure := fun.(funcs.SlangClosure)
				currentFrame.ContinueIndex = index
				index = int(closure.ProgramPosition)
				currentFrame = vm.NewFrameFrom(currentFrame)
				// TODO Pass arguments and whatnot
			default:
				panic("Can't call a non-function")
			}
		case op.APPLY:
			panic("Can't do APPLY yet")
		case op.CONST_A:
			value := program[index]
			index++
			co.Stack.Push(ds.Atom(value))
		case op.CONST_I:
			value := program[index]
			index++
			co.Stack.Push(ds.Value(int(value)))
		case op.CONST_S:
			idx := program[index]
			index++
			str := strings[idx]
			co.Stack.Push(ds.Value(str))
		case op.CONST_TRUE:
			co.Stack.Push(ds.Value(true))
		case op.CONST_FALSE:
			co.Stack.Push(ds.Value(false))
		case op.CONST_NIL:
			co.Stack.Push(ds.Nil)
		case op.JUMP:
			increase := int(program[index])
			index++
			index += increase
		case op.AND:
			increase := program[index]
			index++
			if co.Stack.Pop() != 0 {
				index += int(increase)
			}
		case op.OR:
			panic("Can't do OR yet")
		case op.RETURN:
			value := co.Stack.Pop()
			currentFrame = currentFrame.CallingFrame
			index = currentFrame.ContinueIndex
			co.Stack.Push(value)
		case op.NEW_MAP:
			m := ds.NewMap()
			co.Stack.Push(m)
		case op.NEW_VECTOR:
			v := ds.NewVector()
			co.Stack.Push(v)
		case op.NEW_LIST:
			l := ds.NewList()
			co.Stack.Push(l)
		case op.CONS:
			panic("Can't do CONS yet")
		case op.INSERT:
			panic("Can't do INSERT yet")
		case op.DEFINE:
			id := program[index]
			index++
			funcs.Defs[int(id)] = co.Stack.Pop()
		case op.CLOSURE:
			start := uint(program[index])
			index++
			co.Stack.Push(funcs.SlangClosure{ProgramPosition: start, IsProtocolMethod: false})
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
	coroutine := vm.NewCoroutine()
	startIndex := 0
	strings := ParseStrings(instructions, &startIndex)
	ExpandDefsSlice(instructions, &startIndex)
	fmt.Println(funcs.Defs)
	prog := vm.Program{Instructions: instructions[startIndex:], Strings: strings}
	coroutine.Program = &prog
	Run(coroutine, 0)
}
