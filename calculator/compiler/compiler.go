package compiler

import (
	"math"
	"strconv"
	"strings"

	"github.com/dece2183/hexowl/v2/calculator/types"
	"github.com/dece2183/hexowl/v2/utils/stack"
)

func Compile(ctx *types.Context, tokens []types.Token) (*types.ExecutionSequence, error) {
	var err error
	var op types.OperatorType

	opStack := make([]types.OperatorType, 0, len(tokens)/3)
	seq := types.NewExecutionSequence()

	if len(tokens) == 0 {
		return seq, nil
	}

	if tokens[len(tokens)-1].Type == types.T_OP {
		op := types.ParseOperator(tokens[len(tokens)-1].Literal)
		return nil, NewCompileError(tokens[len(tokens)-1], len(tokens)-1, "missing right operand for the %s operator", op.String())
	}

	for ti, t := range tokens {
		switch t.Type {
		case types.T_NONE:
			continue
		case types.T_OP:
			op = types.ParseOperator(t.Literal)
			if ti == 0 || tokens[ti-1].Type == types.T_OP {
				if !op.IsUnary() {
					return nil, NewCompileError(t, ti, "missing left operand for the %s operator", op.String())
				}
				seq.AppendValue(types.Value{
					Type:  types.V_CONST,
					Value: 0,
				})
			}
			for len(opStack) > 0 {
				_, sop := stack.Pop(opStack)
				if sop <= op {
					break
				}
				opStack, _ = stack.Pop(opStack)
				seq.AppendOperator(sop)
			}
			if op.IsAssign() {
				lastVal, ok := seq.GetValue(seq.Len() - 1)
				if !ok {
					op, _ = seq.GetOperator(seq.Len() - 1)
					return nil, NewCompileError(t, ti, "%s operator is not allowed for assignment", op.String())
				}
				if lastVal.Type == types.V_CONST {
					return nil, NewCompileError(tokens[ti-1], ti-1, "'%v' is not assignable", lastVal.Value)
				}
				if op == types.O_ASSIGNLOCAL {
					seq.SetLocalVariable(lastVal)
				}
			}
			opStack = stack.Push(opStack, op)
		case types.T_CTL:
			if t.Literal == "(" {
				opStack = stack.Push(opStack, types.O_FLOW)
				if seq.Len() > 0 {
					if val, ok := seq.GetValue(seq.Len() - 1); ok {
						if val.Type.IsFunc() {
							opStack = stack.Push(opStack, types.O_CALLFUNC)
						} else {
							return nil, NewCompileError(tokens[ti-1], ti-1, "'%v' is not a function", val.Value)
						}
					}
				}
			} else {
				if len(opStack) == 0 {
					return nil, NewCompileError(t, ti, "missing opening parenthesis")
				}
				for len(opStack) > 0 {
					opStack, op = stack.Pop(opStack)
					if op == types.O_FLOW {
						break
					}
					seq.AppendOperator(op)
				}
				// detect function call with empty args
				lastOp, ok := seq.GetOperator(seq.Len() - 1)
				if ok && lastOp == types.O_CALLFUNC {
					funcName, ok := seq.GetValue(seq.Len() - 2)
					if ok && funcName.Type.IsFunc() {
						seq.InsertValue(seq.Len()-1, types.Value{
							Type:  types.V_CONST,
							Value: 0,
						})
					}
				}
			}
		case types.T_NUM_SCI:
			var mantisse, order float64
			num := strings.Split(t.Literal, "e")
			mantisse, err = strconv.ParseFloat(strings.ReplaceAll(num[0], "_", ""), 64)
			if err != nil {
				return nil, NewCompileError(t, ti, "unable to parse mantisse part of literal '%s'", t.Literal)
			}
			order, err = strconv.ParseFloat(strings.ReplaceAll(num[1], "_", ""), 64)
			if err != nil {
				return nil, NewCompileError(t, ti, "unable to parse order part of literal '%s'", t.Literal)
			}
			seq.AppendValue(types.Value{
				Type:  types.V_CONST,
				Value: mantisse * math.Pow(10, order),
			})
		case types.T_NUM_DEC:
			var val float64
			val, err = strconv.ParseFloat(strings.ReplaceAll(t.Literal, "_", ""), 64)
			if err != nil {
				return nil, NewCompileError(t, ti, "unable to parse literal '%s' as number", t.Literal)
			}
			seq.AppendValue(types.Value{
				Type:  types.V_CONST,
				Value: val,
			})
		case types.T_NUM_HEX:
			var val uint64
			val, err = strconv.ParseUint(strings.ReplaceAll(t.Literal, "_", ""), 16, 64)
			if err != nil {
				return nil, NewCompileError(t, ti, "unable to parse literal '%s' as hex number", t.Literal)
			}
			seq.AppendValue(types.Value{
				Type:  types.V_CONST,
				Value: val,
			})
		case types.T_NUM_BIN:
			var val uint64
			val, err = strconv.ParseUint(strings.ReplaceAll(t.Literal, "_", ""), 2, 64)
			if err != nil {
				return nil, NewCompileError(t, ti, "unable to parse literal '%s' as bin number", t.Literal)
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
				valType = types.V_VARNAME
			}
			seq.AppendValue(types.Value{
				Type:  valType,
				Value: t.Literal,
			})
		default:
			return nil, NewCompileError(t, ti, "unknown token #%d", t.Type)
		}
	}

	for len(opStack) > 0 {
		opStack, op = stack.Pop(opStack)
		seq.AppendOperator(op)
	}

	return seq, nil
}
