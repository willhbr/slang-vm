package op_codes
const (LOAD = 1
STORE = 2
DISPATCH = 3
APPLY = 4
CONST_I = 5
CONST_S = 6
CONST_TRUE = 7
CONST_FALSE = 8
CONST_NIL = 9
JUMP = 10
AND = 11
OR = 12
RETURN = 13
NEW_MAP = 14
NEW_VECTOR = 15
NEW_LIST = 16
CONS = 17
INSERT = 18
)
func ToString(code byte) string {
switch code {
case LOAD: return "LOAD"
case STORE: return "STORE"
case DISPATCH: return "DISPATCH"
case APPLY: return "APPLY"
case CONST_I: return "CONST_I"
case CONST_S: return "CONST_S"
case CONST_TRUE: return "CONST_TRUE"
case CONST_FALSE: return "CONST_FALSE"
case CONST_NIL: return "CONST_NIL"
case JUMP: return "JUMP"
case AND: return "AND"
case OR: return "OR"
case RETURN: return "RETURN"
case NEW_MAP: return "NEW_MAP"
case NEW_VECTOR: return "NEW_VECTOR"
case NEW_LIST: return "NEW_LIST"
case CONS: return "CONS"
case INSERT: return "INSERT"
default: return "UNKNOWN"
}}
