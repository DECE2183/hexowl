package builtin

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"path"
	"time"

	"github.com/dece2183/hexowl/input/terminal"
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

var defaultSystem = System{
	HighlightEnabled: true,
	Stdout:           os.Stdout,
	ClearScreen:      sysClearScreen,
	RandomSeed:       time.Now().UnixNano(),
	ListEnvironments: sysListEnv,
	WriteEnvironment: sysWriteEnv,
	ReadEnvironment:  sysReadEnv,
	Exit:             sysExit,
}

var bDesc descriptor

func init() {
	bDesc.functions = &functions
	SystemInit(defaultSystem)
}

func sysClearScreen() {
	fmt.Fprintf(bDesc.system.Stdout, "\x1bc")
}

func sysGetEnvPath(envName string) (string, error) {
	userDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("unable to get home directory")
	}

	saveDir := path.Join(userDir, ".hexowl/environment")
	err = os.MkdirAll(saveDir, 0666)
	if err != nil {
		return "", fmt.Errorf("unable to make directory")
	}

	return path.Join(saveDir, envName) + ".json", nil
}

func sysListEnv() ([]string, error) {
	userDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("unable to get home directory")
	}
	envDir := path.Join(userDir, ".hexowl/environment")
	err = os.MkdirAll(envDir, 0666)
	if err != nil {
		return nil, fmt.Errorf("unable to make directory")
	}

	dir, err := os.Open(envDir)
	if err != nil {
		return nil, fmt.Errorf("unable to open directory")
	}
	defer dir.Close()

	files, err := dir.Readdir(0)
	if err != nil {
		return nil, fmt.Errorf("unable to read directory")
	}

	envs := make([]string, 0, len(files))
	for _, f := range files {
		if f.IsDir() {
			continue
		}

		name := f.Name()
		ext := path.Ext(name)
		if ext != ".json" {
			continue
		}

		envs = append(envs, name[:len(name)-len(ext)])
	}

	return envs, nil
}

func sysWriteEnv(name string) (io.WriteCloser, error) {
	envPath, err := sysGetEnvPath(name)
	if err != nil {
		return nil, err
	}

	f, err := os.OpenFile(envPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return nil, fmt.Errorf("unable to create file")
	}

	return f, nil
}

func sysReadEnv(name string) (io.ReadCloser, error) {
	envPath, err := sysGetEnvPath(name)
	if err != nil {
		return nil, err
	}

	f, err := os.OpenFile(envPath, os.O_RDONLY, 0755)
	if err != nil {
		return nil, fmt.Errorf("unable to open file")
	}

	return f, nil
}

func sysExit(errCode int) {
	terminal.DisableRawMode()
	os.Exit(errCode)
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
