package types

import "math/big"
import "fmt"

type Closure interface {
	IsBuiltin() bool
}

type SlangError struct {
	Cause Value
}

func (s SlangError) Error() string {
	return fmt.Sprintf("%+v", s.Cause)
}

func NewSlangError(cause Value) error {
	return SlangError{Cause: cause}
}

func NewInt64(value int64) *big.Int {
	return big.NewInt(value)
}

func NewInt(value uint8) *big.Int {
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

type Value interface{}

//go:generate peds -pkg=types -maps="Map<Value,Value>" -sets="Set<Value>" -vectors="VectorImpl<Value>" -file=generated_collections.go

type Vector interface {
	PushBack(values ...Value) Vector
	Get(index int) Value
	Update(index int, value Value) Vector
	GetSlice(start, end int) Vector
	Len() int
}

func (v VectorImpl) PushBack(values ...Value) Vector {
	return v.Append(values)
}

func (v VectorImpl) Update(index int, value Value) Vector {
	return v.Set(index, value)
}

func (v VectorImpl) GetSlice(start, end int) Vector {
	return v.Slice(start, end)
}

func (v VectorImplSlice) PushBack(values ...Value) Vector {
	return v.Append(values)
}

func (v VectorImplSlice) GetSlice(start, end int) Vector {
	return v.Slice(start, end)
}

func (v VectorImplSlice) Update(index int, value Value) Vector {
	return v.Set(index, value)
}

func NewVector(args ...Value) Vector {
	return NewVectorImpl(args...)
}

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

func (l List) Head() Value {
	if l.IsEmpty() {
		return Nil
	} else {
		return *l.value
	}
}

func (l List) Cons(value Value) *List {
	return &List{value: &value, next: &l}
}
