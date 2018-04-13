package types
import "math/big"
// Int: 10
var IntType = &Type{Name: "Int",
ProtocolMethods: map[int]Closure{
1: GoClosure{Function: Int___rArr_string},
}}
// String: 12
var StringType = &Type{Name: "String",
ProtocolMethods: map[int]Closure{
1: GoClosure{Function: String___rArr_string},
}}
// Channel: 14
var ChannelType = &Type{Name: "Channel",
ProtocolMethods: map[int]Closure{
}}
// List: 24
var ListType = &Type{Name: "List",
ProtocolMethods: map[int]Closure{
19: GoClosure{Function: List__conj},
20: GoClosure{Function: List__head},
21: GoClosure{Function: List__tail},
}}
// Vector: 29
var VectorType = &Type{Name: "Vector",
ProtocolMethods: map[int]Closure{
19: GoClosure{Function: Vector__conj},
20: GoClosure{Function: Vector__head},
21: GoClosure{Function: Vector__tail},
}}

    func GetType(object Value) *Type {
      switch object.(type) {
    
case *big.Int:
return IntType
case string:
return StringType
case chan Value:
return ChannelType
case List:
return ListType
case Vector:
return VectorType

    case Instance:
      return object.(Instance).Type
    default:
      println(object)
      panic("Can't find type of object")
    }}
    
