package types

type BuiltinFunction struct {
	// Arguments description.
	Args string
	// Function description.
	Desc string
	// Function that will be executed.
	Exec func(ctx *Context, args []interface{}) (interface{}, error)
}

type BuiltinFunctionMap map[string]BuiltinFunction

type BuiltinConstantMap map[string]any

type BuiltinContainer struct {
	functions BuiltinFunctionMap
	constants BuiltinConstantMap
}

// Is function with name presented in the builtin function map.
func (b *BuiltinContainer) HasFunction(name string) bool {
	_, found := b.functions[name]
	return found
}

// Register a new function and add it to the builtin function map.
func (b *BuiltinContainer) RegisterFunction(name string, function BuiltinFunction) {
	b.functions[name] = function
}

// Get function by name from the builtin function map.
func (b *BuiltinContainer) GetFunction(name string) (function BuiltinFunction, found bool) {
	function, found = b.functions[name]
	return
}

// Return the builtin function map.
func (b *BuiltinContainer) ListFunctions() BuiltinFunctionMap {
	return b.functions
}

// Is constant with name presented in the builtin constant map.
func (b *BuiltinContainer) HasConstant(name string) bool {
	_, found := b.constants[name]
	return found
}

// Register a new constant and add it to the builtin constant map.
func (b *BuiltinContainer) RegisterConstant(name string, value interface{}) {
	b.constants[name] = value
}

// Get constant by name from the builtin constant map.
func (b *BuiltinContainer) GetConstant(name string) (val interface{}, found bool) {
	val, found = b.constants[name]
	return
}

// Return the builtin constant map.
func (b *BuiltinContainer) ListConstants() BuiltinConstantMap {
	return b.constants
}

func (b *BuiltinContainer) Predict(word string) string {
	for c := range b.constants {
		if len(c) < len(word) {
			continue
		}
		if c[:len(word)] == word {
			return c
		}
	}

	for f := range b.functions {
		if len(f) < len(word) {
			continue
		}
		if f[:len(word)] == word {
			return f + "()"
		}
	}

	return ""
}
