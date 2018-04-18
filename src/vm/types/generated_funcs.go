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
// Kernel.+: 9
GoClosure{Function: Kernel___plus_},
// Kernel./: 10
GoClosure{Function: Kernel___div_},
// Kernel.*: 11
GoClosure{Function: Kernel__times},
// Access: 12
Module{Name: "Access"},
// Access.get: 13
ProtocolClosure{ID: 13},
// Equatable: 14
Module{Name: "Equatable"},
// Equatable.=: 15
ProtocolClosure{ID: 15},
// Int: 16
IntType,
// Int.->string: 17
GoClosure{Function: Int___rArr_string},
// String: 18
StringType,
// String.->string: 19
GoClosure{Function: String___rArr_string},
// String.get: 20
GoClosure{Function: String__get},
// String.=: 21
GoClosure{Function: String___eq_},
// Atom: 22
AtomType,
// Atom.->string: 23
GoClosure{Function: Atom___rArr_string},
// Atom.value: 24
GoClosure{Function: Atom__value},
// File: 25
Module{Name: "File"},
// File.read: 26
GoClosure{Function: File__read},
// Channel: 27
ChannelType,
// Channel.new: 28
GoClosure{Function: Channel__new},
// Channel.send: 29
GoClosure{Function: Channel__send},
// Channel.receive: 30
GoClosure{Function: Channel__receive},
// Sequence: 31
Module{Name: "Sequence"},
// Sequence.conj: 32
ProtocolClosure{ID: 32},
// Sequence.head: 33
ProtocolClosure{ID: 33},
// Sequence.tail: 34
ProtocolClosure{ID: 34},
// Enumerable: 35
Module{Name: "Enumerable"},
// Enumerable.reduce: 36
ProtocolClosure{ID: 36},
// List: 37
ListType,
// List.conj: 38
GoClosure{Function: List__conj},
// List.head: 39
GoClosure{Function: List__head},
// List.tail: 40
GoClosure{Function: List__tail},
// List.new: 41
GoClosure{Function: List__new},
// Vector: 42
VectorType,
// Vector.conj: 43
GoClosure{Function: Vector__conj},
// Vector.head: 44
GoClosure{Function: Vector__head},
// Vector.tail: 45
GoClosure{Function: Vector__tail},
}
