package funcs
import "../ds"
var Defs = []ds.Value {
// IO: 0
ds.Module{Name: "IO"},
// IO.puts: 1
GoClosure{Function: IO__puts},
// Kernel: 2
ds.Module{Name: "Kernel"},
// Kernel.type: 3
GoClosure{Function: Kernel__type},
}
