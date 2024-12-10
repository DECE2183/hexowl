package impl

import (
	"fmt"
	"sort"

	"github.com/dece2183/hexowl/v2/calculator/types"
	"github.com/dece2183/hexowl/v2/input/syntax"
	"github.com/dece2183/hexowl/v2/utils"
)

func Funcs(ctx *types.Context, args []interface{}) (interface{}, error) {
	userFuncs := ctx.User.ListFunctions()
	if len(userFuncs) > 0 {
		fmt.Fprintf(ctx.System.GetStdout(), "\n\tUser functions:\n")
		keysList := make([]string, 0, len(userFuncs))
		for key := range userFuncs {
			keysList = append(keysList, key)
		}
		sort.Strings(keysList)
		for _, key := range keysList {
			value := userFuncs[key]
			funcName := fmt.Sprintf("%-12s", key)

			if ctx.System.IsHighlightEnabled() {
				fmt.Fprintf(ctx.System.GetStdout(), "\t\t%s%s\n", syntax.Colorize(funcName, types.T_UNIT), syntax.Highlight(value.Variants[0].String()))
			} else {
				fmt.Fprintf(ctx.System.GetStdout(), "\t\t%s%s\n", funcName, value.Variants[0].String())
			}

			for v := 1; v < len(value.Variants); v++ {
				if ctx.System.IsHighlightEnabled() {
					fmt.Fprintf(ctx.System.GetStdout(), "\t\t%12s%s\n", "", syntax.Highlight(value.Variants[v].String()))
				} else {
					fmt.Fprintf(ctx.System.GetStdout(), "\t\t%12s%s\n", "", value.Variants[v].String())
				}
			}
		}
	} else {
		fmt.Fprintf(ctx.System.GetStdout(), "\n\tThere are no user defined functions.\n")
	}

	builtinFuncs := ctx.Builtin.ListFunctions()
	if len(builtinFuncs) > 0 {
		fmt.Fprintf(ctx.System.GetStdout(), "\n\tBuiltin functions:\n")
		keysList := make([]string, 0, len(builtinFuncs))
		for key := range builtinFuncs {
			keysList = append(keysList, key)
		}
		sort.Strings(keysList)
		for _, key := range keysList {
			value := builtinFuncs[key]
			funcName := fmt.Sprintf("%-12s", key)
			funcArgs := fmt.Sprintf("%-12s", value.Args)

			if ctx.System.IsHighlightEnabled() {
				fmt.Fprintf(ctx.System.GetStdout(), "\t\t%s%s - %s\n", syntax.Colorize(funcName, types.T_UNIT), syntax.Highlight(funcArgs), value.Desc)
			} else {
				fmt.Fprintf(ctx.System.GetStdout(), "\t\t%s%s - %s\n", funcName, funcArgs, value.Desc)
			}
		}
	} else {
		fmt.Fprintf(ctx.System.GetStdout(), "\n\tThere are no builtin functions.\n")
	}

	return uint64(len(userFuncs)), nil
}

func RmFunc(ctx *types.Context, args []interface{}) (interface{}, error) {
	removedFuncs := 0

	for _, arg := range args {
		name, isString := arg.(string)
		if !isString || !ctx.User.HasFunction(name) {
			continue
		}
		ctx.User.DeleteFunction(name)
		removedFuncs++
	}

	return uint64(removedFuncs), nil
}

func RmFuncVar(ctx *types.Context, args []interface{}) (interface{}, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("not enough arguments")
	}

	name, isString := args[0].(string)
	if !isString {
		return nil, fmt.Errorf("the function name must be a string")
	}
	if !ctx.User.HasFunction(name) {
		return nil, fmt.Errorf("the function '%s' does not exists", name)
	}

	varindex := utils.ToNumber[uint64](args[1])
	f, _ := ctx.User.GetFunction(name)
	if int(varindex) >= len(f.Variants) {
		return nil, fmt.Errorf("the function '%s' does not have variant %d", name, varindex)
	}

	ctx.User.DeleteFunctionVariant(name, int(varindex))

	return true, nil
}

func ClearFuncs(ctx *types.Context, args []interface{}) (interface{}, error) {
	ctx.User.DeleteAllFunctions()
	return uint64(0), nil
}
