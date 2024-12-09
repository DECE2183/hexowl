package impl

import (
	"github.com/dece2183/hexowl/v2/calculator/types"
	"github.com/dece2183/hexowl/v2/utils"
)

func Clear(ctx *types.Context, args []interface{}) (interface{}, error) {
	ctx.System.ClearScreen()
	return nil, nil
}

func Exit(ctx *types.Context, args []interface{}) (interface{}, error) {
	exitCode := utils.ToNumber[int64](args[0])
	ctx.System.Exit(int(exitCode))
	return exitCode, nil
}
