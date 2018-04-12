package funcs
import "../ds"
var Defs = []ds.Value {
// Channel: 0
ds.Module{Name: "Channel"},
// Channel.new: 1
GoClosure{Function: Channel__new},
// Channel.receive: 2
GoClosure{Function: Channel__receive},
// Channel.send: 3
GoClosure{Function: Channel__send},
// Enumerable: 4
ds.Module{Name: "Enumerable"},
// Enumerable.reduce: 5 (proto method)
ProtocolClosure{ID: 5},
// IO: 6
ds.Module{Name: "IO"},
// IO.gets: 7
GoClosure{Function: IO__gets},
// IO.puts: 8
GoClosure{Function: IO__puts},
// Kernel: 9
ds.Module{Name: "Kernel"},
// Kernel.*: 10
GoClosure{Function: Kernel__times},
// Kernel.-: 11
GoClosure{Function: Kernel__minus},
// Kernel.<: 12
GoClosure{Function: Kernel__lessThan},
// Kernel.conj: 13
GoClosure{Function: Kernel__conj},
// Kernel.type: 14
GoClosure{Function: Kernel__type},
}
