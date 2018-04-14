package main

import (
	op "./op_codes"
	"./types"
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
			co.Stack.Push(types.Defs[program[index]])
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
				case types.GoClosure:
					arguments := make([]types.Value, argCount, argCount)
					size := len(arguments)
					for i := size - 1; i >= 0; i-- {
						arguments[i] = co.Stack.Pop()
					}
					result, err := fun.(types.GoClosure).Function(arguments...)
					// lol no errors
					if err != nil {
						// TODO jump to next error handle
						panic(err)
					}
					co.Stack.Push(result)
					break Finished
				case types.SlangClosure:
					closure := fun.(types.SlangClosure)
					currentFrame.ContinueIndex = index
					index = int(closure.ProgramPosition)
					currentFrame = vm.NewFrameFrom(currentFrame)
					currentFrame.StackPosition = co.Stack.Len()
					for i, value := range closure.Registers {
						currentFrame.Registers[i] = value
					}
					break Finished
				case types.ProtocolClosure:
					closure := fun.(types.ProtocolClosure)
					subject, ok := co.Stack.PeekFromTopMinus(argCount)
					if !ok {
						panic("Cannot call protocol method with no arguments!")
					}
					t := types.GetType(subject)
					function, ok := t.ProtocolMethods[closure.ID]
					if !ok {
						panic(fmt.Errorf("%s does not implement protocol method", t.Name))
					}
					fun = function
				default:
					panic("Can't call a non-function")
				}
			}
		case op.CONST_A:
			value := program[index]
			index++
			co.Stack.Push(types.Atom(value))
		case op.CONST_I:
			value := program[index]
			index++
			co.Stack.Push(types.NewInt(value))
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
			co.Stack.Push(types.NewInt64(value))
		case op.CONST_S:
			idx := program[index]
			index++
			str := strings[idx]
			co.Stack.Push(types.Value(str))
		case op.CONST_TRUE:
			co.Stack.Push(types.Value(true))
		case op.CONST_FALSE:
			co.Stack.Push(types.Value(false))
		case op.CONST_NIL:
			co.Stack.Push(types.Nil)
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
			if cond == false {
				index += int(increase)
			} else {
				asNil, ok := cond.(types.Instance)
				if ok && asNil.Type == types.NilType {
					index += int(increase)
				}
			}
		case op.RETURN:
			currentFrame = currentFrame.CallingFrame
			index = currentFrame.ContinueIndex
		case op.RAISE:
			// TODO Record some kind of error trace here
			frame := currentFrame
			err := co.Stack.Pop()
			for {
				if frame.CatchIndexes.IsEmpty() {
					frame = frame.CallingFrame
					if frame == nil {
						panic(types.NewSlangError(err))
					}
				} else {
					catchIdx := frame.CatchIndexes.Pop()
					index = catchIdx.Index
					co.Stack.Trim(catchIdx.StackSize)
					co.Stack.Push(err)
					break
				}
			}
		case op.TRY:
			offset := int(program[index])
			index++
			currentFrame.CatchIndexes.Push(vm.CatchIndex{
				Index:     index + offset,
				StackSize: co.Stack.Len(),
			})
		case op.END_TRY:
			currentFrame.CatchIndexes.Pop()
		case op.NEW_MAP:
			count := int(program[index])
			index++
			arguments := make([]types.MapItem, count, count)
			for i := count - 1; i >= 0; i-- {
				item := types.MapItem{}
				item.Value = co.Stack.Pop()
				item.Key = co.Stack.Pop()
				arguments[i] = item
			}
			m := types.NewMap(arguments...)
			co.Stack.Push(m)
		case op.NEW_VECTOR:
			count := int(program[index])
			index++
			arguments := make([]types.Value, count, count)
			for i := count - 1; i >= 0; i-- {
				arguments[i] = co.Stack.Pop()
			}
			v := types.NewVector(arguments...)
			co.Stack.Push(v)
		case op.NEW_LIST:
			l := types.NewList()
			co.Stack.Push(l)
		case op.DEFINE:
			id := program[index]
			index++
			types.Defs[int(id)] = co.Stack.Pop()
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
			newType := types.NewType(co.Program.Strings[nameIndex], attributes)
			types.Defs[int(id)] = newType
		case op.INSTANCE:
			typeID := int(program[index])
			index++
			size := int(program[index])
			index++
			instType := types.Defs[typeID].(*types.Type)
			attributes := make([]types.Value, size, size)
			for i := size - 1; i >= 0; i-- {
				attributes[i] = co.Stack.Pop()
			}
			co.Stack.Push(types.NewInstance(instType, attributes))
		case op.IMPLEMENT:
			panic("Not implemented")
		case op.CLOSURE:
			start := uint(program[index])
			index++
			capturedCount := int(program[index])
			index++
			closure := types.NewSlangClosure(start)
			endAt := index + capturedCount
			for ; index < endAt; index++ {
				register := int(program[index])
				closure.Registers[register] = currentFrame.Registers[register]
			}
			co.Stack.Push(closure)
		case op.PROTOCOL_CLOSURE:
			id := int(program[index])
			index++
			closure := types.ProtocolClosure{ID: id}
			types.Defs[id] = closure
		case op.DISCARD:
			co.Stack.Pop()
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
	defs := make([]types.Value, int(size), int(size))
	for i := range types.Defs {
		defs[i] = types.Defs[i]
	}
	types.Defs = defs
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
	if !coroutine.Stack.IsEmpty() {
		fmt.Println(coroutine.Stack)
	}
}
