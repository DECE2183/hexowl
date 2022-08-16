package builtin

import (
	"encoding/json"
	"fmt"
	"math"
	"math/bits"
	"os"

	"github.com/dece2183/hexowl/user"
	"github.com/dece2183/hexowl/utils"
)

type saveStruct struct {
	UserVars  map[string]interface{}
	UserFuncs map[string]user.Func
}

type Func struct {
	Args string
	Desc string
	Exec func(args ...interface{}) (interface{}, error)
}
type FuncMap map[string]Func

var functions = FuncMap{
	"sin": Func{
		Args: "(x)",
		Desc: "Sine of the radian argument x",
		Exec: sin,
	},
	"cos": Func{
		Args: "(x)",
		Desc: "Cosine of the radian argument x",
		Exec: cos,
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
	if len(args) < 1 {
		return nil, fmt.Errorf("not enough arguments")
	}
	return math.Sin(utils.ToNumber[float64](args[0])), nil
}

func cos(args ...interface{}) (interface{}, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("not enough arguments")
	}
	return math.Cos(utils.ToNumber[float64](args[0])), nil
}

func pow(args ...interface{}) (interface{}, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("not enough arguments")
	}
	return math.Pow(utils.ToNumber[float64](args[0]), utils.ToNumber[float64](args[1])), nil
}

func sqrt(args ...interface{}) (interface{}, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("not enough arguments")
	}
	return math.Sqrt(utils.ToNumber[float64](args[0])), nil
}

func exp(args ...interface{}) (interface{}, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("not enough arguments")
	}
	return math.Exp(utils.ToNumber[float64](args[0])), nil
}

func round(args ...interface{}) (interface{}, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("not enough arguments")
	}
	return math.Round(utils.ToNumber[float64](args[0])), nil
}

func ceil(args ...interface{}) (interface{}, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("not enough arguments")
	}
	return math.Ceil(utils.ToNumber[float64](args[0])), nil
}

func floor(args ...interface{}) (interface{}, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("not enough arguments")
	}
	return math.Floor(utils.ToNumber[float64](args[0])), nil
}

func popcount(args ...interface{}) (interface{}, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("not enough arguments")
	}
	return uint64(bits.OnesCount64(utils.ToNumber[uint64](args[0]))), nil
}

func vars(args ...interface{}) (interface{}, error) {
	userVars := user.ListVariables()
	varsCount := uint64(len(userVars))
	if varsCount > 0 {
		fmt.Printf("\n\tUser variables:\n")
		for key, value := range userVars {
			fmt.Printf("\t\t[%s] = %v\n", key, value)
		}
	} else {
		fmt.Printf("\n\tThere are no user defined variables.\n")
	}
	if len(constants) > 0 {
		fmt.Printf("\n\tBuiltin constants:\n")
		for key, value := range constants {
			if key == "help" {
				continue
			}
			fmt.Printf("\t\t[%s] = %v\n", key, value)
		}
	} else {
		fmt.Printf("\n\tThere are no builtin constants.\n")
	}
	return varsCount, nil
}

func clearVars(args ...interface{}) (interface{}, error) {
	user.DropVariables()
	return uint64(0), nil
}

func funcs(args ...interface{}) (interface{}, error) {
	userFuncs := user.ListFunctions()
	funcsCount := uint64(len(userFuncs))
	if funcsCount > 0 {
		fmt.Printf("\n\tUser functions:\n")
		for key, value := range userFuncs {
			fmt.Printf("\t\t%-8s%s\n", key, value.Variants[0])
			for v := 1; v < len(value.Variants); v++ {
				fmt.Printf("\t\t\t%s\n", value.Variants[v])
			}
		}
	} else {
		fmt.Printf("\n\tThere are no user defined functions.\n")
	}
	if len(*bFuncs) > 0 {
		fmt.Printf("\n\tBuiltin functions:\n")
		for key, value := range *bFuncs {
			fmt.Printf("\t\t%-8s%-8s - %s\n", key, value.Args, value.Desc)
		}
	} else {
		fmt.Printf("\n\tThere are no builtin functions.\n")
	}
	return funcsCount, nil
}

func save(args ...interface{}) (interface{}, error) {
	envID := utils.ToNumber[uint64](args[0])
	userDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("\n\tEnvironment 0x%016X save failed: unable to get user home directory\n", envID)
		return 0, nil
	}
	saveDir := fmt.Sprintf("%s/.hexowl/environment", userDir)
	err = os.MkdirAll(saveDir, 0666)
	if err != nil {
		fmt.Printf("\n\tEnvironment 0x%016X save failed: unable to create save directory\n", envID)
		return 0, nil
	}
	savePath := fmt.Sprintf("%s/0x%016X.json", saveDir, envID)
	saveData := saveStruct{
		UserVars:  user.ListVariables(),
		UserFuncs: user.ListFunctions(),
	}
	saveJson, err := json.Marshal(saveData)
	if err != nil {
		fmt.Printf("\n\tEnvironment 0x%016X save failed: unable to create data\n", envID)
		return 0, nil
	}
	err = os.WriteFile(savePath, saveJson, 0666)
	if err != nil {
		fmt.Printf("\n\tEnvironment 0x%016X save failed: unable to write file\n", envID)
		return 0, nil
	}
	fmt.Printf("\n\tSaving environment as 0x%016X\n", envID)
	return uint64(1), nil
}

func load(args ...interface{}) (interface{}, error) {
	envID := utils.ToNumber[uint64](args[0])
	userDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("\n\tEnvironment 0x%016X load failed: unable to get user home directory\n", envID)
		return 0, nil
	}
	loadPath := fmt.Sprintf("%s/.hexowl/environment/0x%016X.json", userDir, envID)
	loadData := saveStruct{}
	loadBuffer, err := os.ReadFile(loadPath)
	if err != nil {
		fmt.Printf("\n\tEnvironment 0x%016X load failed: environment doesn't exists\n", envID)
		return 0, nil
	}
	err = json.Unmarshal(loadBuffer, &loadData)
	if err != nil {
		fmt.Printf("\n\tEnvironment 0x%016X load failed: unable to parse environment data\n", envID)
		return 0, nil
	}
	user.DropVariables()
	for name, val := range loadData.UserVars {
		user.SetVariable(name, val)
	}
	user.DropFunctions()
	for name, val := range loadData.UserFuncs {
		user.SetFunction(name, val)
	}
	fmt.Printf("\n\tEnvironment 0x%016X loaded\n", envID)
	return uint64(1), nil
}

func clear(args ...interface{}) (interface{}, error) {
	fmt.Printf("\x1bc")
	return nil, nil
}

func exit(args ...interface{}) (interface{}, error) {
	exitCode := utils.ToNumber[int64](args[0])
	os.Exit(int(exitCode))
	return exitCode, nil
}

func FuncsInit() {
	bFuncs = &functions
}

func HasFunction(name string) bool {
	_, found := functions[name]
	return found
}

func GetFunction(name string) (function Func, found bool) {
	function, found = functions[name]
	return
}

func ListFunctions() map[string]Func {
	return functions
}
