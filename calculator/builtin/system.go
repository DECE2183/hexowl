package builtin

import (
	"fmt"
	"io"
	"os"
	"path"
	"time"

	"github.com/dece2183/hexowl/v2/input/terminal"
)

type defaultSystem struct{}

func DefaultSystem() defaultSystem {
	return defaultSystem{}
}

func (sys defaultSystem) IsHighlightEnabled() bool {
	return true
}

func (sys defaultSystem) GetRandomSeed() int64 {
	return time.Now().Unix()
}

func (sys defaultSystem) GetStdout() io.Writer {
	return os.Stdout
}

func (sys defaultSystem) ClearScreen() {
	fmt.Printf("\x1bc")
}

func (sys defaultSystem) ListEnvironments() ([]string, error) {
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

func (sys defaultSystem) WriteEnvironment(name string) (io.WriteCloser, error) {
	envPath, err := getEnvPath(name)
	if err != nil {
		return nil, err
	}

	f, err := os.OpenFile(envPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return nil, fmt.Errorf("unable to create file")
	}

	return f, nil
}

func (sys defaultSystem) ReadEnvironment(name string) (io.ReadCloser, error) {
	envPath, err := getEnvPath(name)
	if err != nil {
		return nil, err
	}

	f, err := os.OpenFile(envPath, os.O_RDONLY, 0755)
	if err != nil {
		return nil, fmt.Errorf("unable to open file")
	}

	return f, nil
}

func (sys defaultSystem) Exit(errCode int) {
	terminal.DisableRawMode()
	os.Exit(errCode)
}

func getEnvPath(envName string) (string, error) {
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
