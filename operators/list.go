package operators

const (
	_VAL_CONST valueType = iota
	_VAL_LOCALVAR
	_VAL_USERVAR
	_VAL_CONSTANT
	_VAL_USERFUNC
	_VAL_BUILTINFUNC
	_VAL_LOCALFUNCPTR
	_VAL_FUNCPTR
)

var valStringRepresent = map[valueType]string{
	_VAL_CONST:        "constant",
	_VAL_LOCALVAR:     "local variable",
	_VAL_USERVAR:      "user variable",
	_VAL_CONSTANT:     "built-in constant",
	_VAL_USERFUNC:     "user function",
	_VAL_BUILTINFUNC:  "built-in function",
	_VAL_LOCALFUNCPTR: "local pointer to function",
	_VAL_FUNCPTR:      "pointer to function",
}

// Operator types
const (
	OP_NONE operatorType = iota

	OP_DECLFUNC
	OP_SEQUENCE

	OP_ASSIGN
	OP_ASSIGNLOCAL
	OP_ASSIGNMINUS
	OP_ASSIGNPLUS
	OP_ASSIGNMUL
	OP_ASSIGNDIV
	OP_ASSIGNBITAND
	OP_ASSIGNBITOR

	OP_ENUMERATE

	OP_LOGICOR
	OP_LOGICAND
	OP_EQUALITY
	OP_NOTEQ

	OP_MORE
	OP_LESS
	OP_MOREEQ
	OP_LESSEQ

	OP_PLUS
	OP_MINUS
	OP_MULTIPLY
	OP_DIVIDE
	OP_MODULO
	OP_POWER

	OP_BITOR
	OP_BITAND
	OP_BITXOR
	OP_BITCLEAR
	OP_BITINVERSE
	OP_LEFTSHIFT
	OP_RIGHTSHIFT

	OP_LOGICNOT
	OP_POPCNT

	OP_LOCALVAR
	OP_USERVAR
	OP_CONSTANT
	OP_USERFUNC
	OP_BUILTINFUNC

	OP_CALLFUNC

	OP_COUNT
	OP_FLOW operatorType = -1
)

var opStringRepresent = map[string]operatorType{
	"->": OP_DECLFUNC,
	";":  OP_SEQUENCE,

	"=":  OP_ASSIGN,
	":=": OP_ASSIGNLOCAL,
	"-=": OP_ASSIGNMINUS,
	"+=": OP_ASSIGNPLUS,
	"*=": OP_ASSIGNMUL,
	"/=": OP_ASSIGNDIV,
	"&=": OP_ASSIGNBITAND,
	"|=": OP_ASSIGNBITOR,

	",": OP_ENUMERATE,

	"||": OP_LOGICOR,
	"&&": OP_LOGICAND,
	"==": OP_EQUALITY,
	"!=": OP_NOTEQ,

	">":  OP_MORE,
	"<":  OP_LESS,
	">=": OP_MOREEQ,
	"<=": OP_LESSEQ,

	"+":  OP_PLUS,
	"-":  OP_MINUS,
	"*":  OP_MULTIPLY,
	"/":  OP_DIVIDE,
	"%":  OP_MODULO,
	"**": OP_POWER,

	"|":  OP_BITOR,
	"&":  OP_BITAND,
	"^":  OP_BITXOR,
	"&^": OP_BITCLEAR,
	"&~": OP_BITCLEAR,
	"~":  OP_BITINVERSE,
	"<<": OP_LEFTSHIFT,
	">>": OP_RIGHTSHIFT,

	"!": OP_LOGICNOT,
	"#": OP_POPCNT,
}

func (op operatorType) IsUnary() bool {
	return op == OP_BITINVERSE || op == OP_POPCNT || op == OP_LOGICNOT
}

func (op operatorType) IsAssign() bool {
	return op >= OP_ASSIGN && op <= OP_ASSIGNBITOR
}
