package impl

import (
	"encoding/json"
	"fmt"

	"github.com/dece2183/hexowl/calculator/types"
	"github.com/dece2183/hexowl/utils"
)

type environment struct {
	Description string
	UserVars    types.UserVariableMap
	UserFuncs   types.UserFunctionMap
}

const (
	_ENV_NAME_VAR        = "name"
	_ENV_DESCRIPTION_VAR = "description"
)

func readEnvironment(ctx *types.Context, envName string) (environment, error) {
	f, err := ctx.System.ReadEnvironment(envName)
	if err != nil {
		return environment{}, err
	}
	defer f.Close()

	var env environment
	dec := json.NewDecoder(f)
	err = dec.Decode(&env)
	if err != nil {
		return environment{}, fmt.Errorf("unable to deserialize data")
	}

	return env, nil
}

func resolveDependencies(ctx *types.Context, env environment, unit string, fn types.UserFunction) int {
	var loadedUnits = 0

	for _, variant := range fn.Variants {
		_ = variant
		//TODO:implement

		// Search for units and function calls that are not presented in current environment.
		// And try to find it in env Environment.

		// nextWord:
		// for _, word := range variant.Body {
		// 	if word.Type != utils.W_UNIT && word.Type != utils.W_FUNC {
		// 		continue
		// 	}

		// 	_, hasConstant := ctx.Builtin.GetConstant(word.Literal)
		// 	_, hasFunction := ctx.Builtin.GetConstant(word.Literal)

		// 	if ctx.User.HasVariable(word.Literal) || ctx.User.HasFunction(word.Literal) || hasConstant || hasFunction {
		// 		continue
		// 	}

		// 	if word.Literal == unit {
		// 		continue
		// 	}

		// 	// Search in user variables
		// 	for key, val := range env.UserVars {
		// 		if key == word.Literal {
		// 			ctx.User.SetVariable(key, val)
		// 			loadedUnits++
		// 			continue nextWord
		// 		}
		// 	}

		// 	// Search in user functions
		// 	for key, val := range env.UserFuncs {
		// 		if key == word.Literal {
		// 			ctx.User.SetFunction(key, val)
		// 			loadedUnits++
		// 			continue nextWord
		// 		}
		// 	}
		// }
	}

	return loadedUnits
}

func Save(ctx *types.Context, args []interface{}) (interface{}, error) {
	var envName string

	// Get file name
	switch a := args[0].(type) {
	case string:
		ctx.User.SetVariable(_ENV_NAME_VAR, a)
		envName = a
	default:
		envNum := utils.ToNumber[uint64](args[0])
		if envNum == 0 && ctx.User.HasVariable(_ENV_NAME_VAR) {
			userName, _ := ctx.User.GetVariable(_ENV_NAME_VAR)
			switch name := userName.(type) {
			case string:
				envName = name
			default:
				envName = fmt.Sprintf("0x%016X", utils.ToNumber[uint64](name))
			}
		} else {
			ctx.User.SetVariable(_ENV_NAME_VAR, envNum)
			envName = fmt.Sprintf("0x%016X", envNum)
		}
	}

	var envDescription string

	// Get env description
	if len(args) > 1 {
		desc, success := args[1].(string)
		if success {
			ctx.User.SetVariable(_ENV_DESCRIPTION_VAR, desc)
			envDescription = desc
		} else {
			return false, fmt.Errorf("second argument must be a string")
		}
	} else if ctx.User.HasVariable(_ENV_DESCRIPTION_VAR) {
		desc, _ := ctx.User.GetVariable(_ENV_DESCRIPTION_VAR)
		envDescription, _ = desc.(string)
	}

	// Save environment
	saveData := environment{
		UserVars:    ctx.User.ListVariables(),
		UserFuncs:   ctx.User.ListFunctions(),
		Description: envDescription,
	}

	// Open file to write
	f, err := ctx.System.WriteEnvironment(envName)
	if err != nil {
		return false, err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	err = enc.Encode(saveData)
	if err != nil {
		return false, fmt.Errorf("unable to serialize data")
	}

	fmt.Fprintf(ctx.System.GetStdout(), "\n\tEnvironment saved as '%s'\n", envName)
	return true, nil
}

func Load(ctx *types.Context, args []interface{}) (interface{}, error) {
	var envName string

	// Get file name
	switch a := args[0].(type) {
	case string:
		envName = a
	default:
		envName = fmt.Sprintf("0x%016X", utils.ToNumber[uint64](args[0]))
	}

	// Load environment
	loadData, err := readEnvironment(ctx, envName)
	if err != nil {
		return false, err
	}

	// Apply loaded environment
	ctx.User.DeleteALlVariables()
	for name, val := range loadData.UserVars {
		ctx.User.SetVariable(name, val)
	}
	ctx.User.DeleteAllFunctions()
	for name, val := range loadData.UserFuncs {
		ctx.User.SetFunction(name, val)
	}

	fmt.Fprintf(ctx.System.GetStdout(), "\n\tEnvironment '%s' loaded\n", envName)
	return true, nil
}

func ImportUnit(ctx *types.Context, args []interface{}) (interface{}, error) {
	var envName string

	// Get file name
	switch a := args[0].(type) {
	case string:
		envName = a
	default:
		envName = fmt.Sprintf("0x%016X", utils.ToNumber[uint64](args[0]))
	}

	// Load environment
	loadedEnv, err := readEnvironment(ctx, envName)
	if err != nil {
		fmt.Fprintf(ctx.System.GetStdout(), "\n\tEnvironment '%s' import failed: %s\n", envName, err)
		return false, nil
	}

	loadedUnits := 0

	if len(args) == 1 {
		// Import all environment
		for name, val := range loadedEnv.UserVars {
			ctx.User.SetVariable(name, val)
			loadedUnits++
		}
		for name, val := range loadedEnv.UserFuncs {
			ctx.User.SetFunction(name, val)
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
				ctx.User.SetVariable(name, userVar)
				loadedUnits++
			}

			userFunc, found := loadedEnv.UserFuncs[name]
			if found {
				loadedUnits += resolveDependencies(ctx, loadedEnv, name, userFunc)
				ctx.User.SetFunction(name, userFunc)
				loadedUnits++
			}
		}
	}

	fmt.Fprintf(ctx.System.GetStdout(), "\n\tImported %d units from environment '%s'\n", loadedUnits, envName)
	return true, nil
}

func ListEnv(ctx *types.Context, args []interface{}) (interface{}, error) {
	envs, err := ctx.System.ListEnvironments()
	if err != nil {
		return false, nil
	}

	if len(envs) == 0 {
		fmt.Fprintf(ctx.System.GetStdout(), "\n\tThere are no saved environments\n")
		return uint64(0), nil
	}

	envCount := 0
	fmt.Fprintf(ctx.System.GetStdout(), "\n\tAvailable environments:\n")
	for _, envName := range envs {
		fmt.Fprintf(ctx.System.GetStdout(), "\t\t%s", envName)

		env, err := readEnvironment(ctx, envName)
		if err != nil {
			fmt.Fprintf(ctx.System.GetStdout(), " - %s", err)
			continue
		}

		if env.Description != "" {
			fmt.Fprintf(ctx.System.GetStdout(), " - %s", env.Description)
		}
		fmt.Fprintf(ctx.System.GetStdout(), "\n")

		envCount++
	}

	return uint64(envCount), nil
}
