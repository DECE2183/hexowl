package builtin

import (
	"math"

	"github.com/dece2183/hexowl/builtin/types"
)

var constants = types.ConstantMap{
	"nil":     nil,
	"inf":     math.Inf(0),
	"nan":     math.NaN(),
	"pi":      math.Pi,
	"e":       math.E,
	"true":    true,
	"false":   false,
	"help":    "Type in the expression you want to calc and press Enter to get the result.\n\tTo define a variable type its name and assign the value with '=' operator.\n\tType 'funcs()' to see all available functions.\n\tType 'vars()' to see all available variables.",
	"version": "1.4.2",
}

// Is constant with name presented in the builtin constant map.
func HasConstant(name string) bool {
	_, found := constants[name]
	return found
}

// Register a new constant and add it to the builtin constant map.
func RegisterConstant[T string | bool | uint64 | int64 | float64](name string, value T) {
	constants[name] = value
}

// Get constant by name from the builtin constant map.
func GetConstant(name string) (val interface{}, found bool) {
	val, found = constants[name]
	return
}

// Return the builtin constant map.
func ListConstants() types.ConstantMap {
	return constants
}
