package ds

import "math/big"

func NewInt64(value int64) Value {
	return big.NewInt(value)
}

func NewInt(value uint8) Value {
	return big.NewInt(int64(value))
}

type Type struct {
	name string
}

func (t Type) String() string {
	return t.name
}

type Instance struct {
	Type *Type
}

type Module struct {
	Name string
}

type Atom int

func (m Module) String() string {
	return m.Name
}

var intType = Type{name: "Int"}
var stringType = Type{name: "String"}
var boolType = Type{name: "Bool"}
var moduleType = Type{name: "Module"}
var atomType = Type{name: "Atom"}

var NilType = Type{name: "Nil"}
var Nil = Instance{Type: &NilType}

type Value interface{}

func GetType(thing interface{}) *Type {
	switch thing.(type) {
	case int:
		return &intType
	case string:
		return &stringType
	case bool:
		return &boolType
	case Instance:
		return thing.(Instance).Type
	case Module:
		return &moduleType
	case Atom:
		return &atomType
	default:
		panic("Cannot get type of variable")
	}
}

//go:generate peds -pkg=ds -maps="Map<Value,Value>" -sets="Set<Value>" -vectors="Vector<Value>" -file=generated_collections.go

type List struct {
	value *Value
	next  *List
}

var emptyList = List{}

func NewList() List {
	return List{}
}

func (l List) IsEmpty() bool {
	return l.value != nil
}

func (l List) Tail() List {
	if l.IsEmpty() {
		return emptyList
	} else {
		return *l.next
	}
}

func (l List) Value() Value {
	if l.IsEmpty() {
		return Nil
	} else {
		return *l.value
	}
}

func (l List) Cons(value Value) List {
	return List{value: &value, next: &l}
}
