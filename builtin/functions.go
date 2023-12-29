package builtin

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"math/bits"
	"math/rand"
	"os"
	"sort"

	"github.com/dece2183/hexowl/input/syntax"
	"github.com/dece2183/hexowl/user"
	"github.com/dece2183/hexowl/utils"
)

type Func struct {
	// Arguments description.
	Args string
	// Function description.
	Desc string
	// Function that will be executed.
	Exec func(args ...interface{}) (interface{}, error)
}

type FuncMap map[string]Func

var loadedEnvDesc = ""

var functions = FuncMap{
	"sin": Func{
		Args: "(x)",
		Desc: "The sine of the radian argument x",
		Exec: sin,
	},
	"cos": Func{
		Args: "(x)",
		Desc: "The cosine of the radian argument x",
		Exec: cos,
	},
	"tan": Func{
		Args: "(x)",
		Desc: "The tangent of the radian argument x",
		Exec: tan,
	},
	"asin": Func{
		Args: "(x)",
		Desc: "The arcsine of the radian argument x",
		Exec: asin,
	},
	"acos": Func{
		Args: "(x)",
		Desc: "The arccosine of the radian argument x",
		Exec: acos,
	},
	"atan": Func{
		Args: "(x)",
		Desc: "The arctangent of the radian argument x",
		Exec: atan,
	},
	"pow": Func{
		Args: "(x,y)",
		Desc: "The base-x exponential of y",
		Exec: pow,
	},
	"sqrt": Func{
		Args: "(x)",
		Desc: "The square root of x",
		Exec: sqrt,
	},
	"exp": Func{
		Args: "(x)",
		Desc: "The base-e exponential of x",
		Exec: exp,
	},
	"logn": Func{
		Args: "(x)",
		Desc: "The natural logarithm of x",
		Exec: logn,
	},
	"log2": Func{
		Args: "(x)",
		Desc: "The binary logarithm of x",
		Exec: log2,
	},
	"log10": Func{
		Args: "(x)",
		Desc: "The decimal logarithm of x",
		Exec: log10,
	},
	"round": Func{
		Args: "(x)",
		Desc: "The nearest integer, rounding half away from zero",
		Exec: round,
	},
	"ceil": Func{
		Args: "(x)",
		Desc: "The least integer value greater than or equal to x",
		Exec: ceil,
	},
	"floor": Func{
		Args: "(x)",
		Desc: "The greatest integer value less than or equal to x",
		Exec: floor,
	},
	"rand": Func{
		Args: "(a,b)",
		Desc: "The random number in the range [a,b) or [0,1) if no arguments are passed",
		Exec: random,
	},
	"popcnt": {
		Args: "(x)",
		Desc: "The number of one bits (\"population count\") in x",
		Exec: popcount,
	},
	"vars": Func{
		Args: "()",
		Desc: "List available variables",
		Exec: vars,
	},
	"rmvar": Func{
		Args: "(name)",
		Desc: "Delete a specific user variable",
		Exec: rmVar,
	},
	"clvars": Func{
		Args: "()",
		Desc: "Delete user defined variables",
		Exec: clearVars,
	},
	"funcs": Func{
		Args: "()",
		Desc: "List alailable functions",
		Exec: funcs,
	},
	"rmfunc": Func{
		Args: "(name)",
		Desc: "Delete a specific user function",
		Exec: rmFunc,
	},
	"rmfuncvar": Func{
		Args: "(name,varid)",
		Desc: "Delete a specific user function variation",
		Exec: rmFuncVar,
	},
	"clfuncs": Func{
		Args: "()",
		Desc: "Delete user defined functions",
		Exec: clearFuncs,
	},
	"save": Func{
		Args: "(id,comment)",
		Desc: "Save working environment with id and optional comment",
		Exec: save,
	},
	"load": Func{
		Args: "(id)",
		Desc: "Load working environment with id",
		Exec: load,
	},
	"import": Func{
		Args: "(id,unit)",
		Desc: "Import unit from the working environment with id",
		Exec: importUnit,
	},
	"envs": Func{
		Args: "()",
		Desc: "List all available environments",
		Exec: listEnv,
	},
	"clear": Func{
		Args: "()",
		Desc: "Clear screen",
		Exec: clear,
	},
	"exit": Func{
		Args: "(code)",
		Desc: "Exit with error code",
		Exec: exit,
	},
}

