package types
var Defs = []Value {
// Printable: 0
Module{Name: "Printable"},
// Printable.->string: 1
ProtocolClosure{ID: 1},
// IO: 2
Module{Name: "IO"},
// IO.puts: 3
GoClosure{Function: IO__puts},
// IO.gets: 4
GoClosure{Function: IO__gets},
// Kernel: 5
Module{Name: "Kernel"},
// Kernel.type: 6
GoClosure{Function: Kernel__type},
// Kernel.<: 7
GoClosure{Function: Kernel__lessThan},
// Kernel.-: 8
GoClosure{Function: Kernel__minus},
// Kernel.*: 9
GoClosure{Function: Kernel__times},
// Int: 10
IntType,
// Int.->string: 11
GoClosure{Function: Int___rArr_string},
// String: 12
StringType,
// String.->string: 13
GoClosure{Function: String___rArr_string},
// Atom: 14
AtomType,
// Atom.->string: 15
GoClosure{Function: Atom___rArr_string},
// Atom.value: 16
GoClosure{Function: Atom__value},
// Channel: 17
ChannelType,
// Channel.new: 18
GoClosure{Function: Channel__new},
// Channel.send: 19
GoClosure{Function: Channel__send},
// Channel.receive: 20
GoClosure{Function: Channel__receive},
// Sequence: 21
Module{Name: "Sequence"},
// Sequence.conj: 22
ProtocolClosure{ID: 22},
// Sequence.head: 23
ProtocolClosure{ID: 23},
// Sequence.tail: 24
ProtocolClosure{ID: 24},
// Enumerable: 25
Module{Name: "Enumerable"},
// Enumerable.reduce: 26
ProtocolClosure{ID: 26},
// List: 27
ListType,
// List.conj: 28
GoClosure{Function: List__conj},
// List.head: 29
GoClosure{Function: List__head},
// List.tail: 30
GoClosure{Function: List__tail},
// List.new: 31
GoClosure{Function: List__new},
// Vector: 32
VectorType,
// Vector.conj: 33
GoClosure{Function: Vector__conj},
// Vector.head: 34
GoClosure{Function: Vector__head},
// Vector.tail: 35
GoClosure{Function: Vector__tail},
}
