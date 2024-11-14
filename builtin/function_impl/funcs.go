package functionimpl

import (
	"fmt"
	"sort"

	"github.com/dece2183/hexowl/builtin/types"
	"github.com/dece2183/hexowl/input/syntax"
	"github.com/dece2183/hexowl/user"
	"github.com/dece2183/hexowl/utils"
)

func Funcs(desc *types.Descriptor, args ...interface{}) (interface{}, error) {
	userFuncs := user.ListFunctions()
	funcsCount := uint64(len(userFuncs))

	if funcsCount > 0 {
		fmt.Fprintf(desc.System.Stdout, "\n\tUser functions:\n")
		keysList := make([]string, 0, len(userFuncs))
		for key := range userFuncs {
			keysList = append(keysList, key)
		}
		sort.Strings(keysList)
		for _, key := range keysList {
			value := userFuncs[key]
			funcName := fmt.Sprintf("%-12s", key)

			if desc.System.HighlightEnabled {
				fmt.Fprintf(desc.System.Stdout, "\t\t%s%s\n", syntax.Colorize(funcName, utils.W_UNIT), syntax.Highlight(value.Variants[0].String()))
			} else {
				fmt.Fprintf(desc.System.Stdout, "\t\t%s%s\n", funcName, value.Variants[0].String())
			}

			for v := 1; v < len(value.Variants); v++ {
				if desc.System.HighlightEnabled {
					fmt.Fprintf(desc.System.Stdout, "\t\t%12s%s\n", "", syntax.Highlight(value.Variants[v].String()))
				} else {
					fmt.Fprintf(desc.System.Stdout, "\t\t%12s%s\n", "", value.Variants[v].String())
				}
			}
		}
	} else {
		fmt.Fprintf(desc.System.Stdout, "\n\tThere are no user defined functions.\n")
	}

	if len(desc.Functions) > 0 {
		fmt.Fprintf(desc.System.Stdout, "\n\tBuiltin functions:\n")
		keysList := make([]string, 0, len(desc.Functions))
		for key := range desc.Functions {
			keysList = append(keysList, key)
		}
		sort.Strings(keysList)
		for _, key := range keysList {
			value := (desc.Functions)[key]
			funcName := fmt.Sprintf("%-12s", key)
			funcArgs := fmt.Sprintf("%-12s", value.Args)

			if desc.System.HighlightEnabled {
				fmt.Fprintf(desc.System.Stdout, "\t\t%s%s - %s\n", syntax.Colorize(funcName, utils.W_FUNC), syntax.Highlight(funcArgs), value.Desc)
			} else {
				fmt.Fprintf(desc.System.Stdout, "\t\t%s%s - %s\n", funcName, funcArgs, value.Desc)
			}
		}
	} else {
		fmt.Fprintf(desc.System.Stdout, "\n\tThere are no builtin functions.\n")
	}

	return funcsCount, nil
}

func RmFunc(desc *types.Descriptor, args ...interface{}) (interface{}, error) {
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

func RmFuncVar(desc *types.Descriptor, args ...interface{}) (interface{}, error) {
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

func ClearFuncs(desc *types.Descriptor, args ...interface{}) (interface{}, error) {
	user.DropFunctions()
	return uint64(0), nil
}
