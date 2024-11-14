package types

import "io"

type System struct {
	// Is syntax highlighting enabled for built-in functions output.
	HighlightEnabled bool

	// Writer for additional output.
	Stdout io.Writer
	// Callback that should clears screen.
	ClearScreen func()

	// Random seed that sets on SystemInit.
	RandomSeed int64

	// Callback that should return list of available environment file names.
	ListEnvironments func() ([]string, error)
	// Callback that should open environment file with provided name for write and return it as io.WriteCloser.
	WriteEnvironment func(name string) (io.WriteCloser, error)
	// Callback that should open environment file with provided name for read and return it as io.ReadCloser.
	ReadEnvironment func(name string) (io.ReadCloser, error)

	// Callback that should terminate the program and perform any necessary cleanup.
	Exit func(errCode int)
}

type Func struct {
	// Arguments description.
	Args string
	// Function description.
	Desc string
	// Function that will be executed.
	Exec func(desc *Descriptor, args ...interface{}) (interface{}, error)
}

type FunctionMap map[string]Func

type ConstantMap map[string]interface{}

type Descriptor struct {
	Constants ConstantMap
	Functions FunctionMap
	System    System
}
