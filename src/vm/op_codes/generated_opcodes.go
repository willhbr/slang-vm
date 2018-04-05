package op_codes
const (LOAD_LOCAL = 1
LOAD_DEF = 21
STORE = 2
CALL_METHOD = 3
CALL_LOCAL = 19
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
DEFINE = 20
CONS = 17
INSERT = 18
)
func ToString(code byte) string {
switch code {
case LOAD_LOCAL: return "LOAD_LOCAL"
case LOAD_DEF: return "LOAD_DEF"
case STORE: return "STORE"
case CALL_METHOD: return "CALL_METHOD"
case CALL_LOCAL: return "CALL_LOCAL"
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
case DEFINE: return "DEFINE"
case CONS: return "CONS"
case INSERT: return "INSERT"
default: return "UNKNOWN"
}}
