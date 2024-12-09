package impl

import (
	"fmt"
	"sort"

	"github.com/dece2183/hexowl/calculator/types"
	"github.com/dece2183/hexowl/input/syntax"
)

func Vars(ctx *types.Context, args []interface{}) (interface{}, error) {
	userVars := ctx.User.ListVariables()
	if len(userVars) > 0 {
		fmt.Fprintf(ctx.System.GetStdout(), "\n\tUser variables:\n")
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

			if ctx.System.IsHighlightEnabled() {
				outstr = syntax.Highlight(outstr)
			}

			fmt.Fprint(ctx.System.GetStdout(), outstr)
		}
	} else {
		fmt.Fprintf(ctx.System.GetStdout(), "\n\tThere are no user defined variables.\n")
	}

	builtinConstants := ctx.Builtin.ListConstants()
	if len(builtinConstants) > 0 {
		fmt.Fprintf(ctx.System.GetStdout(), "\n\tBuiltin constants:\n")
		keysList := make([]string, 0, len(builtinConstants))
		for key := range builtinConstants {
			keysList = append(keysList, key)
		}
		sort.Strings(keysList)
		for _, key := range keysList {
			if key == "help" || key == "version" {
				continue
			}
			outstr := fmt.Sprintf("\t\t[%s] = %v\n", key, builtinConstants[key])

			if ctx.System.IsHighlightEnabled() {
				outstr = syntax.Highlight(outstr)
			}

			fmt.Fprint(ctx.System.GetStdout(), outstr)
		}
	} else {
		fmt.Fprintf(ctx.System.GetStdout(), "\n\tThere are no builtin constants.\n")
	}

	return uint64(len(userVars)), nil
}

func RmVar(ctx *types.Context, args []interface{}) (interface{}, error) {
	removedVars := 0

	for _, arg := range args {
		name, isString := arg.(string)
		if !isString || !ctx.User.HasVariable(name) {
			continue
		}
		ctx.User.DeleteVariable(name)
		removedVars++
	}

	return uint64(removedVars), nil
}

func ClearVars(ctx *types.Context, args []interface{}) (interface{}, error) {
	ctx.User.DeleteALlVariables()
	return uint64(0), nil
}
