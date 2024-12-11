package runtime

import (
	"fmt"

	"github.com/dece2183/hexowl/v2/types"
	"github.com/dece2183/hexowl/v2/utils/stack"
)

type Runtime struct {
	ctx       *types.Context
	localVars map[string]interface{}
}

func NewRuntime(ctx *types.Context) *Runtime {
	return &Runtime{
		ctx:       ctx,
		localVars: make(map[string]interface{}),
	}
}

func (rn *Runtime) Reset() {
	rn.localVars = make(map[string]interface{})
}

func (rn *Runtime) SetLocalVariable(name string, val interface{}) {
	rn.localVars[name] = val
}

func (rn *Runtime) GetLocalVariable(name string) (interface{}, bool) {
	val, ok := rn.localVars[name]
	return val, ok
}

func (rn *Runtime) Execute(seq *types.ExecutionSequence) (interface{}, error) {
	valStack := make([]types.Value, 0, 6)
	var opLeft, opRight types.Value

	for _, n := range seq.GetSequence() {
		switch node := n.(type) {
		case types.Value:
			valStack = stack.Push(valStack, node)
		case types.OperatorType:
			if len(valStack) > 0 {
				valStack, opRight = stack.Pop(valStack)
			} else {
				return nil, fmt.Errorf("missing right operand for the %s operator", node.String())
			}

			if len(valStack) > 0 {
				valStack, opLeft = stack.Pop(valStack)
			} else {
				return nil, fmt.Errorf("missing left operand for the %s operator", node.String())
			}

			handler, ok := actionHandlerMap[node]
			if !ok {
				return nil, fmt.Errorf("unknown operator #%d", node)
			}

			val, err := handler(rn, opLeft, opRight)
			if err != nil {
				return nil, err
			}

			valStack = stack.Push(valStack, types.Value{Type: types.V_CONST, Value: val})
		}
	}

	if len(valStack) == 0 {
		return nil, nil
	}

	_, opLeft = stack.Pop(valStack)
	return rn.obtainVariable(opLeft)
}

func (rn *Runtime) assignValue(variable types.Value, val interface{}) error {
	switch variable.Type {
	case types.V_LOCALVAR:
		rn.localVars[variable.Value.(string)] = val
	case types.V_USERVAR:
		rn.ctx.User.SetVariable(variable.Value.(string), val)
	default:
		return fmt.Errorf("%s '%s' is not allowed for assignment", variable.Type.String(), variable.Value.(string))
	}
	return nil
}

func (rn *Runtime) obtainVariable(variable types.Value) (interface{}, error) {
	if variable.Type == types.V_CONST {
		return variable.Value, nil
	}

	var (
		val interface{}
		ok  bool
	)

	varName := variable.Value.(string)

	switch variable.Type {
	case types.V_LOCALVAR:
		if val, ok = rn.localVars[varName]; !ok {
			return nil, fmt.Errorf("'%s' is not a local variable", varName)
		}
	case types.V_USERVAR:
		if val, ok = rn.ctx.User.GetVariable(varName); !ok {
			return nil, fmt.Errorf("'%s' is not a user variable", varName)
		}
	case types.V_BUILTINCONST:
		if val, ok = rn.ctx.Builtin.GetConstant(varName); !ok {
			return nil, fmt.Errorf("'%s' is not a built-in constant", varName)
		}
	case types.V_USERFUNC:
		if val, ok = rn.ctx.User.GetFunction(varName); !ok {
			return nil, fmt.Errorf("'%s' is not a user function", varName)
		}
	case types.V_BUILTINFUNC:
		if val, ok = rn.ctx.Builtin.GetFunction(varName); !ok {
			return nil, fmt.Errorf("'%s' is not a built-in function", varName)
		}
	case types.V_LOCALFUNCPTR:
	case types.V_FUNCPTR:
	}

	return val, nil
}
