package ds

type Value int

//go:generate peds -pkg=ds -maps="Map<Value,Value>" -sets="Set<Value>" -vectors="Vector<Value>" -file=generated_collections.go

// Required for go generate it seems
func f() {
}
