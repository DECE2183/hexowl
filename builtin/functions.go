package builtin

import (
	"io"

	impl "github.com/dece2183/hexowl/builtin/function_impl"
	"github.com/dece2183/hexowl/builtin/types"
)

var functions = types.FunctionMap{
	"sin": types.Func{
		Args: "(x)",
		Desc: "The sine of the radian argument x",
		Exec: impl.Sin,
	},
	"cos": types.Func{
		Args: "(x)",
		Desc: "The cosine of the radian argument x",
		Exec: impl.Cos,
	},
	"tan": types.Func{
		Args: "(x)",
		Desc: "The tangent of the radian argument x",
		Exec: impl.Tan,
	},
	"asin": types.Func{
		Args: "(x)",
		Desc: "The arcsine of the radian argument x",
		Exec: impl.Asin,
	},
	"acos": types.Func{
		Args: "(x)",
		Desc: "The arccosine of the radian argument x",
		Exec: impl.Acos,
	},
	"atan": types.Func{
		Args: "(x)",
		Desc: "The arctangent of the radian argument x",
		Exec: impl.Atan,
	},
	"pow": types.Func{
		Args: "(x,y)",
		Desc: "The base-x exponential of y",
		Exec: impl.Pow,
	},
	"sqrt": types.Func{
		Args: "(x)",
		Desc: "The square root of x",
		Exec: impl.Sqrt,
	},
	"exp": types.Func{
		Args: "(x)",
		Desc: "The base-e exponential of x",
		Exec: impl.Exp,
	},
	"logn": types.Func{
		Args: "(x)",
		Desc: "The natural logarithm of x",
		Exec: impl.Logn,
	},
	"log2": types.Func{
		Args: "(x)",
		Desc: "The binary logarithm of x",
		Exec: impl.Log2,
	},
	"log10": types.Func{
		Args: "(x)",
		Desc: "The decimal logarithm of x",
		Exec: impl.Log10,
	},
	"round": types.Func{
		Args: "(x)",
		Desc: "The nearest integer, rounding half away from zero",
		Exec: impl.Round,
	},
	"ceil": types.Func{
		Args: "(x)",
		Desc: "The least integer value greater than or equal to x",
		Exec: impl.Ceil,
	},
	"floor": types.Func{
		Args: "(x)",
		Desc: "The greatest integer value less than or equal to x",
		Exec: impl.Floor,
	},
	"rand": types.Func{
		Args: "(a,b)",
		Desc: "The random number in the range [a,b) or [0,1) if no arguments are passed",
		Exec: impl.Random,
	},
	"popcnt": {
		Args: "(x)",
		Desc: "The number of one bits (\"population count\") in x",
		Exec: impl.Popcount,
	},
	"vars": types.Func{
		Args: "()",
		Desc: "List available variables",
		Exec: impl.Vars,
	},
	"rmvar": types.Func{
		Args: "(name)",
		Desc: "Delete a specific user variable",
		Exec: impl.RmVar,
	},
	"clvars": types.Func{
		Args: "()",
		Desc: "Delete user defined variables",
		Exec: impl.ClearVars,
	},
	"funcs": types.Func{
		Args: "()",
		Desc: "List alailable functions",
		Exec: impl.Funcs,
	},
	"rmfunc": types.Func{
		Args: "(name)",
		Desc: "Delete a specific user function",
		Exec: impl.RmFunc,
	},
	"rmfuncvar": types.Func{
		Args: "(name,varid)",
		Desc: "Delete a specific user function variation",
		Exec: impl.RmFuncVar,
	},
	"clfuncs": types.Func{
		Args: "()",
		Desc: "Delete user defined functions",
		Exec: impl.ClearFuncs,
	},
	"save": types.Func{
		Args: "(id,comment)",
		Desc: "Save working environment with id and optional comment",
		Exec: impl.Save,
	},
	"load": types.Func{
		Args: "(id)",
		Desc: "Load working environment with id",
		Exec: impl.Load,
	},
	"import": types.Func{
		Args: "(id,unit)",
		Desc: "Import unit from the working environment with id",
		Exec: impl.ImportUnit,
	},
	"envs": types.Func{
		Args: "()",
		Desc: "List all available environments",
		Exec: impl.ListEnv,
	},
	"clear": types.Func{
		Args: "()",
		Desc: "Clear screen",
		Exec: impl.Clear,
	},
	"exit": types.Func{
		Args: "(code)",
		Desc: "Exit with error code",
		Exec: impl.Exit,
	},
}

// Deprecated: Use builtin.SystemInit instead.
// builtin.SystemInit provides greater portability.
func FuncsInit(out io.Writer) {
	descriptor.System.Stdout = out
}

// Is function with name presented in the builtin function map.
func HasFunction(name string) bool {
	_, found := (descriptor.Functions)[name]
	return found
}

// Register a new function and add it to the builtin function map.
func RegisterFunction(name string, function types.Func) {
	(descriptor.Functions)[name] = function
}

// Get function by name from the builtin function map.
func GetFunction(name string) (function types.Func, found bool) {
	function, found = (descriptor.Functions)[name]
	return
}

// Return the builtin function map.
func ListFunctions() types.FunctionMap {
	return (descriptor.Functions)
}

// Execute builtin function
func Exec(function types.Func, args ...interface{}) (interface{}, error) {
	return function.Exec(&descriptor, args...)
}
