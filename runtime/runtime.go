package runtime

import (
	"fmt"
	"slices"

	"github.com/dece2183/hexowl/v2/types"
	"github.com/dece2183/hexowl/v2/utils/stack"
)

type Runtime struct {
	ctx       *types.Context
	localVars []interface{}
}

func NewRuntime(ctx *types.Context) *Runtime {
	return &Runtime{
		ctx:       ctx,
		localVars: make([]interface{}, 0),
	}
}

func (rn *Runtime) Reset() {
	rn.localVars = make([]interface{}, 0)
}

func (rn *Runtime) SetLocalVariable(index int, val interface{}) {
	rn.localVars = slices.Grow(rn.localVars, index+1)
	for index >= len(rn.localVars) {
		rn.localVars = append(rn.localVars, nil)
	}
	rn.localVars[index] = val
}

func (rn *Runtime) GetLocalVariable(index int) (interface{}, bool) {
	if index >= len(rn.localVars) {
		return nil, false
	}
	return rn.localVars[index], true
}

func (rn *Runtime) Execute(seq *types.ExecutionSequence) (interface{}, error) {
	valStack := make([]types.Value, 0, 6)
	var opLeft, opRight types.Value

	for _, n := range seq.GetSequence() {
		switch node := n.(type) {
		case types.Value:
			valStack = stack.Push(valStack, node)
		case types.Operator:
			if len(valStack) > 0 {
				valStack, opRight = stack.Pop(valStack)
			} else {
				return nil, fmt.Errorf("missing right operand for the %s operator", node.Type.String())
			}

			if len(valStack) > 0 {
				valStack, opLeft = stack.Pop(valStack)
			} else {
				return nil, fmt.Errorf("missing left operand for the %s operator", node.Type.String())
			}

			handler, ok := actionHandlerMap[node.Type]
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

func (rn *Runtime) ExecuteUserFunction(fn types.UserFunction, args []interface{}) (interface{}, error) {
variant:
	for vari := range fn.Variants {
		v := &fn.Variants[vari]
		argsLen := v.Args.Sequence.LocalVariablesLen()
		vArgsPos, hasVArgs := v.Args.Sequence.GetLocalVariableIndex(types.K_VARGS)
		if argsLen > 0 && argsLen != len(args) && !hasVArgs {
			continue
		}
		newRn := NewRuntime(rn.ctx)
		if hasVArgs {
			for i := 0; i < vArgsPos; i++ {
				newRn.SetLocalVariable(i, args[i])
			}
			for i := 1; i < argsLen-vArgsPos; i++ {
				newRn.SetLocalVariable(argsLen-i, args[len(args)-i])
			}
			newRn.SetLocalVariable(vArgsPos, args[vArgsPos:len(args)-(argsLen-vArgsPos)+1])
		} else {
			for i := 0; i < argsLen; i++ {
				newRn.SetLocalVariable(i, args[i])
			}
		}
		res, err := newRn.Execute(v.Args.Sequence)
		if err != nil {
			return nil, err
		}
		switch result := res.(type) {
		case []interface{}:
			for i := range result {
				if success, ok := result[i].(bool); ok && !success {
					continue variant
				}
			}
		case bool:
			if !result {
				continue
			}
		}

		return newRn.Execute(v.Body.Sequence)
	}

	return nil, fmt.Errorf("unable to find user function variant")
}

func (rn *Runtime) assignValue(variable types.Value, val interface{}) error {
	switch varname := variable.Value.(type) {
	case int:
		if varname >= len(rn.localVars) {
			return fmt.Errorf("%s is not assignable", variable.Type.String())
		}
		rn.localVars[varname] = val
	case string:
		rn.ctx.User.SetVariable(varname, val)
	default:
		return fmt.Errorf("%s is not assignable", variable.Type.String())
	}
	return nil
}

func (rn *Runtime) assignLocalValue(variable types.Value, val interface{}) error {
	varIndex, ok := variable.Value.(int)
	if !ok {
		return fmt.Errorf("%s is not assignable", variable.Type.String())
	}

	rn.SetLocalVariable(varIndex, val)
	return nil
}

func (rn *Runtime) obtainVariable(variable types.Value) (interface{}, error) {
	if variable.Type == types.V_CONST {
		return variable.Value, nil
	} else if variable.Type == types.V_LOCALVAR {
		varIndex := variable.Value.(int)
		if varIndex >= len(rn.localVars) {
			return nil, fmt.Errorf("there is no #%d local variable", varIndex)
		}
		return rn.localVars[varIndex], nil
	}

	var (
		val interface{}
		ok  bool
	)

	varName := variable.Value.(string)

	switch variable.Type {
	case types.V_USERVAR:
		if val, ok = rn.ctx.User.GetVariable(varName); !ok {
			return nil, fmt.Errorf("not found user variable '%s'", varName)
		}
	case types.V_BUILTINCONST:
		if val, ok = rn.ctx.Builtin.GetConstant(varName); !ok {
			return nil, fmt.Errorf("not found built-in constant '%s'", varName)
		}
	case types.V_USERFUNC:
		if val, ok = rn.ctx.User.GetFunction(varName); !ok {
			return nil, fmt.Errorf("not found user function '%s'", varName)
		}
	case types.V_BUILTINFUNC:
		if val, ok = rn.ctx.Builtin.GetFunction(varName); !ok {
			return nil, fmt.Errorf("not found built-in function '%s'", varName)
		}
	case types.V_LOCALFUNCPTR:
	case types.V_FUNCPTR:
	}

	return val, nil
}
