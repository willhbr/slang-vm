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
// IO.puts: 5
GoClosure{Function: IO__puts},
// Kernel: 6
ds.Module{Name: "Kernel"},
// Kernel.type: 7
GoClosure{Function: Kernel__type},
}
