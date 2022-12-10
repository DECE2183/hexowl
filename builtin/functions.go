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
	"time"

	"github.com/dece2183/hexowl/user"
	"github.com/dece2183/hexowl/utils"
)

type saveStruct struct {
	Description string
	UserVars    map[string]interface{}
	UserFuncs   map[string]user.Func
}

type Func struct {
	Args string
	Desc string
	Exec func(args ...interface{}) (interface{}, error)
}
type FuncMap map[string]Func

var loadedEnvDesc = ""
var outStream io.Writer

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
		Args: "(id)",
		Desc: "Save working environment with id",
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

var bFuncs *FuncMap

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
		fmt.Fprintf(outStream, "\n\tUser variables:\n")
		keysList := make([]string, 0, len(userVars))
		for key := range userVars {
			keysList = append(keysList, key)
		}
		sort.Strings(keysList)
		for _, key := range keysList {
			fmt.Fprintf(outStream, "\t\t[%s] = %v\n", key, userVars[key])
		}
	} else {
		fmt.Fprintf(outStream, "\n\tThere are no user defined variables.\n")
	}
	if len(constants) > 0 {
		fmt.Fprintf(outStream, "\n\tBuiltin constants:\n")
		keysList := make([]string, 0, len(constants))
		for key := range constants {
			keysList = append(keysList, key)
		}
		sort.Strings(keysList)
		for _, key := range keysList {
			if key == "help" || key == "version" {
				continue
			}
			fmt.Fprintf(outStream, "\t\t[%s] = %v\n", key, constants[key])
		}
	} else {
		fmt.Fprintf(outStream, "\n\tThere are no builtin constants.\n")
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
		fmt.Fprintf(outStream, "\n\tUser functions:\n")
		keysList := make([]string, 0, len(userFuncs))
		for key := range userFuncs {
			keysList = append(keysList, key)
		}
		sort.Strings(keysList)
		for _, key := range keysList {
			value := userFuncs[key]
			fmt.Fprintf(outStream, "\t\t%-12s%s\n", key, value.Variants[0])
			for v := 1; v < len(value.Variants); v++ {
				fmt.Fprintf(outStream, "\t\t%12s%s\n", "", value.Variants[v])
			}
		}
	} else {
		fmt.Fprintf(outStream, "\n\tThere are no user defined functions.\n")
	}
	if len(*bFuncs) > 0 {
		fmt.Fprintf(outStream, "\n\tBuiltin functions:\n")
		keysList := make([]string, 0, len(*bFuncs))
		for key := range *bFuncs {
			keysList = append(keysList, key)
		}
		sort.Strings(keysList)
		for _, key := range keysList {
			value := (*bFuncs)[key]
			fmt.Fprintf(outStream, "\t\t%-12s%-12s - %s\n", key, value.Args, value.Desc)
		}
	} else {
		fmt.Fprintf(outStream, "\n\tThere are no builtin functions.\n")
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

func getEnvPath(envName string) (string, error) {
	userDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	saveDir := fmt.Sprintf("%s/.hexowl/environment", userDir)
	err = os.MkdirAll(saveDir, 0666)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/%s.json", saveDir, envName), nil
}

func getEnvContent(envName string) (saveStruct, error) {
	path, err := getEnvPath(envName)
	if err != nil {
		return saveStruct{}, fmt.Errorf("unable to get access to file")
	}
	data := saveStruct{}
	buffer, err := os.ReadFile(path)
	if err != nil {
		return saveStruct{}, fmt.Errorf("environment doesn't exists")
	}
	err = json.Unmarshal(buffer, &data)
	if err != nil {
		return saveStruct{}, fmt.Errorf("unable to parse environment data")
	}
	return data, nil
}

func save(args ...interface{}) (interface{}, error) {
	var envName string

	// Get file name
	switch a := args[0].(type) {
	case string:
		envName = a
	default:
		envName = fmt.Sprintf("0x%016X", utils.ToNumber[uint64](args[0]))
	}

	// Save environment
	savePath, err := getEnvPath(envName)
	if err != nil {
		fmt.Fprintf(outStream, "\n\tEnvironment '%s' save failed: unable to get access to file\n", envName)
		return false, nil
	}
	saveData := saveStruct{
		UserVars:  user.ListVariables(),
		UserFuncs: user.ListFunctions(),
	}
	// Add description if it is provided
	if len(args) > 1 {
		desc, success := args[1].(string)
		if success {
			saveData.Description = desc
		} else {
			fmt.Fprintf(outStream, "\n\tEnvironment '%s' save failed: second argument must be a string\n", envName)
			return false, nil
		}
	} else {
		saveData.Description = loadedEnvDesc
	}
	saveJson, err := json.Marshal(saveData)
	if err != nil {
		fmt.Fprintf(outStream, "\n\tEnvironment '%s' save failed: unable to create data\n", envName)
		return false, nil
	}
	err = os.WriteFile(savePath, saveJson, 0666)
	if err != nil {
		fmt.Fprintf(outStream, "\n\tEnvironment '%s' save failed: unable to write file\n", envName)
		return false, nil
	}

	fmt.Fprintf(outStream, "\n\tSaving environment as '%s'\n", envName)
	return true, nil
}

func load(args ...interface{}) (interface{}, error) {
	var envName string

	// Get file name
	switch a := args[0].(type) {
	case string:
		envName = a
	default:
		envName = fmt.Sprintf("0x%016X", utils.ToNumber[uint64](args[0]))
	}

	// Load environment
	loadData, err := getEnvContent(envName)
	if err != nil {
		fmt.Fprintf(outStream, "\n\tEnvironment '%s' load failed: %s\n", envName, err)
		return false, nil
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

	fmt.Fprintf(outStream, "\n\tEnvironment '%s' loaded\n", envName)
	return true, nil
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
	loadData, err := getEnvContent(envName)
	if err != nil {
		fmt.Fprintf(outStream, "\n\tEnvironment '%s' import failed: %s\n", envName, err)
		return false, nil
	}

	loadedUnits := 0

	if len(args) == 1 {
		// Import all environment
		for name, val := range loadData.UserVars {
			user.SetVariable(name, val)
			loadedUnits++
		}
		for name, val := range loadData.UserFuncs {
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

			userVar, found := loadData.UserVars[name]
			if found {
				user.SetVariable(name, userVar)
				loadedUnits++
			}

			userFunc, found := loadData.UserFuncs[name]
			if found {
				user.SetFunction(name, userFunc)
				loadedUnits++
			}
		}
	}

	fmt.Fprintf(outStream, "\n\tImported %d units from environment '%s'\n", loadedUnits, envName)
	return true, nil
}

func listEnv(args ...interface{}) (interface{}, error) {
	userDir, err := os.UserHomeDir()
	if err != nil {
		return false, nil
	}
	envDir := fmt.Sprintf("%s/.hexowl/environment", userDir)
	err = os.MkdirAll(envDir, 0666)
	if err != nil {
		return false, err
	}

	f, err := os.Open(envDir)
	if err != nil {
		fmt.Fprintf(outStream, "\n\tEnvironment list failed: unable to get access to the directory with environments\n")
		return false, nil
	}
	files, err := f.Readdir(0)
	if err != nil {
		fmt.Fprintf(outStream, "\n\tEnvironment list failed: unable to read the directory with environments\n")
		return false, nil
	}

	if len(files) > 0 {
		envCount := 0
		fmt.Fprintf(outStream, "\n\tAvailable environments:\n")
		for _, v := range files {
			fname := v.Name()
			extIndex := len(fname) - 5
			if extIndex < 0 {
				continue
			}

			fname = fname[:extIndex]

			env, err := getEnvContent(fname)
			if err != nil {
				continue
			}

			fmt.Fprintf(outStream, "\t\t%s", fname)
			if env.Description != "" {
				fmt.Fprintf(outStream, " - %s", env.Description)
			}
			fmt.Fprintf(outStream, "\n")

			envCount++
		}
		return uint64(envCount), nil
	} else {
		fmt.Fprintf(outStream, "\n\tThere are no saved environments\n")
		return uint64(0), nil
	}
}

func clear(args ...interface{}) (interface{}, error) {
	fmt.Fprintf(outStream, "\x1bc")
	return nil, nil
}

func exit(args ...interface{}) (interface{}, error) {
	exitCode := utils.ToNumber[int64](args[0])
	os.Exit(int(exitCode))
	return exitCode, nil
}

func FuncsInit(out io.Writer) {
	bFuncs = &functions
	outStream = out
	rand.Seed(time.Now().UnixNano())
}

func HasFunction(name string) bool {
	_, found := functions[name]
	return found
}

func RegisterFunction(name string, function Func) {
	functions[name] = function
}

func GetFunction(name string) (function Func, found bool) {
	function, found = functions[name]
	return
}

func ListFunctions() map[string]Func {
	return functions
}