func sin(args ...interface{}) (interface{}, error) {
	return math.Sin(utils.ToNumber[float64](args[0])), nil
}

func cos(args ...interface{}) (interface{}, error) {
	return math.Cos(utils.ToNumber[float64](args[0])), nil
}

func asin(args ...interface{}) (interface{}, error) {
	return math.Asin(utils.ToNumber[float64](args[0])), nil
}

func acos(args ...interface{}) (interface{}, error) {
	return math.Acos(utils.ToNumber[float64](args[0])), nil
}

func tan(args ...interface{}) (interface{}, error) {
	return math.Tan(utils.ToNumber[float64](args[0])), nil
}

func atan(args ...interface{}) (interface{}, error) {
	return math.Atan(utils.ToNumber[float64](args[0])), nil
}

func pow(args ...interface{}) (interface{}, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("not enough arguments")
	}
	return math.Pow(utils.ToNumber[float64](args[0]), utils.ToNumber[float64](args[1])), nil
}

func sqrt(args ...interface{}) (interface{}, error) {
	return math.Sqrt(utils.ToNumber[float64](args[0])), nil
}

func logn(args ...interface{}) (interface{}, error) {
	return math.Log(utils.ToNumber[float64](args[0])), nil
}

func log2(args ...interface{}) (interface{}, error) {
	return math.Log2(utils.ToNumber[float64](args[0])), nil
}

func log10(args ...interface{}) (interface{}, error) {
	return math.Log10(utils.ToNumber[float64](args[0])), nil
}

func exp(args ...interface{}) (interface{}, error) {
	return math.Exp(utils.ToNumber[float64](args[0])), nil
}

func round(args ...interface{}) (interface{}, error) {
	return math.Round(utils.ToNumber[float64](args[0])), nil
}

func ceil(args ...interface{}) (interface{}, error) {
	return math.Ceil(utils.ToNumber[float64](args[0])), nil
}

func floor(args ...interface{}) (interface{}, error) {
	return math.Floor(utils.ToNumber[float64](args[0])), nil
}

func random(args ...interface{}) (interface{}, error) {
	argslen := len(args)
	if argslen == 0 || args[0] == nil {
		return rand.Float64(), nil
	} else {
		if argslen == 1 {
			a := utils.ToNumber[int64](args[0])
			if a < 0 {
				return 0, fmt.Errorf("the first argument must be positive")
			}
			return rand.Int63n(a), nil
		} else {
			a := utils.ToNumber[int64](args[0])
			b := utils.ToNumber[int64](args[1])
			if b < a {
				return 0, fmt.Errorf("the first argument must be greater")
			}
			return rand.Int63n(b-a) + a, nil
		}
	}
}

func popcount(args ...interface{}) (interface{}, error) {
	return uint64(bits.OnesCount64(utils.ToNumber[uint64](args[0]))), nil
}

func vars(args ...interface{}) (interface{}, error) {
	userVars := user.ListVariables()
	varsCount := uint64(len(userVars))
	if varsCount > 0 {
		fmt.Fprintf(bDesc.system.Stdout, "\n\tUser variables:\n")
		keysList := make([]string, 0, len(userVars))
		for key := range userVars {
			keysList = append(keysList, key)
		}
		sort.Strings(keysList)
		for _, key := range keysList {
			var outstr string
			if str, isStr := userVars[key].(string); isStr {
				outstr = fmt.Sprintf("\t\t[%s] = \"%s\"\n", key, str)
			} else {
				outstr = fmt.Sprintf("\t\t[%s] = %v\n", key, userVars[key])
			}

			if bDesc.system.HighlightEnabled {
				outstr = syntax.Highlight(outstr)
			}

			fmt.Fprint(bDesc.system.Stdout, outstr)
		}
	} else {
		fmt.Fprintf(bDesc.system.Stdout, "\n\tThere are no user defined variables.\n")
	}
	if len(constants) > 0 {
		fmt.Fprintf(bDesc.system.Stdout, "\n\tBuiltin constants:\n")
		keysList := make([]string, 0, len(constants))
		for key := range constants {
			keysList = append(keysList, key)
		}
		sort.Strings(keysList)
		for _, key := range keysList {
			if key == "help" || key == "version" {
				continue
			}
			outstr := fmt.Sprintf("\t\t[%s] = %v\n", key, constants[key])

			if bDesc.system.HighlightEnabled {
				outstr = syntax.Highlight(outstr)
			}

			fmt.Fprint(bDesc.system.Stdout, outstr)
		}
	} else {
		fmt.Fprintf(bDesc.system.Stdout, "\n\tThere are no builtin constants.\n")
	}
	return varsCount, nil
}

