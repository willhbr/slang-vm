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
// Channel: 14
ChannelType,
// Channel.new: 15
GoClosure{Function: Channel__new},
// Channel.send: 16
GoClosure{Function: Channel__send},
// Channel.receive: 17
GoClosure{Function: Channel__receive},
// Sequence: 18
Module{Name: "Sequence"},
// Sequence.conj: 19
ProtocolClosure{ID: 19},
// Sequence.head: 20
ProtocolClosure{ID: 20},
// Sequence.tail: 21
ProtocolClosure{ID: 21},
// Enumerable: 22
Module{Name: "Enumerable"},
// Enumerable.reduce: 23
ProtocolClosure{ID: 23},
// List: 24
ListType,
// List.conj: 25
GoClosure{Function: List__conj},
// List.head: 26
GoClosure{Function: List__head},
// List.tail: 27
GoClosure{Function: List__tail},
// List.new: 28
GoClosure{Function: List__new},
// Vector: 29
VectorType,
// Vector.conj: 30
GoClosure{Function: Vector__conj},
// Vector.head: 31
GoClosure{Function: Vector__head},
// Vector.tail: 32
GoClosure{Function: Vector__tail},
}
