package types

import (
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"math"
)

const shiftSize = 5
const nodeSize = 32
const shiftBitMask = 0x1F

type commonNode interface{}

var emptyCommonNode commonNode = []commonNode{}

func uintMin(a, b uint) uint {
	if a < b {
		return a
	}

	return b
}

func newPath(shift uint, node commonNode) commonNode {
	if shift == 0 {
		return node
	}

	return newPath(shift-shiftSize, commonNode([]commonNode{node}))
}

func assertSliceOk(start, stop, len int) {
	if start < 0 {
		panic(fmt.Sprintf("Invalid slice index %d (index must be non-negative)", start))
	}

	if start > stop {
		panic(fmt.Sprintf("Invalid slice index: %d > %d", start, stop))
	}

	if stop > len {
		panic(fmt.Sprintf("Slice bounds out of range, start=%d, stop=%d, len=%d", start, stop, len))
	}
}

const upperMapLoadFactor float64 = 8.0
const lowerMapLoadFactor float64 = 2.0
const initialMapLoadFactor float64 = (upperMapLoadFactor + lowerMapLoadFactor) / 2

//////////////////////////
//// Hash functions //////
//////////////////////////

func hash(x []byte) uint32 {
	return crc32.ChecksumIEEE(x)
}

func interfaceHash(x interface{}) uint32 {
	return hash([]byte(fmt.Sprintf("%v", x)))
}

func byteHash(x byte) uint32 {
	return hash([]byte{x})
}

func uint8Hash(x uint8) uint32 {
	return byteHash(byte(x))
}

func int8Hash(x int8) uint32 {
	return uint8Hash(uint8(x))
}

func uint16Hash(x uint16) uint32 {
	bX := make([]byte, 2)
	binary.LittleEndian.PutUint16(bX, x)
	return hash(bX)
}

func int16Hash(x int16) uint32 {
	return uint16Hash(uint16(x))
}

func uint32Hash(x uint32) uint32 {
	bX := make([]byte, 4)
	binary.LittleEndian.PutUint32(bX, x)
	return hash(bX)
}

func int32Hash(x int32) uint32 {
	return uint32Hash(uint32(x))
}

func uint64Hash(x uint64) uint32 {
	bX := make([]byte, 8)
	binary.LittleEndian.PutUint64(bX, x)
	return hash(bX)
}

func int64Hash(x int64) uint32 {
	return uint64Hash(uint64(x))
}

func intHash(x int) uint32 {
	return int64Hash(int64(x))
}

func uintHash(x uint) uint32 {
	return uint64Hash(uint64(x))
}

func boolHash(x bool) uint32 {
	if x {
		return 1
	}

	return 0
}

func runeHash(x rune) uint32 {
	return int32Hash(int32(x))
}

func stringHash(x string) uint32 {
	return hash([]byte(x))
}

func float64Hash(x float64) uint32 {
	return uint64Hash(math.Float64bits(x))
}

func float32Hash(x float32) uint32 {
	return uint32Hash(math.Float32bits(x))
}

//////////////
/// Vector ///
//////////////

// A Vector is an ordered persistent/immutable collection of items corresponding roughly
// to the use cases for a slice.
type Vector struct {
	tail  []Value
	root  commonNode
	len   uint
	shift uint
}

var emptyVectorTail = make([]Value, 0)
var emptyVector *Vector = &Vector{root: emptyCommonNode, shift: shiftSize, tail: emptyVectorTail}

// NewVector returns a new Vector containing the items provided in items.
func NewVector(items ...Value) *Vector {
	return emptyVector.Append(items...)
}

// Get returns the element at position i.
func (v *Vector) Get(i int) Value {
	if i < 0 || uint(i) >= v.len {
		panic("Index out of bounds")
	}

	return v.sliceFor(uint(i))[i&shiftBitMask]
}

func (v *Vector) sliceFor(i uint) []Value {
	if i >= v.tailOffset() {
		return v.tail
	}

	node := v.root
	for level := v.shift; level > 0; level -= shiftSize {
		node = node.([]commonNode)[(i>>level)&shiftBitMask]
	}

	return node.([]Value)
}

func (v *Vector) tailOffset() uint {
	if v.len < nodeSize {
		return 0
	}

	return ((v.len - 1) >> shiftSize) << shiftSize
}

// Set returns a new vector with the element at position i set to item.
func (v *Vector) Set(i int, item Value) *Vector {
	if i < 0 || uint(i) >= v.len {
		panic("Index out of bounds")
	}

	if uint(i) >= v.tailOffset() {
		newTail := make([]Value, len(v.tail))
		copy(newTail, v.tail)
		newTail[i&shiftBitMask] = item
		return &Vector{root: v.root, tail: newTail, len: v.len, shift: v.shift}
	}

	return &Vector{root: v.doAssoc(v.shift, v.root, uint(i), item), tail: v.tail, len: v.len, shift: v.shift}
}

