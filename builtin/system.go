package builtin

import (
	"io"
	"math/rand"

	"github.com/dece2183/hexowl/user"
)

type Environment struct {
	Description string
	UserVars    map[string]interface{}
	UserFuncs   map[string]user.Func
}

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

type descriptor struct {
	functions *FuncMap
	system    System
}

var bDesc descriptor

func init() {
	bDesc.functions = &functions
	bDesc.system.Stdout = io.Discard
}

// Provide your own system description.
//
// Use this function to implement native integration into your application.
func SystemInit(sys System) {
	bDesc.system = sys
	if sys.RandomSeed != 0 {
		rand.Seed(sys.RandomSeed)
	}
}
