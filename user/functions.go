package user

import (
	"github.com/dece2183/hexowl/utils"
)

type FuncVariant struct {
	Args []utils.Word
	Body []utils.Word
}

type Func struct {
	Variants []FuncVariant
}

var functions = map[string]Func{}

func HasFunction(name string) bool {
	_, found := functions[name]
	return found
}

func GetFunction(name string) (function Func, found bool) {
	function, found = functions[name]
	return
}

func SetFunction(name string, function Func) {
	functions[name] = function
}

func SetFunctionVariant(name string, variant FuncVariant) {
	if !HasFunction(name) {
		functions[name] = Func{
			Variants: make([]FuncVariant, 0),
		}
	}
	currentFunc := functions[name]
	currentFunc.Variants = append(currentFunc.Variants, variant)
	functions[name] = currentFunc
}

func DeleteFunction(name string) {
	delete(functions, name)
}

func ListFunctions() map[string]Func {
	return functions
}

func DropFunctions() {
	for name := range functions {
		delete(functions, name)
	}
}
