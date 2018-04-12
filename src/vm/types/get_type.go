package types

import "math/big"

var intType = Type{Name: "Int"}
var stringType = Type{Name: "String"}
var boolType = Type{Name: "Bool"}
var moduleType = Type{Name: "Module"}
var atomType = Type{Name: "Atom"}
var listType = Type{Name: "List"}
var mapType = Type{Name: "Map"}

var NilType = Type{Name: "Nil"}
var Nil = Instance{Type: &NilType}

func GetType(thing interface{}) *Type {
	switch thing.(type) {
	case *big.Int:
		return &intType
	case string:
		return &stringType
	case bool:
		return &boolType
	case *Instance:
		return thing.(*Instance).Type
	case *Module:
		return &moduleType
	case Atom:
		return &atomType
	case *Vector:
		return VectorType
	case *List:
		return &listType
	case *Map:
		return &mapType
	default:
		panic("Cannot get type of variable")
	}
}