func (v *Vector) doAssoc(level uint, node commonNode, i uint, item Value) commonNode {
	if level == 0 {
		ret := make([]Value, nodeSize)
		copy(ret, node.([]Value))
		ret[i&shiftBitMask] = item
		return ret
	}

	ret := make([]commonNode, nodeSize)
	copy(ret, node.([]commonNode))
	subidx := (i >> level) & shiftBitMask
	ret[subidx] = v.doAssoc(level-shiftSize, ret[subidx], i, item)
	return ret
}

func (v *Vector) pushTail(level uint, parent commonNode, tailNode []Value) commonNode {
	subIdx := ((v.len - 1) >> level) & shiftBitMask
	parentNode := parent.([]commonNode)
	ret := make([]commonNode, subIdx+1)
	copy(ret, parentNode)
	var nodeToInsert commonNode

	if level == shiftSize {
		nodeToInsert = tailNode
	} else if subIdx < uint(len(parentNode)) {
		nodeToInsert = v.pushTail(level-shiftSize, parentNode[subIdx], tailNode)
	} else {
		nodeToInsert = newPath(level-shiftSize, tailNode)
	}

	ret[subIdx] = nodeToInsert
	return ret
}

// Append returns a new vector with item(s) appended to it.
func (v *Vector) Append(item ...Value) *Vector {
	result := v
	itemLen := uint(len(item))
	for insertOffset := uint(0); insertOffset < itemLen; {
		tailLen := result.len - result.tailOffset()
		tailFree := nodeSize - tailLen
		if tailFree == 0 {
			result = result.pushLeafNode(result.tail)
			result.tail = emptyVector.tail
			tailFree = nodeSize
			tailLen = 0
		}

		batchLen := uintMin(itemLen-insertOffset, tailFree)
		newTail := make([]Value, 0, tailLen+batchLen)
		newTail = append(newTail, result.tail...)
		newTail = append(newTail, item[insertOffset:insertOffset+batchLen]...)
		result = &Vector{root: result.root, tail: newTail, len: result.len + batchLen, shift: result.shift}
		insertOffset += batchLen
	}

	return result
}

func (v *Vector) pushLeafNode(node []Value) *Vector {
	var newRoot commonNode
	newShift := v.shift

	// Root overflow?
	if (v.len >> shiftSize) > (1 << v.shift) {
		newNode := newPath(v.shift, node)
		newRoot = commonNode([]commonNode{v.root, newNode})
		newShift = v.shift + shiftSize
	} else {
		newRoot = v.pushTail(v.shift, v.root, node)
	}

	return &Vector{root: newRoot, tail: v.tail, len: v.len, shift: newShift}
}

// Slice returns a VectorSlice that refers to all elements [start,stop) in v.
func (v *Vector) Slice(start, stop int) *VectorSlice {
	assertSliceOk(start, stop, v.Len())
	return &VectorSlice{vector: v, start: start, stop: stop}
}

// Len returns the length of v.
func (v *Vector) Len() int {
	return int(v.len)
}

// Range calls f repeatedly passing it each element in v in order as argument until either
// all elements have been visited or f returns false.
func (v *Vector) Range(f func(Value) bool) {
	var currentNode []Value
	for i := uint(0); i < v.len; i++ {
		if i&shiftBitMask == 0 {
			currentNode = v.sliceFor(i)
		}

		if !f(currentNode[i&shiftBitMask]) {
			return
		}
	}
}

// ToNativeSlice returns a Go slice containing all elements of v
func (v *Vector) ToNativeSlice() []Value {
	result := make([]Value, 0, v.len)
	for i := uint(0); i < v.len; i += nodeSize {
		result = append(result, v.sliceFor(i)...)
	}

	return result
}

////////////////
//// Slice /////
////////////////

// VectorSlice is a slice type backed by a Vector.
type VectorSlice struct {
	vector      *Vector
	start, stop int
}

// NewVectorSlice returns a new NewVectorSlice containing the items provided in items.
func NewVectorSlice(items ...Value) *VectorSlice {
	return &VectorSlice{vector: emptyVector.Append(items...), start: 0, stop: len(items)}
}

// Len returns the length of s.
func (s *VectorSlice) Len() int {
	return s.stop - s.start
}

// Get returns the element at position i.
func (s *VectorSlice) Get(i int) Value {
	if i < 0 || s.start+i >= s.stop {
		panic("Index out of bounds")
	}

	return s.vector.Get(s.start + i)
}

// Set returns a new slice with the element at position i set to item.
func (s *VectorSlice) Set(i int, item Value) *VectorSlice {
	if i < 0 || s.start+i >= s.stop {
		panic("Index out of bounds")
	}

	return s.vector.Set(s.start+i, item).Slice(s.start, s.stop)
}

