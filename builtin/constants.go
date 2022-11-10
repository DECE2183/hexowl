package builtin

import "math"

var constants = map[string]interface{}{
	"nil":   nil,
	"inf":   math.Inf(0),
	"nan":   math.NaN(),
	"pi":    math.Pi,
	"e":     math.E,
	"true":  true,
	"false": false,
	"help":  "Type in the expression you want to calc and press Enter to get the result.\n\tTo define a variable type its name and assign the value with '=' operator.\n\tType 'funcs()' to see all available functions.\n\tType 'vars()' to see all available variables.",
}

func HasConstant(name string) bool {
	_, found := constants[name]
	return found
}

func RegisterConstant[T string | bool | uint64 | int64 | float64](name string, value T) {
	constants[name] = value
}

func GetConstant(name string) (val interface{}, found bool) {
	val, found = constants[name]
	return
}

func ListConstants() map[string]interface{} {
	return constants
}
