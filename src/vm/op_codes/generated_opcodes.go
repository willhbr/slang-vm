package op_codes
const (LOAD_LOCAL = 1
LOAD_DEF = 2
STORE = 3
INVOKE = 4
APPLY = 5
CONST_I = 6
CONST_S = 7
CONST_A = 8
CONST_TRUE = 9
CONST_FALSE = 10
CONST_NIL = 11
JUMP = 12
AND = 13
OR = 14
RETURN = 15
CLOSURE = 16
NEW_MAP = 17
NEW_VECTOR = 18
NEW_LIST = 19
DEFINE = 20
CONS = 21
INSERT = 22
)
func ToString(code byte) string {
switch code {
case LOAD_LOCAL: return "LOAD_LOCAL"
case LOAD_DEF: return "LOAD_DEF"
case STORE: return "STORE"
case INVOKE: return "INVOKE"
case APPLY: return "APPLY"
case CONST_I: return "CONST_I"
case CONST_S: return "CONST_S"
case CONST_A: return "CONST_A"
case CONST_TRUE: return "CONST_TRUE"
case CONST_FALSE: return "CONST_FALSE"
case CONST_NIL: return "CONST_NIL"
case JUMP: return "JUMP"
case AND: return "AND"
case OR: return "OR"
case RETURN: return "RETURN"
case CLOSURE: return "CLOSURE"
case NEW_MAP: return "NEW_MAP"
case NEW_VECTOR: return "NEW_VECTOR"
case NEW_LIST: return "NEW_LIST"
case DEFINE: return "DEFINE"
case CONS: return "CONS"
case INSERT: return "INSERT"
default: return "UNKNOWN"
}}
