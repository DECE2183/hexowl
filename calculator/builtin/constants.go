package builtin

import (
	"math"

	"github.com/dece2183/hexowl/calculator/types"
)

var constants = types.BuiltinConstantMap{
	"nil":     nil,
	"inf":     math.Inf(0),
	"nan":     math.NaN(),
	"pi":      math.Pi,
	"e":       math.E,
	"true":    true,
	"false":   false,
	"help":    "Type in the expression you want to calc and press Enter to get the result.\n\tTo define a variable type its name and assign the value with '=' operator.\n\tType 'funcs()' to see all available functions.\n\tType 'vars()' to see all available variables.",
	"version": "2.0.0",
}

func RegisterConstants(ctx *types.Context) {
	for key, val := range constants {
		ctx.Builtin.RegisterConstant(key, val)
	}
}
