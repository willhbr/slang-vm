package types

import "math/big"

type Channel chan Value
type Int *big.Int
type String string

var NilType = &Type{Name: "Nil"}
var Nil = Instance{Type: NilType}