// Append returns a new slice with item(s) appended to it.
func (s *VectorSlice) Append(items ...Value) *VectorSlice {
	newSlice := VectorSlice{vector: s.vector, start: s.start, stop: s.stop + len(items)}

	// If this is v slice that has an upper bound that is lower than the backing
	// vector then set the values in the backing vector to achieve some structural
	// sharing.
	itemPos := 0
	for ; s.stop+itemPos < s.vector.Len() && itemPos < len(items); itemPos++ {
		newSlice.vector = newSlice.vector.Set(s.stop+itemPos, items[itemPos])
	}

	// For the rest just append it to the underlying vector
	newSlice.vector = newSlice.vector.Append(items[itemPos:]...)
	return &newSlice
}

// Slice returns a VectorSlice that refers to all elements [start,stop) in s.
func (s *VectorSlice) Slice(start, stop int) *VectorSlice {
	assertSliceOk(start, stop, s.stop-s.start)
	return &VectorSlice{vector: s.vector, start: s.start + start, stop: s.start + stop}
}

// Range calls f repeatedly passing it each element in s in order as argument until either
// all elements have been visited or f returns false.
func (s *VectorSlice) Range(f func(Value) bool) {
	var currentNode []Value
	for i := uint(s.start); i < uint(s.stop); i++ {
		if i&shiftBitMask == 0 || i == uint(s.start) {
			currentNode = s.vector.sliceFor(uint(i))
		}

		if !f(currentNode[i&shiftBitMask]) {
			return
		}
	}
}

///////////
/// Map ///
///////////

//////////////////////
/// Backing vector ///
//////////////////////

type privateMapItemBucketVector struct {
	tail  []privateMapItemBucket
	root  commonNode
	len   uint
	shift uint
}

type MapItem struct {
	Key   Value
	Value Value
}

type privateMapItemBucket []MapItem

var emptyMapItemBucketVectorTail = make([]privateMapItemBucket, 0)
var emptyMapItemBucketVector *privateMapItemBucketVector = &privateMapItemBucketVector{root: emptyCommonNode, shift: shiftSize, tail: emptyMapItemBucketVectorTail}

func (v *privateMapItemBucketVector) Get(i int) privateMapItemBucket {
	if i < 0 || uint(i) >= v.len {
		panic("Index out of bounds")
	}

	return v.sliceFor(uint(i))[i&shiftBitMask]
}

func (v *privateMapItemBucketVector) sliceFor(i uint) []privateMapItemBucket {
	if i >= v.tailOffset() {
		return v.tail
	}

	node := v.root
	for level := v.shift; level > 0; level -= shiftSize {
		node = node.([]commonNode)[(i>>level)&shiftBitMask]
	}

	return node.([]privateMapItemBucket)
}

func (v *privateMapItemBucketVector) tailOffset() uint {
	if v.len < nodeSize {
		return 0
	}

	return ((v.len - 1) >> shiftSize) << shiftSize
}

func (v *privateMapItemBucketVector) Set(i int, item privateMapItemBucket) *privateMapItemBucketVector {
	if i < 0 || uint(i) >= v.len {
		panic("Index out of bounds")
	}

	if uint(i) >= v.tailOffset() {
		newTail := make([]privateMapItemBucket, len(v.tail))
		copy(newTail, v.tail)
		newTail[i&shiftBitMask] = item
		return &privateMapItemBucketVector{root: v.root, tail: newTail, len: v.len, shift: v.shift}
	}

	return &privateMapItemBucketVector{root: v.doAssoc(v.shift, v.root, uint(i), item), tail: v.tail, len: v.len, shift: v.shift}
}

func (v *privateMapItemBucketVector) doAssoc(level uint, node commonNode, i uint, item privateMapItemBucket) commonNode {
	if level == 0 {
		ret := make([]privateMapItemBucket, nodeSize)
		copy(ret, node.([]privateMapItemBucket))
		ret[i&shiftBitMask] = item
		return ret
	}

	ret := make([]commonNode, nodeSize)
	copy(ret, node.([]commonNode))
	subidx := (i >> level) & shiftBitMask
	ret[subidx] = v.doAssoc(level-shiftSize, ret[subidx], i, item)
	return ret
}

func (v *privateMapItemBucketVector) pushTail(level uint, parent commonNode, tailNode []privateMapItemBucket) commonNode {
	subIdx := ((v.len - 1) >> level) & shiftBitMask
	parentNode := parent.([]commonNode)
	ret := make([]commonNode, subIdx+1)
	copy(ret, parentNode)
	var nodeToInsert commonNode

	if level == shiftSize {
		nodeToInsert = tailNode
	} else if subIdx < uint(len(parentNode)) {
		nodeToInsert = v.pushTail(level-shiftSize, parentNode[subIdx], tailNode)
	} else {
		nodeToInsert = newPath(level-shiftSize, tailNode)
	}

	ret[subIdx] = nodeToInsert
	return ret
}

