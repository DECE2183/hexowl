package calculator

import (
	"github.com/dece2183/hexowl/v2/calculator/builtin"
	"github.com/dece2183/hexowl/v2/calculator/compiler"
	"github.com/dece2183/hexowl/v2/calculator/lexer"
	"github.com/dece2183/hexowl/v2/calculator/runtime"
	"github.com/dece2183/hexowl/v2/calculator/types"
)

type Calculator struct {
	ctx *types.Context
}

func NewCalculator(system types.SystemInterface) *Calculator {
	ctx := types.NewEmptyContext(system)
	builtin.RegisterConstants(ctx)
	builtin.RegisterFunctions(ctx)
	return &Calculator{
		ctx: ctx,
	}
}

func (calc *Calculator) GetBuiltinContainer() *types.BuiltinContainer {
	return &calc.ctx.Builtin
}

func (calc *Calculator) GetUserContainer() *types.UserContainer {
	return &calc.ctx.User
}

func (calc *Calculator) Eval(str string) (interface{}, error) {
	tokens := lexer.Parse(str)

	seq, err := compiler.Compile(calc.ctx, tokens)
	if err != nil {
		return nil, err
	}

	rn := runtime.NewRuntime(calc.ctx)
	result, err := rn.Execute(seq)
	if err != nil {
		return nil, err
	}

	return result, nil
}