func rmVar(args ...interface{}) (interface{}, error) {
	removedVars := 0

	for _, arg := range args {
		name, isString := arg.(string)
		if !isString || !user.HasVariable(name) {
			continue
		}
		user.DeleteVariable(name)
		removedVars++
	}

	return uint64(removedVars), nil
}

func clearVars(args ...interface{}) (interface{}, error) {
	user.DropVariables()
	return uint64(0), nil
}

func funcs(args ...interface{}) (interface{}, error) {
	userFuncs := user.ListFunctions()
	funcsCount := uint64(len(userFuncs))
	if funcsCount > 0 {
		fmt.Fprintf(bDesc.system.Stdout, "\n\tUser functions:\n")
		keysList := make([]string, 0, len(userFuncs))
		for key := range userFuncs {
			keysList = append(keysList, key)
		}
		sort.Strings(keysList)
		for _, key := range keysList {
			value := userFuncs[key]
			funcName := fmt.Sprintf("%-12s", key)

			if bDesc.system.HighlightEnabled {
				fmt.Fprintf(bDesc.system.Stdout, "\t\t%s%s\n", syntax.Colorize(funcName, utils.W_UNIT), syntax.Highlight(value.Variants[0].String()))
			} else {
				fmt.Fprintf(bDesc.system.Stdout, "\t\t%s%s\n", funcName, value.Variants[0].String())
			}

			for v := 1; v < len(value.Variants); v++ {
				if bDesc.system.HighlightEnabled {
					fmt.Fprintf(bDesc.system.Stdout, "\t\t%12s%s\n", "", syntax.Highlight(value.Variants[v].String()))
				} else {
					fmt.Fprintf(bDesc.system.Stdout, "\t\t%12s%s\n", "", value.Variants[v].String())
				}
			}
		}
	} else {
		fmt.Fprintf(bDesc.system.Stdout, "\n\tThere are no user defined functions.\n")
	}
	if len(*bDesc.functions) > 0 {
		fmt.Fprintf(bDesc.system.Stdout, "\n\tBuiltin functions:\n")
		keysList := make([]string, 0, len(*bDesc.functions))
		for key := range *bDesc.functions {
			keysList = append(keysList, key)
		}
		sort.Strings(keysList)
		for _, key := range keysList {
			value := (*bDesc.functions)[key]
			funcName := fmt.Sprintf("%-12s", key)
			funcArgs := fmt.Sprintf("%-12s", value.Args)

			if bDesc.system.HighlightEnabled {
				fmt.Fprintf(bDesc.system.Stdout, "\t\t%s%s - %s\n", syntax.Colorize(funcName, utils.W_FUNC), syntax.Highlight(funcArgs), value.Desc)
			} else {
				fmt.Fprintf(bDesc.system.Stdout, "\t\t%s%s - %s\n", funcName, funcArgs, value.Desc)
			}
		}
	} else {
		fmt.Fprintf(bDesc.system.Stdout, "\n\tThere are no builtin functions.\n")
	}
	return funcsCount, nil
}

func rmFunc(args ...interface{}) (interface{}, error) {
	removedFuncs := 0

	for _, arg := range args {
		name, isString := arg.(string)
		if !isString || !user.HasFunction(name) {
			continue
		}
		user.DeleteFunction(name)
		removedFuncs++
	}

	return uint64(removedFuncs), nil
}