func (v *privateMapItemBucketVector) Append(item ...privateMapItemBucket) *privateMapItemBucketVector {
	result := v
	itemLen := uint(len(item))
	for insertOffset := uint(0); insertOffset < itemLen; {
		tailLen := result.len - result.tailOffset()
		tailFree := nodeSize - tailLen
		if tailFree == 0 {
			result = result.pushLeafNode(result.tail)
			result.tail = emptyMapItemBucketVector.tail
			tailFree = nodeSize
			tailLen = 0
		}

		batchLen := uintMin(itemLen-insertOffset, tailFree)
		newTail := make([]privateMapItemBucket, 0, tailLen+batchLen)
		newTail = append(newTail, result.tail...)
		newTail = append(newTail, item[insertOffset:insertOffset+batchLen]...)
		result = &privateMapItemBucketVector{root: result.root, tail: newTail, len: result.len + batchLen, shift: result.shift}
		insertOffset += batchLen
	}

	return result
}

func (v *privateMapItemBucketVector) pushLeafNode(node []privateMapItemBucket) *privateMapItemBucketVector {
	var newRoot commonNode
	newShift := v.shift

	// Root overflow?
	if (v.len >> shiftSize) > (1 << v.shift) {
		newNode := newPath(v.shift, node)
		newRoot = commonNode([]commonNode{v.root, newNode})
		newShift = v.shift + shiftSize
	} else {
		newRoot = v.pushTail(v.shift, v.root, node)
	}

	return &privateMapItemBucketVector{root: newRoot, tail: v.tail, len: v.len, shift: newShift}
}

func (v *privateMapItemBucketVector) Len() int {
	return int(v.len)
}

func (v *privateMapItemBucketVector) Range(f func(privateMapItemBucket) bool) {
	var currentNode []privateMapItemBucket
	for i := uint(0); i < v.len; i++ {
		if i&shiftBitMask == 0 {
			currentNode = v.sliceFor(uint(i))
		}

		if !f(currentNode[i&shiftBitMask]) {
			return
		}
	}
}

// Map is a persistent key - value map
type Map struct {
	backingVector *privateMapItemBucketVector
	len           int
}

func (m *Map) pos(key Value) int {
	return int(uint64(interfaceHash(key)) % uint64(m.backingVector.Len()))
}

// Helper type used during map creation and reallocation
type privateMapItemBuckets struct {
	buckets []privateMapItemBucket
	length  int
}

func newPrivateMapItemBuckets(itemCount int) *privateMapItemBuckets {
	size := int(float64(itemCount)/initialMapLoadFactor) + 1
	buckets := make([]privateMapItemBucket, size)
	return &privateMapItemBuckets{buckets: buckets}
}

func (b *privateMapItemBuckets) AddItem(item MapItem) {
	ix := int(uint64(interfaceHash(item.Key)) % uint64(len(b.buckets)))
	bucket := b.buckets[ix]
	if bucket != nil {
		// Hash collision, merge with existing bucket
		for keyIx, bItem := range bucket {
			if item.Key == bItem.Key {
				bucket[keyIx] = item
				return
			}
		}

		b.buckets[ix] = append(bucket, MapItem{Key: item.Key, Value: item.Value})
		b.length++
	} else {
		bucket := make(privateMapItemBucket, 0, int(math.Max(initialMapLoadFactor, 1.0)))
		b.buckets[ix] = append(bucket, item)
		b.length++
	}
}

func (b *privateMapItemBuckets) AddItemsFromMap(m *Map) {
	m.backingVector.Range(func(bucket privateMapItemBucket) bool {
		for _, item := range bucket {
			b.AddItem(item)
		}
		return true
	})
}

func newMap(items []MapItem) *Map {
	buckets := newPrivateMapItemBuckets(len(items))
	for _, item := range items {
		buckets.AddItem(item)
	}

	return &Map{backingVector: emptyMapItemBucketVector.Append(buckets.buckets...), len: buckets.length}
}

// Len returns the number of items in m.
func (m *Map) Len() int {
	return int(m.len)
}

// Load returns value identified by key. ok is set to true if key exists in the map, false otherwise.
func (m *Map) Load(key Value) (value Value, ok bool) {
	bucket := m.backingVector.Get(m.pos(key))
	if bucket != nil {
		for _, item := range bucket {
			if item.Key == key {
				return item.Value, true
			}
		}
	}

	var zeroValue Value
	return zeroValue, false
}

