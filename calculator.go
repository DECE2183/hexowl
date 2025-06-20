package hexowl

import (
	"strings"

	"github.com/dece2183/hexowl/v2/builtin"
	"github.com/dece2183/hexowl/v2/compiler"
	"github.com/dece2183/hexowl/v2/lexer"
	"github.com/dece2183/hexowl/v2/runtime"
	"github.com/dece2183/hexowl/v2/types"
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

func (calc *Calculator) Eval(str string) (result interface{}, err error) {
	lines := strings.Split(str, "\n")
	for _, ln := range lines {
		tokens := lexer.Parse(ln)

		var seq *types.ExecutionSequence
		seq, err = compiler.Compile(calc.ctx, tokens)
		if err != nil {
			return
		}

		rn := runtime.NewRuntime(calc.ctx)
		result, err = rn.Execute(seq)
		if err != nil {
			return
		}
	}
	return
}
