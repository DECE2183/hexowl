package user

var variables = map[string]interface{}{}

// Is variable with name presented in the user variables map.
func HasVariable(name string) bool {
	_, found := variables[name]
	return found
}

// Set user variable with name and value.
func SetVariable(name string, val interface{}) {
	variables[name] = val
}

// Get user variable with given name.
func GetVariable(name string) (val interface{}, found bool) {
	val, found = variables[name]
	return
}

// Delete user variable with given name.
func DeleteVariable(name string) {
	delete(variables, name)
}

// Return the user variables map.
func ListVariables() map[string]interface{} {
	return variables
}

// Delete all user variables.
func DropVariables() {
	for name := range variables {
		delete(variables, name)
	}
}