// Store returns a new Map containing value identified by key.
func (m *Map) Store(key Value, value Value) *Map {
	// Grow backing vector if load factor is too high
	if m.Len() >= m.backingVector.Len()*int(upperMapLoadFactor) {
		buckets := newPrivateMapItemBuckets(m.Len() + 1)
		buckets.AddItemsFromMap(m)
		buckets.AddItem(MapItem{Key: key, Value: value})
		return &Map{backingVector: emptyMapItemBucketVector.Append(buckets.buckets...), len: buckets.length}
	}

	pos := m.pos(key)
	bucket := m.backingVector.Get(pos)
	if bucket != nil {
		for ix, item := range bucket {
			if item.Key == key {
				// Overwrite existing item
				newBucket := make(privateMapItemBucket, len(bucket))
				copy(newBucket, bucket)
				newBucket[ix] = MapItem{Key: key, Value: value}
				return &Map{backingVector: m.backingVector.Set(pos, newBucket), len: m.len}
			}
		}

		// Add new item to bucket
		newBucket := make(privateMapItemBucket, len(bucket), len(bucket)+1)
		copy(newBucket, bucket)
		newBucket = append(newBucket, MapItem{Key: key, Value: value})
		return &Map{backingVector: m.backingVector.Set(pos, newBucket), len: m.len + 1}
	}

	item := MapItem{Key: key, Value: value}
	newBucket := privateMapItemBucket{item}
	return &Map{backingVector: m.backingVector.Set(pos, newBucket), len: m.len + 1}
}

// Delete returns a new Map without the element identified by key.
func (m *Map) Delete(key Value) *Map {
	pos := m.pos(key)
	bucket := m.backingVector.Get(pos)
	if bucket != nil {
		newBucket := make(privateMapItemBucket, 0)
		for _, item := range bucket {
			if item.Key != key {
				newBucket = append(newBucket, item)
			}
		}

		removedItemCount := len(bucket) - len(newBucket)
		if removedItemCount == 0 {
			return m
		}

		if len(newBucket) == 0 {
			newBucket = nil
		}

		newMap := &Map{backingVector: m.backingVector.Set(pos, newBucket), len: m.len - removedItemCount}
		if newMap.backingVector.Len() > 1 && newMap.Len() < newMap.backingVector.Len()*int(lowerMapLoadFactor) {
			// Shrink backing vector if needed to avoid occupying excessive space
			buckets := newPrivateMapItemBuckets(newMap.Len())
			buckets.AddItemsFromMap(newMap)
			return &Map{backingVector: emptyMapItemBucketVector.Append(buckets.buckets...), len: buckets.length}
		}

		return newMap
	}

	return m
}

// Range calls f repeatedly passing it each key and value as argument until either
// all elements have been visited or f returns false.
func (m *Map) Range(f func(Value, Value) bool) {
	m.backingVector.Range(func(bucket privateMapItemBucket) bool {
		for _, item := range bucket {
			if !f(item.Key, item.Value) {
				return false
			}
		}
		return true
	})
}

// ToNativeMap returns a native Go map containing all elements of m.
func (m *Map) ToNativeMap() map[Value]Value {
	result := make(map[Value]Value)
	m.Range(func(key Value, value Value) bool {
		result[key] = value
		return true
	})

	return result
}

////////////////////
/// Constructors ///
////////////////////

// NewMap returns a new Map containing all items in items.
func NewMap(items ...MapItem) *Map {
	return newMap(items)
}

// NewMapFromNativeMap returns a new Map containing all items in m.
func NewMapFromNativeMap(m map[Value]Value) *Map {
	buckets := newPrivateMapItemBuckets(len(m))
	for key, value := range m {
		buckets.AddItem(MapItem{Key: key, Value: value})
	}

	return &Map{backingVector: emptyMapItemBucketVector.Append(buckets.buckets...), len: buckets.length}
}

///////////
/// Map ///
///////////

//////////////////////
/// Backing vector ///
//////////////////////

type privateprivateSetMapItemBucketVector struct {
	tail  []privateprivateSetMapItemBucket
	root  commonNode
	len   uint
	shift uint
}

type privateSetMapItem struct {
	Key   Value
	Value struct{}
}

type privateprivateSetMapItemBucket []privateSetMapItem

var emptyprivateSetMapItemBucketVectorTail = make([]privateprivateSetMapItemBucket, 0)
var emptyprivateSetMapItemBucketVector *privateprivateSetMapItemBucketVector = &privateprivateSetMapItemBucketVector{root: emptyCommonNode, shift: shiftSize, tail: emptyprivateSetMapItemBucketVectorTail}

func (v *privateprivateSetMapItemBucketVector) Get(i int) privateprivateSetMapItemBucket {
	if i < 0 || uint(i) >= v.len {
		panic("Index out of bounds")
	}

	return v.sliceFor(uint(i))[i&shiftBitMask]
}

func (v *privateprivateSetMapItemBucketVector) sliceFor(i uint) []privateprivateSetMapItemBucket {
	if i >= v.tailOffset() {
		return v.tail
	}

	node := v.root
	for level := v.shift; level > 0; level -= shiftSize {
		node = node.([]commonNode)[(i>>level)&shiftBitMask]
	}

	return node.([]privateprivateSetMapItemBucket)
}

