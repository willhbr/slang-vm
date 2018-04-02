package op_codes
const (LOAD = 1
STORE = 2
DISPATCH = 3
APPLY = 4
CONST_I = 5
CONST_S = 6
JUMP = 7
AND = 8
OR = 9
RETURN = 10
NEW_MAP = 11
NEW_VECTOR = 12
NEW_LIST = 13
CONS = 14
INSERT = 15
)
func ToString(code byte) string {
switch code {
case LOAD: return "LOAD"
case STORE: return "STORE"
case DISPATCH: return "DISPATCH"
case APPLY: return "APPLY"
case CONST_I: return "CONST_I"
case CONST_S: return "CONST_S"
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
