package user

import (
	"strings"

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
	var currentFunc Func

	if !HasFunction(name) {
		currentFunc = Func{
			Variants: make([]FuncVariant, 0),
		}
	} else {
		currentFunc = functions[name]
	}

	for i, v := range currentFunc.Variants {
		// if variant with such arguments already exists replace it
		if utils.WordsEqual(v.Args, variant.Args) {
			currentFunc.Variants[i] = variant
			functions[name] = currentFunc
			return
		}
	}

	currentFunc.Variants = append(currentFunc.Variants, variant)
	functions[name] = currentFunc
}

func DeleteFunction(name string) {
	delete(functions, name)
}

func DeleteFunctionVariant(name string, idx int) {
	f := functions[name]
	if idx < len(f.Variants)-1 {
		f.Variants = append(f.Variants[:idx], f.Variants[idx+1:]...)
	} else if idx == len(f.Variants)-1 {
		f.Variants = f.Variants[:idx]
	}
	functions[name] = f
}

func ListFunctions() map[string]Func {
	return functions
}

func DropFunctions() {
	for name := range functions {
		delete(functions, name)
	}
}

func (v FuncVariant) ArgNames() (pos []string) {
	argFound := false
	for _, w := range v.Args {
		if w.Type == utils.W_UNIT && !argFound {
			pos = append(pos, w.Literal)
			argFound = true
		} else if w.Type == utils.W_OP && w.Literal == "," {
			argFound = false
		} else {
			continue
		}
	}
	return
}

func (v FuncVariant) String() string {
	str := "("
	for _, lw := range v.Args {
		str += lw.Literal
	}
	str += ") -> "
	for _, rw := range v.Body {
		str += rw.Literal
	}
	return str
}

func PredictFunction(word string) string {
	for k := range functions {
		if strings.Contains(k, word) {
			return k
		}
	}
	return ""
}