func (v *privateprivateSetMapItemBucketVector) tailOffset() uint {
	if v.len < nodeSize {
		return 0
	}

	return ((v.len - 1) >> shiftSize) << shiftSize
}

func (v *privateprivateSetMapItemBucketVector) Set(i int, item privateprivateSetMapItemBucket) *privateprivateSetMapItemBucketVector {
	if i < 0 || uint(i) >= v.len {
		panic("Index out of bounds")
	}

	if uint(i) >= v.tailOffset() {
		newTail := make([]privateprivateSetMapItemBucket, len(v.tail))
		copy(newTail, v.tail)
		newTail[i&shiftBitMask] = item
		return &privateprivateSetMapItemBucketVector{root: v.root, tail: newTail, len: v.len, shift: v.shift}
	}

	return &privateprivateSetMapItemBucketVector{root: v.doAssoc(v.shift, v.root, uint(i), item), tail: v.tail, len: v.len, shift: v.shift}
}

func (v *privateprivateSetMapItemBucketVector) doAssoc(level uint, node commonNode, i uint, item privateprivateSetMapItemBucket) commonNode {
	if level == 0 {
		ret := make([]privateprivateSetMapItemBucket, nodeSize)
		copy(ret, node.([]privateprivateSetMapItemBucket))
		ret[i&shiftBitMask] = item
		return ret
	}

	ret := make([]commonNode, nodeSize)
	copy(ret, node.([]commonNode))
	subidx := (i >> level) & shiftBitMask
	ret[subidx] = v.doAssoc(level-shiftSize, ret[subidx], i, item)
	return ret
}

func (v *privateprivateSetMapItemBucketVector) pushTail(level uint, parent commonNode, tailNode []privateprivateSetMapItemBucket) commonNode {
	subIdx := ((v.len - 1) >> level) & shiftBitMask
	parentNode := parent.([]commonNode)
	ret := make([]commonNode, subIdx+1)
	copy(ret, parentNode)
	var nodeToInsert commonNode

	if level == shiftSize {
		nodeToInsert = tailNode
	} else if subIdx < uint(len(parentNode)) {
		nodeToInsert = v.pushTail(level-shiftSize, parentNode[subIdx], tailNode)
	} else {
		nodeToInsert = newPath(level-shiftSize, tailNode)
	}

	ret[subIdx] = nodeToInsert
	return ret
}

func (v *privateprivateSetMapItemBucketVector) Append(item ...privateprivateSetMapItemBucket) *privateprivateSetMapItemBucketVector {
	result := v
	itemLen := uint(len(item))
	for insertOffset := uint(0); insertOffset < itemLen; {
		tailLen := result.len - result.tailOffset()
		tailFree := nodeSize - tailLen
		if tailFree == 0 {
			result = result.pushLeafNode(result.tail)
			result.tail = emptyprivateSetMapItemBucketVector.tail
			tailFree = nodeSize
			tailLen = 0
		}

		batchLen := uintMin(itemLen-insertOffset, tailFree)
		newTail := make([]privateprivateSetMapItemBucket, 0, tailLen+batchLen)
		newTail = append(newTail, result.tail...)
		newTail = append(newTail, item[insertOffset:insertOffset+batchLen]...)
		result = &privateprivateSetMapItemBucketVector{root: result.root, tail: newTail, len: result.len + batchLen, shift: result.shift}
		insertOffset += batchLen
	}

	return result
}

func (v *privateprivateSetMapItemBucketVector) pushLeafNode(node []privateprivateSetMapItemBucket) *privateprivateSetMapItemBucketVector {
	var newRoot commonNode
	newShift := v.shift

	// Root overflow?
	if (v.len >> shiftSize) > (1 << v.shift) {
		newNode := newPath(v.shift, node)
		newRoot = commonNode([]commonNode{v.root, newNode})
		newShift = v.shift + shiftSize
	} else {
		newRoot = v.pushTail(v.shift, v.root, node)
	}

	return &privateprivateSetMapItemBucketVector{root: newRoot, tail: v.tail, len: v.len, shift: newShift}
}

func (v *privateprivateSetMapItemBucketVector) Len() int {
	return int(v.len)
}

func (v *privateprivateSetMapItemBucketVector) Range(f func(privateprivateSetMapItemBucket) bool) {
	var currentNode []privateprivateSetMapItemBucket
	for i := uint(0); i < v.len; i++ {
		if i&shiftBitMask == 0 {
			currentNode = v.sliceFor(uint(i))
		}

		if !f(currentNode[i&shiftBitMask]) {
			return
		}
	}
}

// privateSetMap is a persistent key - value map
type privateSetMap struct {
	backingVector *privateprivateSetMapItemBucketVector
	len           int
}

func (m *privateSetMap) pos(key Value) int {
	return int(uint64(interfaceHash(key)) % uint64(m.backingVector.Len()))
}

