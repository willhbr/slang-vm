package utils


func ReadInteger(program []byte, index int) (byte, int) {
  return program[index], index + 1
}

func ReadString(program []byte, index int) (byte, int) {
  return program[index], index + 1
}

func ReadAtom(program []byte, index int) (byte, int) {
  return program[index], index + 1
}

func ReadGlobal(program []byte, index int) (byte, int) {
  return program[index], index + 1
}

func ReadArgCount(program []byte, index int) (byte, int) {
  return program[index], index + 1
}

func ReadOffset(program []byte, index int) (byte, int) {
  return program[index], index + 1
}

func ReadPosition(program []byte, index int) (byte, int) {
  return program[index], index + 1
}

func ReadDefCount(program []byte, index int) (byte, int) {
  return program[index], index + 1
}

func ReadStringLen(program []byte, index int) (byte, int) {
  return program[index], index + 1
}


func ReadBigInteger(program []byte, index int) (int64, int) {
  var value int64
  var read int64
  value = 0
  end := index + 8
  for ; index < end; index++ {
    read = int64(program[index])
    value = value | read
    value = value << 8
  }
  return value, index
}
