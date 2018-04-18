package types
import "math/big"
// Int: 16
var IntType = &Type{Name: "Int",
ProtocolMethods: map[int]Closure{
1: GoClosure{Function: Int___rArr_string},
}}
// String: 18
var StringType = &Type{Name: "String",
ProtocolMethods: map[int]Closure{
1: GoClosure{Function: String___rArr_string},
13: GoClosure{Function: String__get},
15: GoClosure{Function: String___eq_},
}}
// Atom: 22
var AtomType = &Type{Name: "Atom",
ProtocolMethods: map[int]Closure{
1: GoClosure{Function: Atom___rArr_string},
}}
// Channel: 27
var ChannelType = &Type{Name: "Channel",
ProtocolMethods: map[int]Closure{
}}
// List: 37
var ListType = &Type{Name: "List",
ProtocolMethods: map[int]Closure{
32: GoClosure{Function: List__conj},
33: GoClosure{Function: List__head},
34: GoClosure{Function: List__tail},
}}
// Vector: 42
var VectorType = &Type{Name: "Vector",
ProtocolMethods: map[int]Closure{
32: GoClosure{Function: Vector__conj},
33: GoClosure{Function: Vector__head},
34: GoClosure{Function: Vector__tail},
}}

    func GetType(object Value) *Type {
      switch object.(type) {
    
case *big.Int:
return IntType
case string:
return StringType
case Atom:
return AtomType
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
    
