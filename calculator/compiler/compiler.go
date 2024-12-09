package compiler

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/dece2183/hexowl/calculator/compiler/stack"
	"github.com/dece2183/hexowl/calculator/types"
)

func Compile(ctx *types.Context, tokens []types.Token) (*types.ExecutionSequence, error) {
	var err error

	var op types.OperatorType
	opStack := make([]types.OperatorType, 0, len(tokens)/3)
	seq := types.NewExecutionSequence()

	for _, t := range tokens {
		switch t.Type {
		case types.T_NONE:
			seq.AppendValue(types.Value{
				Type:  types.V_CONST,
				Value: nil,
			})
		case types.T_OP:
			op = types.ParseOperator(t.Literal)
			if op == types.O_ASSIGNLOCAL {
				val, _ := seq.GetValue(seq.Len() - 1)
				seq.SetLocalVariable(val.Value.(string))
			}
			for len(opStack) > 0 {
				_, sop := stack.Pop(opStack)
				if sop <= op {
					break
				}
				opStack, _ = stack.Pop(opStack)
				seq.AppendOperator(sop)
			}
			opStack = stack.Push(opStack, op)
		case types.T_CTL:
			if t.Literal == "(" {
				opStack = stack.Push(opStack, types.O_FLOW)
				if seq.Len() > 0 {
					if val, ok := seq.GetValue(seq.Len() - 1); ok && val.Type.IsFunc() {
						opStack = stack.Push(opStack, types.O_CALLFUNC)
					}
				}
			} else {
				for len(opStack) > 0 {
					opStack, op = stack.Pop(opStack)
					if op == types.O_FLOW {
						break
					}
					seq.AppendOperator(op)
				}
			}
		case types.T_NUM_SCI:
			var mantisse, order float64
			num := strings.Split(t.Literal, "e")
			mantisse, err = strconv.ParseFloat(strings.ReplaceAll(num[0], "_", ""), 64)
			if err != nil {
				return nil, fmt.Errorf("unable to parse mantisse part of literal '%s'", t.Literal)
			}
			order, err = strconv.ParseFloat(strings.ReplaceAll(num[1], "_", ""), 64)
			if err != nil {
				return nil, fmt.Errorf("unable to parse order part of literal '%s'", t.Literal)
			}
			seq.AppendValue(types.Value{
				Type:  types.V_CONST,
				Value: mantisse * math.Pow(10, order),
			})
		case types.T_NUM_DEC:
			var val float64
			val, err = strconv.ParseFloat(strings.ReplaceAll(t.Literal, "_", ""), 64)
			if err != nil {
				return nil, fmt.Errorf("unable to parse literal '%s' as number", t.Literal)
			}
			seq.AppendValue(types.Value{
				Type:  types.V_CONST,
				Value: val,
			})
		case types.T_NUM_HEX:
			var val uint64
			val, err = strconv.ParseUint(strings.ReplaceAll(t.Literal, "_", ""), 16, 64)
			if err != nil {
				return nil, fmt.Errorf("unable to parse literal '%s' as hex number", t.Literal)
			}
			seq.AppendValue(types.Value{
				Type:  types.V_CONST,
				Value: val,
			})
		case types.T_NUM_BIN:
			var val uint64
			val, err = strconv.ParseUint(strings.ReplaceAll(t.Literal, "_", ""), 2, 64)
			if err != nil {
				return nil, fmt.Errorf("unable to parse literal '%s' as bin number", t.Literal)
			}
			seq.AppendValue(types.Value{
				Type:  types.V_CONST,
				Value: val,
			})
		case types.T_STR:
			seq.AppendValue(types.Value{
				Type:  types.V_CONST,
				Value: t.Literal,
			})
		case types.T_UNIT:
			// Try to find variable
			var valType types.ValueType
			if seq.HasLocalVariable(t.Literal) {
				valType = types.V_LOCALVAR
			} else if ctx.User.HasVariable(t.Literal) {
				valType = types.V_USERVAR
			} else if ctx.Builtin.HasConstant(t.Literal) {
				valType = types.V_BUILTINCONST
			} else if ctx.User.HasFunction(t.Literal) {
				valType = types.V_USERFUNC
			} else if ctx.Builtin.HasFunction(t.Literal) {
				valType = types.V_BUILTINFUNC
			} else {
				valType = types.V_CONST
			}
			seq.AppendValue(types.Value{
				Type:  valType,
				Value: t.Literal,
			})
		case types.T_FUNC:
			// Try to find function
			var valType types.ValueType
			if seq.HasLocalVariable(t.Literal) {
				valType = types.V_LOCALFUNCPTR
			} else if ctx.User.HasVariable(t.Literal) {
				valType = types.V_FUNCPTR
			} else if ctx.User.HasFunction(t.Literal) {
				valType = types.V_USERFUNC
			} else if ctx.Builtin.HasFunction(t.Literal) {
				valType = types.V_BUILTINFUNC
			} else {
				return nil, fmt.Errorf("there is no function named '%s'", t.Literal)
			}
			seq.AppendValue(types.Value{
				Type:  valType,
				Value: t.Literal,
			})
		default:
			return nil, fmt.Errorf("unknown token #%d", t.Type)
		}
	}

	for len(opStack) > 0 {
		opStack, op = stack.Pop(opStack)
		seq.AppendOperator(op)
	}

	return seq, nil
}