func rmFuncVar(args ...interface{}) (interface{}, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("not enough arguments")
	}

	name, isString := args[0].(string)
	if !isString {
		return nil, fmt.Errorf("the function name must be a string")
	}
	if !user.HasFunction(name) {
		return nil, fmt.Errorf("the function '%s' does not exists", name)
	}

	varindex := utils.ToNumber[uint64](args[1])
	f, _ := user.GetFunction(name)
	if int(varindex) >= len(f.Variants) {
		return nil, fmt.Errorf("the function '%s' does not have variant %d", name, varindex)
	}

	user.DeleteFunctionVariant(name, int(varindex))

	return true, nil
}

func clearFuncs(args ...interface{}) (interface{}, error) {
	user.DropFunctions()
	return uint64(0), nil
}

func readEnvironment(envName string) (Environment, error) {
	if bDesc.system.ReadEnvironment == nil {
		return Environment{}, fmt.Errorf("not implemented")
	}

	f, err := bDesc.system.ReadEnvironment(envName)
	if err != nil {
		return Environment{}, err
	}
	defer f.Close()

	var env Environment
	dec := json.NewDecoder(f)
	err = dec.Decode(&env)
	if err != nil {
		return Environment{}, fmt.Errorf("unable to deserialize data")
	}

	return env, nil
}

