package types

import "io"

type SystemInterface interface {
	// Is syntax highlighting enabled for built-in functions output.
	IsHighlightEnabled() bool
	// Random seed that sets on SystemInit.
	GetRandomSeed() int64
	// Writer for additional output.
	GetStdout() io.Writer

	// Callback that should clears screen.
	ClearScreen()

	// Callback that should return list of available environment file names.
	ListEnvironments() ([]string, error)
	// Callback that should open environment file with provided name for write and return it as io.WriteCloser.
	WriteEnvironment(name string) (io.WriteCloser, error)
	// Callback that should open environment file with provided name for read and return it as io.ReadCloser.
	ReadEnvironment(name string) (io.ReadCloser, error)

	// Callback that should terminate the program and perform any necessary cleanup.
	Exit(errCode int)
}
