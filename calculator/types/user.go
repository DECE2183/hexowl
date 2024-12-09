package types

import "slices"

type UserFunctionVariant struct {
}

// Return argument names of function variant as strings slice.
func (v UserFunctionVariant) ArgNames() []string {
	//TODO: implement
	return []string{}
}

// fmt.Stringer interface implementation.
func (v UserFunctionVariant) String() string {
	//TODO: implement
	return ""
}

type UserFunction struct {
	Variants []UserFunctionVariant
}

type UserFunctionMap map[string]UserFunction

type UserVariableMap map[string]interface{}

type UserContainer struct {
	functions UserFunctionMap
	variables UserVariableMap
}

// Is function presented in the user functions map.
func (u *UserContainer) HasFunction(name string) bool {
	_, found := u.functions[name]
	return found
}

// Set user function with given name.
func (u *UserContainer) SetFunction(name string, function UserFunction) {
	u.functions[name] = function
}

// Set function varian for the funtion with given name.
func (u *UserContainer) SetFunctionVariant(name string, variant UserFunctionVariant) {
	var currentFunc UserFunction

	if !u.HasFunction(name) {
		currentFunc = UserFunction{
			Variants: make([]UserFunctionVariant, 0),
		}
	} else {
		currentFunc = u.functions[name]
	}

	for i, v := range currentFunc.Variants {
		// if variant with such arguments already exists replace it
		// TODO: implement
		_, _ = i, v
	}

	currentFunc.Variants = append(currentFunc.Variants, variant)
	u.functions[name] = currentFunc
}

// Get user function by name from the user functions map.
func (u *UserContainer) GetFunction(name string) (function UserFunction, found bool) {
	function, found = u.functions[name]
	return
}

// Return the user function map.
func (u *UserContainer) ListFunctions() UserFunctionMap {
	return u.functions
}

// Delete user function with name.
func (u *UserContainer) DeleteFunction(name string) {
	delete(u.functions, name)
}

// Delete user function variant by id.
func (u *UserContainer) DeleteFunctionVariant(name string, idx int) {
	f := u.functions[name]
	f.Variants = slices.Delete(f.Variants, idx, idx+1)
	u.functions[name] = f
}

// Delete all user defined functions.
func (u *UserContainer) DeleteAllFunctions() {
	u.functions = make(UserFunctionMap)
}

// Is variable with name presented in the user variables map.
func (u *UserContainer) HasVariable(name string) bool {
	_, found := u.variables[name]
	return found
}

// Set user variable with name and value.
func (u *UserContainer) SetVariable(name string, val interface{}) {
	u.variables[name] = val
}

// Get user variable with given name.
func (u *UserContainer) GetVariable(name string) (val interface{}, found bool) {
	val, found = u.variables[name]
	return
}

// Return the user variables map.
func (u *UserContainer) ListVariables() UserVariableMap {
	return u.variables
}

// Delete user variable with given name.
func (u *UserContainer) DeleteVariable(name string) {
	delete(u.variables, name)
}

// Delete all user variables.
func (u *UserContainer) DeleteALlVariables() {
	u.variables = make(UserVariableMap)
}

func (u *UserContainer) Predict(word string) string {
	for v := range u.variables {
		if len(v) < len(word) {
			continue
		}
		if v[:len(word)] == word {
			return v
		}
	}

	for f := range u.functions {
		if len(f) < len(word) {
			continue
		}
		if f[:len(word)] == word {
			return f + "()"
		}
	}

	return ""
}
