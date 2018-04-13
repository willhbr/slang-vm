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
// String: 11
StringType,
// Channel: 12
ChannelType,
// Channel.new: 13
GoClosure{Function: Channel__new},
// Channel.send: 14
GoClosure{Function: Channel__send},
// Channel.receive: 15
GoClosure{Function: Channel__receive},
// Sequence: 16
Module{Name: "Sequence"},
// Sequence.conj: 17
ProtocolClosure{ID: 17},
// Sequence.head: 18
ProtocolClosure{ID: 18},
// Sequence.tail: 19
ProtocolClosure{ID: 19},
// Enumerable: 20
Module{Name: "Enumerable"},
// Enumerable.reduce: 21
ProtocolClosure{ID: 21},
// List: 22
ListType,
// List.conj: 23
GoClosure{Function: List__conj},
// List.head: 24
GoClosure{Function: List__head},
// List.tail: 25
GoClosure{Function: List__tail},
// List.new: 26
GoClosure{Function: List__new},
// Vector: 27
VectorType,
// Vector.conj: 28
GoClosure{Function: Vector__conj},
// Vector.head: 29
GoClosure{Function: Vector__head},
// Vector.tail: 30
GoClosure{Function: Vector__tail},
}
