package builtin

import "math"

var Constants = map[string]interface{}{
	"nil":   nil,
	"inf":   math.Inf(0),
	"nan":   math.NaN(),
	"pi":    math.Pi,
	"e":     math.E,
	"true":  true,
	"false": false,
	"help":  "Enter expression you want to calc and press Enter to get the result.\n\tTo define a variable type its name and assign the value with '=' operator.\n\tType 'funcs()' to see all available functions.\n\tType 'vars()' to see all available variables.",
}
