package types
// Channel: 9
var ChannelType = &Type{Name: "Channel",
ProtocolMethods: map[int]Closure{
}}
// List: 20
var ListType = &Type{Name: "List",
ProtocolMethods: map[int]Closure{
15: GoClosure{Function: List__conj},
16: GoClosure{Function: List__head},
17: GoClosure{Function: List__tail},
}}
// Vector: 25
var VectorType = &Type{Name: "Vector",
ProtocolMethods: map[int]Closure{
15: GoClosure{Function: Vector__conj},
16: GoClosure{Function: Vector__head},
17: GoClosure{Function: Vector__tail},
}}
