package user

import "strings"

var variables = map[string]interface{}{}

func HasVariable(name string) bool {
	_, found := variables[name]
	return found
}

func SetVariable(name string, val interface{}) {
	variables[name] = val
}

func GetVariable(name string) (val interface{}, found bool) {
	val, found = variables[name]
	return
}

func DeleteVariable(name string) {
	delete(variables, name)
}

func ListVariables() map[string]interface{} {
	return variables
}

func DropVariables() {
	for name := range variables {
		delete(variables, name)
	}
}

func PredictVariable(word string) string {
	for k := range variables {
		if strings.Contains(k, word) {
			return k
		}
	}
	return ""
}
