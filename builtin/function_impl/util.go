package functionimpl

import (
	"fmt"

	"github.com/dece2183/hexowl/builtin/types"
	"github.com/dece2183/hexowl/utils"
)

func Clear(desc *types.Descriptor, args ...interface{}) (interface{}, error) {
	if desc.System.ClearScreen == nil {
		return nil, fmt.Errorf("'clear' not implemented")
	}

	desc.System.ClearScreen()
	return nil, nil
}

func Exit(desc *types.Descriptor, args ...interface{}) (interface{}, error) {
	if desc.System.Exit == nil {
		return nil, fmt.Errorf("'exit' not implemented")
	}

	exitCode := utils.ToNumber[int64](args[0])
	desc.System.Exit(int(exitCode))
	return exitCode, nil
}
