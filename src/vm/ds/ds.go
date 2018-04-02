package ds

// The one and only nil
// TODO make own type
const Nil = 0

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