// Helper type used during map creation and reallocation
type privateprivateSetMapItemBuckets struct {
	buckets []privateprivateSetMapItemBucket
	length  int
}

func newPrivateprivateSetMapItemBuckets(itemCount int) *privateprivateSetMapItemBuckets {
	size := int(float64(itemCount)/initialMapLoadFactor) + 1
	buckets := make([]privateprivateSetMapItemBucket, size)
	return &privateprivateSetMapItemBuckets{buckets: buckets}
}

func (b *privateprivateSetMapItemBuckets) AddItem(item privateSetMapItem) {
	ix := int(uint64(interfaceHash(item.Key)) % uint64(len(b.buckets)))
	bucket := b.buckets[ix]
	if bucket != nil {
		// Hash collision, merge with existing bucket
		for keyIx, bItem := range bucket {
			if item.Key == bItem.Key {
				bucket[keyIx] = item
				return
			}
		}

		b.buckets[ix] = append(bucket, privateSetMapItem{Key: item.Key, Value: item.Value})
		b.length++
	} else {
		bucket := make(privateprivateSetMapItemBucket, 0, int(math.Max(initialMapLoadFactor, 1.0)))
		b.buckets[ix] = append(bucket, item)
		b.length++
	}
}

func (b *privateprivateSetMapItemBuckets) AddItemsFromMap(m *privateSetMap) {
	m.backingVector.Range(func(bucket privateprivateSetMapItemBucket) bool {
		for _, item := range bucket {
			b.AddItem(item)
		}
		return true
	})
}

func newprivateSetMap(items []privateSetMapItem) *privateSetMap {
	buckets := newPrivateprivateSetMapItemBuckets(len(items))
	for _, item := range items {
		buckets.AddItem(item)
	}

	return &privateSetMap{backingVector: emptyprivateSetMapItemBucketVector.Append(buckets.buckets...), len: buckets.length}
}

// Len returns the number of items in m.
func (m *privateSetMap) Len() int {
	return int(m.len)
}

// Load returns value identified by key. ok is set to true if key exists in the map, false otherwise.
func (m *privateSetMap) Load(key Value) (value struct{}, ok bool) {
	bucket := m.backingVector.Get(m.pos(key))
	if bucket != nil {
		for _, item := range bucket {
			if item.Key == key {
				return item.Value, true
			}
		}
	}

	var zeroValue struct{}
	return zeroValue, false
}

// Store returns a new privateSetMap containing value identified by key.
func (m *privateSetMap) Store(key Value, value struct{}) *privateSetMap {
	// Grow backing vector if load factor is too high
	if m.Len() >= m.backingVector.Len()*int(upperMapLoadFactor) {
		buckets := newPrivateprivateSetMapItemBuckets(m.Len() + 1)
		buckets.AddItemsFromMap(m)
		buckets.AddItem(privateSetMapItem{Key: key, Value: value})
		return &privateSetMap{backingVector: emptyprivateSetMapItemBucketVector.Append(buckets.buckets...), len: buckets.length}
	}

	pos := m.pos(key)
	bucket := m.backingVector.Get(pos)
	if bucket != nil {
		for ix, item := range bucket {
			if item.Key == key {
				// Overwrite existing item
				newBucket := make(privateprivateSetMapItemBucket, len(bucket))
				copy(newBucket, bucket)
				newBucket[ix] = privateSetMapItem{Key: key, Value: value}
				return &privateSetMap{backingVector: m.backingVector.Set(pos, newBucket), len: m.len}
			}
		}

		// Add new item to bucket
		newBucket := make(privateprivateSetMapItemBucket, len(bucket), len(bucket)+1)
		copy(newBucket, bucket)
		newBucket = append(newBucket, privateSetMapItem{Key: key, Value: value})
		return &privateSetMap{backingVector: m.backingVector.Set(pos, newBucket), len: m.len + 1}
	}

	item := privateSetMapItem{Key: key, Value: value}
	newBucket := privateprivateSetMapItemBucket{item}
	return &privateSetMap{backingVector: m.backingVector.Set(pos, newBucket), len: m.len + 1}
}

// Delete returns a new privateSetMap without the element identified by key.
func (m *privateSetMap) Delete(key Value) *privateSetMap {
	pos := m.pos(key)
	bucket := m.backingVector.Get(pos)
	if bucket != nil {
		newBucket := make(privateprivateSetMapItemBucket, 0)
		for _, item := range bucket {
			if item.Key != key {
				newBucket = append(newBucket, item)
			}
		}

		removedItemCount := len(bucket) - len(newBucket)
		if removedItemCount == 0 {
			return m
		}

		if len(newBucket) == 0 {
			newBucket = nil
		}

		newMap := &privateSetMap{backingVector: m.backingVector.Set(pos, newBucket), len: m.len - removedItemCount}
		if newMap.backingVector.Len() > 1 && newMap.Len() < newMap.backingVector.Len()*int(lowerMapLoadFactor) {
			// Shrink backing vector if needed to avoid occupying excessive space
			buckets := newPrivateprivateSetMapItemBuckets(newMap.Len())
			buckets.AddItemsFromMap(newMap)
			return &privateSetMap{backingVector: emptyprivateSetMapItemBucketVector.Append(buckets.buckets...), len: buckets.length}
		}

		return newMap
	}

	return m
}

