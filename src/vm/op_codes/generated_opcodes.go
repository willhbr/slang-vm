package op_codes
const (LOAD_LOCAL = 1
LOAD_DEF = 2
STORE = 3
INVOKE = 4
CONST_I = 5
CONST_I_BIG = 6
CONST_S = 7
CONST_A = 8
CONST_TRUE = 9
CONST_FALSE = 10
CONST_NIL = 11
JUMP = 12
JUMP_BACK = 13
AND = 14
RETURN = 15
CLOSURE = 16
PROTOCOL_CLOSURE = 17
NEW_MAP = 18
NEW_VECTOR = 19
NEW_LIST = 20
DEFINE = 21
TYPE = 22
INSTANCE = 23
IMPLEMENT = 24
RAISE = 25
TRY = 26
END_TRY = 27
DISCARD = 28
)
func ToString(code byte) string {
switch code {
case LOAD_LOCAL: return "LOAD_LOCAL"
case LOAD_DEF: return "LOAD_DEF"
case STORE: return "STORE"
case INVOKE: return "INVOKE"
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
case TYPE: return "TYPE"
case INSTANCE: return "INSTANCE"
case IMPLEMENT: return "IMPLEMENT"
case RAISE: return "RAISE"
case TRY: return "TRY"
case END_TRY: return "END_TRY"
case DISCARD: return "DISCARD"
default: return "UNKNOWN"
}}
