package types

type ValueType uint8

const (
	V_CONST ValueType = iota
	V_VARNAME
	V_FUNCNAME
	V_LOCALVAR
	V_USERVAR
	V_BUILTINCONST
	V_USERFUNC
	V_BUILTINFUNC
	V_LOCALFUNCPTR
	V_FUNCPTR
	V_FUNCARG
	V_FUNCBODY
)

var valueToStringMap = map[ValueType]string{
	V_CONST:        "constant value",
	V_VARNAME:      "variable name",
	V_FUNCNAME:     "function name",
	V_LOCALVAR:     "local variable",
	V_USERVAR:      "user variable",
	V_BUILTINCONST: "built-in constant",
	V_USERFUNC:     "user function",
	V_BUILTINFUNC:  "built-in function",
	V_LOCALFUNCPTR: "local pointer to function",
	V_FUNCPTR:      "pointer to function",
	V_FUNCARG:      "function arguments",
	V_FUNCBODY:     "function body",
}

func (v ValueType) String() string {
	str, ok := valueToStringMap[v]
	if !ok {
		return "unknown"
	}
	return str
}

func (v ValueType) IsAssignable() bool {
	return v > V_CONST && v < V_BUILTINCONST
}

func (v ValueType) IsFunc() bool {
	return v == V_USERFUNC || v == V_BUILTINFUNC
}

func (v ValueType) IsFuncPtr() bool {
	return v == V_LOCALFUNCPTR || v == V_FUNCPTR
}

type Value struct {
	Type       ValueType
	Value      interface{}
	TokenIndex int
}
