package types

type TokenType uint8

const (
	// Not a token.
	T_NONE TokenType = iota
	// Number in scientific notation.
	T_NUM_SCI
	// Number in decimal representation.
	T_NUM_DEC
	// Number in hexadecimal representation.
	T_NUM_HEX
	// Number in binary representation.
	T_NUM_BIN
	// Some variable, constant or function.
	T_UNIT
	// Operator.
	T_OP
	// Flow control (brackets).
	T_CTL
	// Detected function call.
	T_FUNC
	// String.
	T_STR
)

type Token struct {
	Type    TokenType
	Literal string
}
