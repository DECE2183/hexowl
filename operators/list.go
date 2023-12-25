package operators

// Operator types
const (
	OP_NONE operatorType = iota

	OP_DECLFUNC operatorType = iota
	OP_SEQUENCE operatorType = iota

	OP_ASSIGN       operatorType = iota
	OP_LOCALASSIGN  operatorType = iota
	OP_DECREMENT    operatorType = iota
	OP_INCREMENT    operatorType = iota
	OP_ASSIGNMUL    operatorType = iota
	OP_ASSIGNDIV    operatorType = iota
	OP_ASSIGNBITAND operatorType = iota
	OP_ASSIGNBITOR  operatorType = iota

	OP_ENUMERATE operatorType = iota

	OP_LOGICOR  operatorType = iota
	OP_LOGICAND operatorType = iota
	OP_EQUALITY operatorType = iota
	OP_NOTEQ    operatorType = iota

	OP_MORE   operatorType = iota
	OP_LESS   operatorType = iota
	OP_MOREEQ operatorType = iota
	OP_LESSEQ operatorType = iota

	OP_PLUS     operatorType = iota
	OP_MINUS    operatorType = iota
	OP_MULTIPLY operatorType = iota
	OP_DIVIDE   operatorType = iota
	OP_MODULO   operatorType = iota
	OP_POWER    operatorType = iota

	OP_BITOR      operatorType = iota
	OP_BITAND     operatorType = iota
	OP_BITXOR     operatorType = iota
	OP_BITCLEAR   operatorType = iota
	OP_BITINVERSE operatorType = iota
	OP_LEFTSHIFT  operatorType = iota
	OP_RIGHTSHIFT operatorType = iota

	OP_LOGICNOT operatorType = iota
	OP_POPCNT   operatorType = iota

	OP_LOCALVAR    operatorType = iota
	OP_USERVAR     operatorType = iota
	OP_CONSTANT    operatorType = iota
	OP_USERFUNC    operatorType = iota
	OP_BUILTINFUNC operatorType = iota

	OP_COUNT operatorType = iota
)

var opStringRepresent = map[string]operatorType{
	"->": OP_DECLFUNC,
	";":  OP_SEQUENCE,

	"=":  OP_ASSIGN,
	":=": OP_LOCALASSIGN,
	"-=": OP_DECREMENT,
	"+=": OP_INCREMENT,
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
