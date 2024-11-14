package functionimpl

import (
	"encoding/json"
	"fmt"

	"github.com/dece2183/hexowl/builtin/types"
	"github.com/dece2183/hexowl/user"
	"github.com/dece2183/hexowl/utils"
)

type environment struct {
	Description string
	UserVars    map[string]interface{}
	UserFuncs   map[string]user.Func
}

var loadedEnvDesc = ""

func readEnvironment(desc *types.Descriptor, envName string) (environment, error) {
	if desc.System.ReadEnvironment == nil {
		return environment{}, fmt.Errorf("'load' not implemented")
	}

	f, err := desc.System.ReadEnvironment(envName)
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

func resolveDependencies(desc *types.Descriptor, env environment, unit string, fn user.Func) int {
	var loadedUnits = 0

	for _, variant := range fn.Variants {
		// Search for units and function calls that are not presented in current environment.
		// And try to find it in env Environment.
	nextWord:
		for _, word := range variant.Body {
			if word.Type != utils.W_UNIT && word.Type != utils.W_FUNC {
				continue
			}

			_, hasConstant := desc.Constants[word.Literal]
			_, hasFunction := desc.Constants[word.Literal]

			if user.HasVariable(word.Literal) || user.HasFunction(word.Literal) || hasConstant || hasFunction {
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

func Save(desc *types.Descriptor, args ...interface{}) (interface{}, error) {
	if desc.System.WriteEnvironment == nil {
		return environment{}, fmt.Errorf("'save' not implemented")
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
	saveData := environment{
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
	f, err := desc.System.WriteEnvironment(envName)
	if err != nil {
		return false, err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	err = enc.Encode(saveData)
	if err != nil {
		return false, fmt.Errorf("unable to serialize data")
	}

	fmt.Fprintf(desc.System.Stdout, "\n\tEnvironment saved as '%s'\n", envName)
	return true, nil
}

func Load(desc *types.Descriptor, args ...interface{}) (interface{}, error) {
	if desc.System.ReadEnvironment == nil {
		return environment{}, fmt.Errorf("'load' not implemented")
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
	loadData, err := readEnvironment(desc, envName)
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

	fmt.Fprintf(desc.System.Stdout, "\n\tEnvironment '%s' loaded\n", envName)
	return true, nil
}

func ImportUnit(desc *types.Descriptor, args ...interface{}) (interface{}, error) {
	var envName string

	// Get file name
	switch a := args[0].(type) {
	case string:
		envName = a
	default:
		envName = fmt.Sprintf("0x%016X", utils.ToNumber[uint64](args[0]))
	}

	// Load environment
	loadedEnv, err := readEnvironment(desc, envName)
	if err != nil {
		fmt.Fprintf(desc.System.Stdout, "\n\tEnvironment '%s' import failed: %s\n", envName, err)
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
				loadedUnits += resolveDependencies(desc, loadedEnv, name, userFunc)
				user.SetFunction(name, userFunc)
				loadedUnits++
			}
		}
	}

	fmt.Fprintf(desc.System.Stdout, "\n\tImported %d units from environment '%s'\n", loadedUnits, envName)
	return true, nil
}

func ListEnv(desc *types.Descriptor, args ...interface{}) (interface{}, error) {
	if desc.System.ListEnvironments == nil {
		return false, fmt.Errorf("'envs' not implemented")
	}

	envs, err := desc.System.ListEnvironments()
	if err != nil {
		return false, nil
	}

	if len(envs) == 0 {
		fmt.Fprintf(desc.System.Stdout, "\n\tThere are no saved environments\n")
		return uint64(0), nil
	}

	envCount := 0
	fmt.Fprintf(desc.System.Stdout, "\n\tAvailable environments:\n")
	for _, envName := range envs {
		fmt.Fprintf(desc.System.Stdout, "\t\t%s", envName)

		env, err := readEnvironment(desc, envName)
		if err != nil {
			fmt.Fprintf(desc.System.Stdout, " - %s", err)
			continue
		}

		if env.Description != "" {
			fmt.Fprintf(desc.System.Stdout, " - %s", env.Description)
		}
		fmt.Fprintf(desc.System.Stdout, "\n")

		envCount++
	}

	return uint64(envCount), nil
}
