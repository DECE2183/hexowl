package builtin

import (
	"github.com/dece2183/hexowl/v2/calculator/builtin/impl"
	"github.com/dece2183/hexowl/v2/calculator/types"
)

var functions = types.BuiltinFunctionMap{
	"sin": types.BuiltinFunction{
		Args: "(x)",
		Desc: "The sine of the radian argument x",
		Exec: impl.Sin,
	},
	"cos": types.BuiltinFunction{
		Args: "(x)",
		Desc: "The cosine of the radian argument x",
		Exec: impl.Cos,
	},
	"tan": types.BuiltinFunction{
		Args: "(x)",
		Desc: "The tangent of the radian argument x",
		Exec: impl.Tan,
	},
	"asin": types.BuiltinFunction{
		Args: "(x)",
		Desc: "The arcsine of the radian argument x",
		Exec: impl.Asin,
	},
	"acos": types.BuiltinFunction{
		Args: "(x)",
		Desc: "The arccosine of the radian argument x",
		Exec: impl.Acos,
	},
	"atan": types.BuiltinFunction{
		Args: "(x)",
		Desc: "The arctangent of the radian argument x",
		Exec: impl.Atan,
	},
	"pow": types.BuiltinFunction{
		Args: "(x,y)",
		Desc: "The base-x exponential of y",
		Exec: impl.Pow,
	},
	"sqrt": types.BuiltinFunction{
		Args: "(x)",
		Desc: "The square root of x",
		Exec: impl.Sqrt,
	},
	"exp": types.BuiltinFunction{
		Args: "(x)",
		Desc: "The base-e exponential of x",
		Exec: impl.Exp,
	},
	"logn": types.BuiltinFunction{
		Args: "(x)",
		Desc: "The natural logarithm of x",
		Exec: impl.Logn,
	},
	"log2": types.BuiltinFunction{
		Args: "(x)",
		Desc: "The binary logarithm of x",
		Exec: impl.Log2,
	},
	"log10": types.BuiltinFunction{
		Args: "(x)",
		Desc: "The decimal logarithm of x",
		Exec: impl.Log10,
	},
	"round": types.BuiltinFunction{
		Args: "(x)",
		Desc: "The nearest integer, rounding half away from zero",
		Exec: impl.Round,
	},
	"ceil": types.BuiltinFunction{
		Args: "(x)",
		Desc: "The least integer value greater than or equal to x",
		Exec: impl.Ceil,
	},
	"floor": types.BuiltinFunction{
		Args: "(x)",
		Desc: "The greatest integer value less than or equal to x",
		Exec: impl.Floor,
	},
	"rand": types.BuiltinFunction{
		Args: "(a,b)",
		Desc: "The random number in the range [a,b) or [0,1) if no arguments are passed",
		Exec: impl.Random,
	},
	"popcnt": {
		Args: "(x)",
		Desc: "The number of one bits (\"population count\") in x",
		Exec: impl.Popcount,
	},
	"vars": types.BuiltinFunction{
		Args: "()",
		Desc: "List available variables",
		Exec: impl.Vars,
	},
	"rmvar": types.BuiltinFunction{
		Args: "(name)",
		Desc: "Delete a specific user variable",
		Exec: impl.RmVar,
	},
	"clvars": types.BuiltinFunction{
		Args: "()",
		Desc: "Delete user defined variables",
		Exec: impl.ClearVars,
	},
	"funcs": types.BuiltinFunction{
		Args: "()",
		Desc: "List alailable functions",
		Exec: impl.Funcs,
	},
	"rmfunc": types.BuiltinFunction{
		Args: "(name)",
		Desc: "Delete a specific user function",
		Exec: impl.RmFunc,
	},
	"rmfuncvar": types.BuiltinFunction{
		Args: "(name,varid)",
		Desc: "Delete a specific user function variation",
		Exec: impl.RmFuncVar,
	},
	"clfuncs": types.BuiltinFunction{
		Args: "()",
		Desc: "Delete user defined functions",
		Exec: impl.ClearFuncs,
	},
	"save": types.BuiltinFunction{
		Args: "(id,comment)",
		Desc: "Save working environment with id and optional comment",
		Exec: impl.Save,
	},
	"load": types.BuiltinFunction{
		Args: "(id)",
		Desc: "Load working environment with id",
		Exec: impl.Load,
	},
	"import": types.BuiltinFunction{
		Args: "(id,unit)",
		Desc: "Import unit from the working environment with id",
		Exec: impl.ImportUnit,
	},
	"envs": types.BuiltinFunction{
		Args: "()",
		Desc: "List all available environments",
		Exec: impl.ListEnv,
	},
	"clear": types.BuiltinFunction{
		Args: "()",
		Desc: "Clear screen",
		Exec: impl.Clear,
	},
	"exit": types.BuiltinFunction{
		Args: "(code)",
		Desc: "Exit with error code",
		Exec: impl.Exit,
	},
}

func RegisterFunctions(ctx *types.Context) {
	for key, val := range functions {
		ctx.Builtin.RegisterFunction(key, val)
	}
}
