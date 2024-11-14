package functionimpl

import (
	"fmt"
	"sort"

	"github.com/dece2183/hexowl/builtin/types"
	"github.com/dece2183/hexowl/input/syntax"
	"github.com/dece2183/hexowl/user"
)

func Vars(desc *types.Descriptor, args ...interface{}) (interface{}, error) {
	userVars := user.ListVariables()
	varsCount := uint64(len(userVars))

	if varsCount > 0 {
		fmt.Fprintf(desc.System.Stdout, "\n\tUser variables:\n")
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

			if desc.System.HighlightEnabled {
				outstr = syntax.Highlight(outstr)
			}

			fmt.Fprint(desc.System.Stdout, outstr)
		}
	} else {
		fmt.Fprintf(desc.System.Stdout, "\n\tThere are no user defined variables.\n")
	}

	if len(desc.Constants) > 0 {
		fmt.Fprintf(desc.System.Stdout, "\n\tBuiltin constants:\n")
		keysList := make([]string, 0, len(desc.Constants))
		for key := range desc.Constants {
			keysList = append(keysList, key)
		}
		sort.Strings(keysList)
		for _, key := range keysList {
			if key == "help" || key == "version" {
				continue
			}
			outstr := fmt.Sprintf("\t\t[%s] = %v\n", key, desc.Constants[key])

			if desc.System.HighlightEnabled {
				outstr = syntax.Highlight(outstr)
			}

			fmt.Fprint(desc.System.Stdout, outstr)
		}
	} else {
		fmt.Fprintf(desc.System.Stdout, "\n\tThere are no builtin constants.\n")
	}

	return varsCount, nil
}

func RmVar(desc *types.Descriptor, args ...interface{}) (interface{}, error) {
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

func ClearVars(desc *types.Descriptor, args ...interface{}) (interface{}, error) {
	user.DropVariables()
	return uint64(0), nil
}
