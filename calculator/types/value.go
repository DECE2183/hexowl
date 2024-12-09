package types

type ValueType uint8

const (
	V_CONST ValueType = iota
	V_LOCALVAR
	V_USERVAR
	V_BUILTINCONST
	V_USERFUNC
	V_BUILTINFUNC
	V_LOCALFUNCPTR
	V_FUNCPTR
)

var valueToStringMap = map[ValueType]string{
	V_CONST:        "constant value",
	V_LOCALVAR:     "local variable",
	V_USERVAR:      "user variable",
	V_BUILTINCONST: "built-in constant",
	V_USERFUNC:     "user function",
	V_BUILTINFUNC:  "built-in function",
	V_LOCALFUNCPTR: "local pointer to function",
	V_FUNCPTR:      "pointer to function",
}

func (v ValueType) String() string {
	str, ok := valueToStringMap[v]
	if !ok {
		return "unknown"
	}
	return str
}

func (v ValueType) IsFunc() bool {
	return v > V_BUILTINCONST
}

type Value struct {
	Type  ValueType
	Value interface{}
}
