package ds

type Type struct {
	name string
}

type Instance struct {
	Type *Type
}

func GetType(thing interface{}) *Type {
	// TODO get a type struct for the thing
	// Fallback on things that are Instances
	return &Type{name: "Unknown"}
}

// The one and only nil
// TODO make own type
var NilType = Type{name: "NilType"}

var Nil = Instance{Type: &NilType}

type Value interface{}

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
