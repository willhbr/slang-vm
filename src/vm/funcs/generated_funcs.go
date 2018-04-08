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
// IO: 4
ds.Module{Name: "IO"},
// IO.gets: 5
GoClosure{Function: IO__gets},
// IO.puts: 6
GoClosure{Function: IO__puts},
// Kernel: 7
ds.Module{Name: "Kernel"},
// Kernel.*: 8
GoClosure{Function: Kernel__times},
// Kernel.-: 9
GoClosure{Function: Kernel__minus},
// Kernel.<: 10
GoClosure{Function: Kernel__lessThan},
// Kernel.conj: 11
GoClosure{Function: Kernel__conj},
// Kernel.type: 12
GoClosure{Function: Kernel__type},
}
