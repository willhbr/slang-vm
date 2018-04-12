package ds

import "math/big"

type Closure interface {
	IsBuiltin() bool
}

func NewInt64(value int64) Value {
	return big.NewInt(value)
}

func NewInt(value uint8) Value {
	return big.NewInt(int64(value))
}

type Type struct {
	Name            string
	ProtocolMethods map[int]Closure
	Attributes      []uint8
}

func (t Type) String() string {
	return t.Name
}

// TODO This just does a linear search, maybe work something better out?
func (t Type) GetAttrFrom(inst Instance, atomValue uint8) Value {
	for instanceIndex, atomIndex := range t.Attributes {
		if atomIndex == atomValue {
			return inst.Attributes[instanceIndex]
		}
	}
	panic("Could not get attribute from instance!")
}

func NewType(name string, attributes []uint8) *Type {
	return &Type{Name: name, Attributes: attributes}
}

type Instance struct {
	Type       *Type
	Attributes []Value
}

func NewInstance(t *Type, attributes []Value) *Instance {
	return &Instance{Type: t, Attributes: attributes}
}

type Module struct {
	Name string
}

type Atom int

func (m Module) String() string {
	return m.Name
}

var intType = Type{Name: "Int"}
var stringType = Type{Name: "String"}
var boolType = Type{Name: "Bool"}
var moduleType = Type{Name: "Module"}
var atomType = Type{Name: "Atom"}

var NilType = Type{Name: "Nil"}
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
