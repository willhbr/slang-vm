package types
// Channel: 9
var ChannelType = &Type{Name: "Channel",
ProtocolMethods: map[int]Closure{
}}
// Vector: 19
var VectorType = &Type{Name: "Vector",
ProtocolMethods: map[int]Closure{
14: GoClosure{Function: Vector__conj},
15: GoClosure{Function: Vector__head},
16: GoClosure{Function: Vector__tail},
}}
