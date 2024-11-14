//go:build !nodefsystem
// +build !nodefsystem

package defaultsystem

import (
	"fmt"
	"io"
	"os"
	"path"
	"time"

	"github.com/dece2183/hexowl/builtin"
	"github.com/dece2183/hexowl/builtin/types"
	"github.com/dece2183/hexowl/input/terminal"
)

var DefaultSystem types.System

func init() {
	DefaultSystem = types.System{
		HighlightEnabled: true,
		Stdout:           os.Stdout,
		ClearScreen:      sysClearScreen,
		RandomSeed:       time.Now().UnixNano(),
		ListEnvironments: sysListEnv,
		WriteEnvironment: sysWriteEnv,
		ReadEnvironment:  sysReadEnv,
		Exit:             sysExit,
	}
	builtin.SystemInit(DefaultSystem)
}

func sysClearScreen() {
	fmt.Fprintf(DefaultSystem.Stdout, "\x1bc")
}

func sysGetEnvPath(envName string) (string, error) {
	userDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("unable to get home directory")
	}

	saveDir := path.Join(userDir, ".hexowl/environment")
	err = os.MkdirAll(saveDir, 0755)
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
	err = os.MkdirAll(envDir, 0755)
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
