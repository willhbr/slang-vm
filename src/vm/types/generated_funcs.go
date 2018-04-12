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
// Sequence.conj: 14
ProtocolClosure{ID: 14},
// Sequence.head: 15
ProtocolClosure{ID: 15},
// Sequence.tail: 16
ProtocolClosure{ID: 16},
// Enumerable: 17
Module{Name: "Enumerable"},
// Enumerable.reduce: 18
ProtocolClosure{ID: 18},
// Vector: 19
VectorType,
// Vector.conj: 20
GoClosure{Function: Vector__conj},
// Vector.head: 21
GoClosure{Function: Vector__head},
// Vector.tail: 22
GoClosure{Function: Vector__tail},
}
