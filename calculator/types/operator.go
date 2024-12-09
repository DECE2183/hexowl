package types

type OperatorType int8

const (
	// Not an operator.
	O_NONE OperatorType = iota

	// Declare function operator.
	O_DECLFUNC
	// Execution sequence operator.
	O_SEQUENCE

	// Assign operator.
	O_ASSIGN
	// Local assign operator.
	O_ASSIGNLOCAL
	// Assign after subtraction operator.
	O_ASSIGNMINUS
	// Assign after addition operator.
	O_ASSIGNPLUS
	// Assign after multiplication operator.
	O_ASSIGNMUL
	// Assign after division operator.
	O_ASSIGNDIV
	// Assign after bitwise AND operator.
	O_ASSIGNBITAND
	// Assign after bitwise OR operator.
	O_ASSIGNBITOR

	// Enumerate values operator.
	O_ENUMERATE

	// Logic OR operator.
	O_LOGICOR
	// Logic AND operator.
	O_LOGICAND

	// Equal to comparison operator
	O_EQUALITY
	// Not equal to comparison operator.
	O_NOTEQ
	// Greater than comparison operator.
	O_MORE
	// Less than comparison operator.
	O_LESS
	// Greater than or equal to comparison operator.
	O_MOREEQ
	// Less than or equal to comparison operator.
	O_LESSEQ

	// Addition operator.
	O_PLUS
	// Subtraction operator.
	O_MINUS
	// Multiplication operator.
	O_MULTIPLY
	// Division operator.
	O_DIVIDE
	// Modulo operator.
	O_MODULO
	// Exponentiation operator.
	O_POWER

	// Bitwise OR operator.
	O_BITOR
	// Bitwise AND operator.
	O_BITAND
	// Bitwise XOR operator.
	O_BITXOR
	// Set bit to zero operator.
	O_BITCLEAR
	// Bitwise inverse operator.
	O_BITINVERSE
	// Bit shift left operator.
	O_LEFTSHIFT
	// Bit shift right operator.
	O_RIGHTSHIFT

	// Logic NOT operator.
	O_LOGICNOT
	// Non zero bits count operator.
	O_POPCNT

	O_LOCALVAR
	O_USERVAR
	O_CONSTANT
	O_USERFUNC
	O_BUILTINFUNC

	// Call function operator.
	O_CALLFUNC

	// Operator count (utility value).
	O_COUNT
	// Flow control operator.
	O_FLOW OperatorType = -1
)

var stringToOperatorMap = map[string]OperatorType{
	"->": O_DECLFUNC,
	";":  O_SEQUENCE,

	"=":  O_ASSIGN,
	":=": O_ASSIGNLOCAL,
	"-=": O_ASSIGNMINUS,
	"+=": O_ASSIGNPLUS,
	"*=": O_ASSIGNMUL,
	"/=": O_ASSIGNDIV,
	"&=": O_ASSIGNBITAND,
	"|=": O_ASSIGNBITOR,

	",": O_ENUMERATE,

	"||": O_LOGICOR,
	"&&": O_LOGICAND,
	"==": O_EQUALITY,
	"!=": O_NOTEQ,

	">":  O_MORE,
	"<":  O_LESS,
	">=": O_MOREEQ,
	"<=": O_LESSEQ,

	"+":  O_PLUS,
	"-":  O_MINUS,
	"*":  O_MULTIPLY,
	"/":  O_DIVIDE,
	"%":  O_MODULO,
	"**": O_POWER,

	"|":  O_BITOR,
	"&":  O_BITAND,
	"^":  O_BITXOR,
	"&^": O_BITCLEAR,
	"&~": O_BITCLEAR,
	"~":  O_BITINVERSE,
	"<<": O_LEFTSHIFT,
	">>": O_RIGHTSHIFT,

	"!": O_LOGICNOT,
	"#": O_POPCNT,
}

func ParseOperator(str string) OperatorType {
	return stringToOperatorMap[str]
}

func (op OperatorType) IsUnary() bool {
	return op == O_BITINVERSE || op == O_POPCNT || op == O_LOGICNOT
}

func (op OperatorType) IsAssign() bool {
	return op >= O_ASSIGN && op <= O_ASSIGNBITOR
}
