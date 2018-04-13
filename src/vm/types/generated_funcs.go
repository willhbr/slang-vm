package types
var Defs = []Value {
// IO: 0
Module{Name: "IO"},
// IO.puts: 1
GoClosure{Function: IO__puts},
// IO.gets: 2
GoClosure{Function: IO__gets},
// Kernel: 3
Module{Name: "Kernel"},
// Kernel.type: 4
GoClosure{Function: Kernel__type},
// Kernel.<: 5
GoClosure{Function: Kernel__lessThan},
// Kernel.-: 6
GoClosure{Function: Kernel__minus},
// Kernel.*: 7
GoClosure{Function: Kernel__times},
// Kernel.conj: 8
GoClosure{Function: Kernel__conj},
// Channel: 9
ChannelType,
// Channel.new: 10
GoClosure{Function: Channel__new},
// Channel.send: 11
GoClosure{Function: Channel__send},
// Channel.receive: 12
GoClosure{Function: Channel__receive},
// Sequence: 13
Module{Name: "Sequence"},
// Sequence.cons: 14
ProtocolClosure{ID: 14},
// Sequence.conj: 15
ProtocolClosure{ID: 15},
// Sequence.head: 16
ProtocolClosure{ID: 16},
// Sequence.tail: 17
ProtocolClosure{ID: 17},
// Enumerable: 18
Module{Name: "Enumerable"},
// Enumerable.reduce: 19
ProtocolClosure{ID: 19},
// List: 20
ListType,
// List.new: 21
GoClosure{Function: List__new},
// List.conj: 22
GoClosure{Function: List__conj},
// List.head: 23
GoClosure{Function: List__head},
// List.tail: 24
GoClosure{Function: List__tail},
// Vector: 25
VectorType,
// Vector.conj: 26
GoClosure{Function: Vector__conj},
// Vector.head: 27
GoClosure{Function: Vector__head},
// Vector.tail: 28
GoClosure{Function: Vector__tail},
}
