package types
// Int: 10
var IntType = &Type{Name: "Int",
ProtocolMethods: map[int]Closure{
}}
// String: 11
var StringType = &Type{Name: "String",
ProtocolMethods: map[int]Closure{
}}
// Channel: 12
var ChannelType = &Type{Name: "Channel",
ProtocolMethods: map[int]Closure{
}}
// List: 22
var ListType = &Type{Name: "List",
ProtocolMethods: map[int]Closure{
17: GoClosure{Function: List__conj},
18: GoClosure{Function: List__head},
19: GoClosure{Function: List__tail},
}}
// Vector: 27
var VectorType = &Type{Name: "Vector",
ProtocolMethods: map[int]Closure{
17: GoClosure{Function: Vector__conj},
18: GoClosure{Function: Vector__head},
19: GoClosure{Function: Vector__tail},
}}

    func GetType(object Value) *Type {
      switch object.(type) {
    

      case Int:
        return IntType
      

      case String:
        return StringType
      

      case Channel:
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
    