func save(args ...interface{}) (interface{}, error) {
	if bDesc.system.WriteEnvironment == nil {
		return Environment{}, fmt.Errorf("not implemented")
	}

	var envName string

	// Get file name
	switch a := args[0].(type) {
	case string:
		envName = a
	default:
		envName = fmt.Sprintf("0x%016X", utils.ToNumber[uint64](args[0]))
	}

	// Save environment
	saveData := Environment{
		UserVars:  user.ListVariables(),
		UserFuncs: user.ListFunctions(),
	}
	// Add description if it is provided
	if len(args) > 1 {
		desc, success := args[1].(string)
		if success {
			saveData.Description = desc
		} else {
			return false, fmt.Errorf("second argument must be a string")
		}
	} else {
		saveData.Description = loadedEnvDesc
	}

	// Open file to write
	f, err := bDesc.system.WriteEnvironment(envName)
	if err != nil {
		return false, err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	err = enc.Encode(saveData)
	if err != nil {
		return false, fmt.Errorf("unable to serialize data")
	}

	fmt.Fprintf(bDesc.system.Stdout, "\n\tEnvironment saved as '%s'\n", envName)
	return true, nil
}

func load(args ...interface{}) (interface{}, error) {
	if bDesc.system.ReadEnvironment == nil {
		return Environment{}, fmt.Errorf("not implemented")
	}

	var envName string

	// Get file name
	switch a := args[0].(type) {
	case string:
		envName = a
	default:
		envName = fmt.Sprintf("0x%016X", utils.ToNumber[uint64](args[0]))
	}

	// Load environment
	loadData, err := readEnvironment(envName)
	if err != nil {
		return false, err
	}

	// Apply loaded environment
	loadedEnvDesc = loadData.Description
	user.DropVariables()
	for name, val := range loadData.UserVars {
		user.SetVariable(name, val)
	}
	user.DropFunctions()
	for name, val := range loadData.UserFuncs {
		user.SetFunction(name, val)
	}

	fmt.Fprintf(bDesc.system.Stdout, "\n\tEnvironment '%s' loaded\n", envName)
	return true, nil
}

func resolveDependencies(env Environment, unit string, fn user.Func) int {
	var loadedUnits = 0

	for _, variant := range fn.Variants {
		// Search for units and function calls that are not presented in current environment.
		// And try to find it in env Environment.
	nextWord:
		for _, word := range variant.Body {
			if word.Type != utils.W_UNIT && word.Type != utils.W_FUNC {
				continue
			}

			if user.HasVariable(word.Literal) || user.HasFunction(word.Literal) || HasConstant(word.Literal) || HasFunction(word.Literal) {
				continue
			}

			if word.Literal == unit {
				continue
			}

			// Search in user variables
			for key, val := range env.UserVars {
				if key == word.Literal {
					user.SetVariable(key, val)
					loadedUnits++
					continue nextWord
				}
			}

			// Search in user functions
			for key, val := range env.UserFuncs {
				if key == word.Literal {
					user.SetFunction(key, val)
					loadedUnits++
					continue nextWord
				}
			}
		}
	}

	return loadedUnits
}

func importUnit(args ...interface{}) (interface{}, error) {
	var envName string

	// Get file name
	switch a := args[0].(type) {
	case string:
		envName = a
	default:
		envName = fmt.Sprintf("0x%016X", utils.ToNumber[uint64](args[0]))
	}

	// Load environment
	loadedEnv, err := readEnvironment(envName)
	if err != nil {
		fmt.Fprintf(bDesc.system.Stdout, "\n\tEnvironment '%s' import failed: %s\n", envName, err)
		return false, nil
	}

	loadedUnits := 0

	if len(args) == 1 {
		// Import all environment
		for name, val := range loadedEnv.UserVars {
			user.SetVariable(name, val)
			loadedUnits++
		}
		for name, val := range loadedEnv.UserFuncs {
			user.SetFunction(name, val)
			loadedUnits++
		}
	} else {
		// Try to find units
		for i := 1; i < len(args); i++ {
			var name string
			switch a := args[i].(type) {
			case string:
				name = a
			default:
				continue
			}

			userVar, found := loadedEnv.UserVars[name]
			if found {
				user.SetVariable(name, userVar)
				loadedUnits++
			}

			userFunc, found := loadedEnv.UserFuncs[name]
			if found {
				loadedUnits += resolveDependencies(loadedEnv, name, userFunc)
				user.SetFunction(name, userFunc)
				loadedUnits++
			}
		}
	}

	fmt.Fprintf(bDesc.system.Stdout, "\n\tImported %d units from environment '%s'\n", loadedUnits, envName)
	return true, nil
}

func listEnv(args ...interface{}) (interface{}, error) {
	if bDesc.system.ListEnvironments == nil {
		return false, fmt.Errorf("not implemented")
	}

	envs, err := bDesc.system.ListEnvironments()
	if err != nil {
		return false, nil
	}

	if len(envs) == 0 {
		fmt.Fprintf(bDesc.system.Stdout, "\n\tThere are no saved environments\n")
		return uint64(0), nil
	}

	envCount := 0
	fmt.Fprintf(bDesc.system.Stdout, "\n\tAvailable environments:\n")
	for _, envName := range envs {
		fmt.Fprintf(bDesc.system.Stdout, "\t\t%s", envName)

		env, err := readEnvironment(envName)
		if err != nil {
			fmt.Fprintf(bDesc.system.Stdout, " - %s", err)
			continue
		}

		if env.Description != "" {
			fmt.Fprintf(bDesc.system.Stdout, " - %s", env.Description)
		}
		fmt.Fprintf(bDesc.system.Stdout, "\n")

		envCount++
	}

	return uint64(envCount), nil
}

func clear(args ...interface{}) (interface{}, error) {
	if bDesc.system.ClearScreen == nil {
		return nil, fmt.Errorf("not implemented")
	}

	bDesc.system.ClearScreen()
	return nil, nil
}

func exit(args ...interface{}) (interface{}, error) {
	exitCode := utils.ToNumber[int64](args[0])
	os.Exit(int(exitCode))
	return exitCode, nil
}

// Deprecated: Use builtin.SystemInit instead.
// builtin.SystemInit provides greater portability.
func FuncsInit(out io.Writer) {
	bDesc.system.Stdout = out
}

// Is function with name presented in the builtin function map.
func HasFunction(name string) bool {
	_, found := (*bDesc.functions)[name]
	return found
}

// Register a new function and add it to the builtin function map.
func RegisterFunction(name string, function Func) {
	(*bDesc.functions)[name] = function
}

// Get function by name from the builtin function map.
func GetFunction(name string) (function Func, found bool) {
	function, found = (*bDesc.functions)[name]
	return
}

// Return the builtin function map.
func ListFunctions() FuncMap {
	return (*bDesc.functions)
}
