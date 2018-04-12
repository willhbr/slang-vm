package op_codes
const (LOAD_LOCAL = 1
LOAD_DEF = 2
STORE = 3
INVOKE = 4
APPLY = 5
CONST_I = 6
CONST_I_BIG = 7
CONST_S = 8
CONST_A = 9
CONST_TRUE = 10
CONST_FALSE = 11
CONST_NIL = 12
JUMP = 13
JUMP_BACK = 14
AND = 15
RETURN = 16
CLOSURE = 17
PROTOCOL_CLOSURE = 23
NEW_MAP = 18
NEW_VECTOR = 19
NEW_LIST = 20
DEFINE = 21
SPAWN = 22
TYPE = 24
INSTANCE = 25
)
func ToString(code byte) string {
switch code {
case LOAD_LOCAL: return "LOAD_LOCAL"
case LOAD_DEF: return "LOAD_DEF"
case STORE: return "STORE"
case INVOKE: return "INVOKE"
case APPLY: return "APPLY"
case CONST_I: return "CONST_I"
case CONST_I_BIG: return "CONST_I_BIG"
case CONST_S: return "CONST_S"
case CONST_A: return "CONST_A"
case CONST_TRUE: return "CONST_TRUE"
case CONST_FALSE: return "CONST_FALSE"
case CONST_NIL: return "CONST_NIL"
case JUMP: return "JUMP"
case JUMP_BACK: return "JUMP_BACK"
case AND: return "AND"
case RETURN: return "RETURN"
case CLOSURE: return "CLOSURE"
case PROTOCOL_CLOSURE: return "PROTOCOL_CLOSURE"
case NEW_MAP: return "NEW_MAP"
case NEW_VECTOR: return "NEW_VECTOR"
case NEW_LIST: return "NEW_LIST"
case DEFINE: return "DEFINE"
case SPAWN: return "SPAWN"
case TYPE: return "TYPE"
case INSTANCE: return "INSTANCE"
default: return "UNKNOWN"
}}
