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
			argCount := int(program[index])
			index++
		Finished:
			for {
				switch fun.(type) {
				case funcs.GoClosure:
					arguments := make([]ds.Value, argCount, argCount)
					size := len(arguments)
					for i := size - 1; i >= 0; i-- {
						arguments[i] = co.Stack.Pop()
					}
					result := fun.(funcs.GoClosure).Function(co, arguments...)
					co.Stack.Push(result)
					break Finished
				case funcs.SlangClosure:
					closure := fun.(funcs.SlangClosure)
					currentFrame.ContinueIndex = index
					index = int(closure.ProgramPosition)
					currentFrame = vm.NewFrameFrom(currentFrame)
					for i, value := range closure.Registers {
						currentFrame.Registers[i] = value
					}
					break Finished
				case funcs.ProtocolClosure:
					closure := fun.(funcs.ProtocolClosure)
					subject, ok := co.Stack.PeekFromTopMinus(argCount)
					if !ok {
						panic("Cannot call protocol method with no arguments!")
					}
					t := ds.GetType(subject)
					function, ok := t.ProtocolMethods[closure.ID]
					if !ok {
						panic(fmt.Errorf("%s does not implement protocol method", t.Name))
					}
					fun = function
				default:
					panic("Can't call a non-function")
				}
			}
		case op.SPAWN:
			skip := program[index]
			index++
			continueIndex := index
			index += int(skip)
			newCo := vm.NewCoroutine()
			newCo.Program = co.Program
			go Run(newCo, int(continueIndex))
		case op.APPLY:
			panic("Can't do APPLY yet")
		case op.CONST_A:
			value := program[index]
			index++
			co.Stack.Push(ds.Atom(value))
		case op.CONST_I:
			value := program[index]
			index++
			co.Stack.Push(ds.NewInt(value))
		case op.CONST_I_BIG:
			var value int64
			var read int64
			value = 0
			end := index + 8
			for ; index < end; index++ {
				read = int64(program[index])
				value = value | read
				value = value << 8
			}
			co.Stack.Push(ds.NewInt64(value))
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
		case op.JUMP_BACK:
			decrease := int(program[index])
			index++
			index -= decrease
		case op.AND:
			increase := program[index]
			index++
			cond := co.Stack.Pop()
			if cond != true {
				index += int(increase)
			}
		case op.RETURN:
			currentFrame = currentFrame.CallingFrame
			index = currentFrame.ContinueIndex
		case op.NEW_MAP:
			count := int(program[index])
			index++
			arguments := make([]ds.MapItem, count, count)
			for i := count - 1; i >= 0; i-- {
				item := ds.MapItem{}
				item.Value = co.Stack.Pop()
				item.Key = co.Stack.Pop()
				arguments[i] = item
			}
			m := ds.NewMap(arguments...)
			co.Stack.Push(m)
		case op.NEW_VECTOR:
			count := int(program[index])
			index++
			arguments := make([]ds.Value, count, count)
			for i := count - 1; i >= 0; i-- {
				arguments[i] = co.Stack.Pop()
			}
			v := ds.NewVector(arguments...)
			co.Stack.Push(v)
		case op.NEW_LIST:
			l := ds.NewList()
			co.Stack.Push(l)
		case op.DEFINE:
			id := program[index]
			index++
			funcs.Defs[int(id)] = co.Stack.Pop()
		case op.TYPE:
			id := program[index]
			index++
			nameIndex := int(program[index])
			index++
			attrCount := int(program[index])
			index++
			attributes := make([]uint8, attrCount, attrCount)
			for i := range attributes {
				attributes[i] = program[index]
				index++
			}
			newType := ds.NewType(co.Program.Strings[nameIndex], attributes)
			funcs.Defs[int(id)] = newType
		case op.INSTANCE:
			typeID := int(program[index])
			index++
			size := int(program[index])
			index++
			instType := funcs.Defs[typeID].(*ds.Type)
			attributes := make([]ds.Value, size, size)
			for i := size - 1; i >= 0; i-- {
				attributes[i] = co.Stack.Pop()
			}
			co.Stack.Push(ds.NewInstance(instType, attributes))
		case op.CLOSURE:
			start := uint(program[index])
			index++
			capturedCount := int(program[index])
			index++
			closure := funcs.NewSlangClosure(start)
			endAt := index + capturedCount
			for ; index < endAt; index++ {
				register := int(program[index])
				closure.Registers[register] = currentFrame.Registers[register]
			}
			co.Stack.Push(closure)
		case op.PROTOCOL_CLOSURE:
			id := int(program[index])
			index++
			closure := funcs.ProtocolClosure{ID: id}
			funcs.Defs[id] = closure
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
	if err != nil {
		panic(err)
	}
	coroutine := vm.NewCoroutine()
	startIndex := 0
	strings := ParseStrings(instructions, &startIndex)
	ExpandDefsSlice(instructions, &startIndex)
	prog := vm.Program{Instructions: instructions[startIndex:], Strings: strings}
	coroutine.Program = &prog
	Run(coroutine, 0)
}