// Range calls f repeatedly passing it each key and value as argument until either
// all elements have been visited or f returns false.
func (m *privateSetMap) Range(f func(Value, struct{}) bool) {
	m.backingVector.Range(func(bucket privateprivateSetMapItemBucket) bool {
		for _, item := range bucket {
			if !f(item.Key, item.Value) {
				return false
			}
		}
		return true
	})
}

// ToNativeMap returns a native Go map containing all elements of m.
func (m *privateSetMap) ToNativeMap() map[Value]struct{} {
	result := make(map[Value]struct{})
	m.Range(func(key Value, value struct{}) bool {
		result[key] = value
		return true
	})

	return result
}

// Set is a persistent set
type Set struct {
	backingMap *privateSetMap
}

// NewSet returns a new Set containing items.
func NewSet(items ...Value) *Set {
	mapItems := make([]privateSetMapItem, 0, len(items))
	var mapValue struct{}
	for _, x := range items {
		mapItems = append(mapItems, privateSetMapItem{Key: x, Value: mapValue})
	}

	return &Set{backingMap: newprivateSetMap(mapItems)}
}

// Add returns a new Set containing item.
func (s *Set) Add(item Value) *Set {
	var mapValue struct{}
	return &Set{backingMap: s.backingMap.Store(item, mapValue)}
}

// Delete returns a new Set without item.
func (s *Set) Delete(item Value) *Set {
	newMap := s.backingMap.Delete(item)
	if newMap == s.backingMap {
		return s
	}

	return &Set{backingMap: newMap}
}

// Contains returns true if item is present in s, false otherwise.
func (s *Set) Contains(item Value) bool {
	_, ok := s.backingMap.Load(item)
	return ok
}

// Range calls f repeatedly passing it each element in s as argument until either
// all elements have been visited or f returns false.
func (s *Set) Range(f func(Value) bool) {
	s.backingMap.Range(func(k Value, _ struct{}) bool {
		return f(k)
	})
}

// IsSubset returns true if all elements in s are present in other, false otherwise.
func (s *Set) IsSubset(other *Set) bool {
	if other.Len() < s.Len() {
		return false
	}

	isSubset := true
	s.Range(func(item Value) bool {
		if !other.Contains(item) {
			isSubset = false
		}

		return isSubset
	})

	return isSubset
}

// IsSuperset returns true if all elements in other are present in s, false otherwise.
func (s *Set) IsSuperset(other *Set) bool {
	return other.IsSubset(s)
}

// Union returns a new Set containing all elements present
// in either s or other.
func (s *Set) Union(other *Set) *Set {
	result := s

	// Simplest possible solution right now. Would probable be more efficient
	// to concatenate two slices of elements from the two sets and create a
	// new set from that slice for many cases.
	other.Range(func(item Value) bool {
		result = result.Add(item)
		return true
	})

	return result
}

// Equals returns true if s and other contains the same elements, false otherwise.
func (s *Set) Equals(other *Set) bool {
	return s.Len() == other.Len() && s.IsSubset(other)
}

func (s *Set) difference(other *Set) []Value {
	items := make([]Value, 0)
	s.Range(func(item Value) bool {
		if !other.Contains(item) {
			items = append(items, item)
		}

		return true
	})

	return items
}

// Difference returns a new Set containing all elements present
// in s but not in other.
func (s *Set) Difference(other *Set) *Set {
	return NewSet(s.difference(other)...)
}

// SymmetricDifference returns a new Set containing all elements present
// in either s or other but not both.
func (s *Set) SymmetricDifference(other *Set) *Set {
	items := s.difference(other)
	items = append(items, other.difference(s)...)
	return NewSet(items...)
}

// Intersection returns a new Set containing all elements present in both
// s and other.
func (s *Set) Intersection(other *Set) *Set {
	items := make([]Value, 0)
	s.Range(func(item Value) bool {
		if other.Contains(item) {
			items = append(items, item)
		}

		return true
	})

	return NewSet(items...)
}

// Len returns the number of elements in s.
func (s *Set) Len() int {
	return s.backingMap.Len()
}

// ToNativeSlice returns a native Go slice containing all elements of s.
func (s *Set) ToNativeSlice() []Value {
	items := make([]Value, 0, s.Len())
	s.Range(func(item Value) bool {
		items = append(items, item)
		return true
	})

	return items
}
