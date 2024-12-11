package builtin

import (
	"math"

	"github.com/dece2183/hexowl/v2/types"
)

const _VERSION = "2.0.0-alpha"
const _HELP_MESSAGE = `Type in the expression you want to calc and press Enter to get the result.
	
	| Operator               | Syntax  |
	|------------------------|---------|
	| Positive bits count    | #       |
	| Bitwise NOT            | ~       |
	| Bitclear (AND NOT)     | &~ &^   |
	| Bitwise XOR            | ^       |
	| Bitwise AND            | &       |
	| Bitwise OR             | |       |
	| Right shift            | >>      |
	| Left shift             | <<      |
	| Modulo                 | %       |
	| Division               | /       |
	| Exponentiation         | **      |
	| Multiplication         | *       |
	| Subtraction            | -       |
	| Addition               | +       |
	| Logical NOT            | !       |
	| Less or equal          | <=      |
	| More or equal          | >=      |
	| Less                   | <       |
	| More                   | >       |
	| Not equal              | !=      |
	| Equal                  | ==      |
	| Logical AND            | &&      |
	| Logical OR             | ||      |
	| Enumerate              | ,       |
	| Bitwise OR and assign  | |=      |
	| Bitwise AND and assign | &=      |
	| Divide and assign      | /=      |
	| Mutiply and assign     | *=      |
	| Add and assign         | +=      |
	| Subtract and assign    | -=      |
	| Local assign           | :=      |
	| Assign                 | =       |
	| Sequence               | ;       |
	| Declare function       | ->      |

	To define a variable type its name and assign the value with '=' operator.
	To define a function use the '->' operator. For example: 'foo(x) -> x*x'.
	You can sequence multiple expressions using the ';' operator.
	The result of the last expression in the sequence will be returned.
	
	Type 'funcs()' to see all available functions.
	Type 'vars()' to see all available variables.`

var constants = types.BuiltinConstantMap{
	"nil":     nil,
	"inf":     math.Inf(0),
	"nan":     math.NaN(),
	"pi":      math.Pi,
	"e":       math.E,
	"true":    true,
	"false":   false,
	"help":    _HELP_MESSAGE,
	"version": _VERSION,
}

func RegisterConstants(ctx *types.Context) {
	for key, val := range constants {
		ctx.Builtin.RegisterConstant(key, val)
	}
}
