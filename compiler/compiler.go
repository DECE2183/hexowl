package compiler

import (
	"math"
	"strconv"
	"strings"

	"github.com/dece2183/hexowl/v2/types"
	"github.com/dece2183/hexowl/v2/utils/stack"
)

type assignment struct {
	tokenPos    int
	sequencePos int
}

func Compile(ctx *types.Context, tokens []types.Token) (*types.ExecutionSequence, error) {
	var err error
	var op types.Operator

	declarationStack := make([]assignment, 0)
	functionStack := make([]assignment, 0)
	opStack := make([]types.Operator, 0, len(tokens)/3)
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
			opType := types.ParseOperator(t.Literal)
			if ti == 0 || tokens[ti-1].Type == types.T_OP {
				if !opType.IsUnary() {
					return nil, NewCompileError(t, ti, "missing left operand for the %s operator", opType.String())
				}
				seq.AppendValue(types.Value{
					Type:       types.V_CONST,
					Value:      0,
					TokenIndex: -1,
				})
			}
			for len(opStack) > 0 {
				_, sop := stack.Pop(opStack)
				if sop.Type <= opType {
					break
				}
				opStack, _ = stack.Pop(opStack)
				if sop.Type == types.O_DECLFUNC {
					var fn assignment
					declarationStack, fn = stack.Pop(declarationStack)
					bodySeq := seq.ExtractSubsequence(fn.sequencePos, seq.Len())
					seq.AppendValue(types.Value{
						Type:       types.V_CONST,
						Value:      bodySeq,
						TokenIndex: -1,
					})
				}
				seq.AppendOperator(sop)
			}
			if opType.IsAssign() {
				lastVal, ok := seq.GetValue(seq.Len() - 1)
				if !ok {
					return nil, NewCompileError(t, ti, "there is no variable for assignment")
				}
				if lastVal.Type == types.V_CONST {
					return nil, NewCompileError(tokens[ti-1], ti-1, "'%v' is not assignable", lastVal.Value)
				}
				if opType == types.O_ASSIGNLOCAL {
					lastVal.Type = types.V_LOCALVAR
				} else if opType == types.O_ASSIGN {
					lastVal.Type = types.V_USERVAR
				}
				seq.SetValue(seq.Len()-1, lastVal)
			} else if opType == types.O_DECLFUNC {
				declarationStack = stack.Push(declarationStack, assignment{
					tokenPos:    ti,
					sequencePos: seq.Len(),
				})
			}
			opStack = stack.Push(opStack, types.Operator{
				Type:       opType,
				TokenIndex: ti,
			})
		case types.T_CTL:
			if t.Literal == "(" {
				opStack = stack.Push(opStack, types.Operator{
					Type:       types.O_FLOW,
					TokenIndex: ti,
				})
				if seq.Len() > 0 {
					if val, ok := seq.GetValue(seq.Len() - 1); ok {
						if val.Type.IsFunc() || val.Type == types.V_FUNCNAME {
							opStack = stack.Push(opStack, types.Operator{
								Type:       types.O_CALLFUNC,
								TokenIndex: ti,
							})
						}
					}
				}
			} else {
				var flowFound bool
				for len(opStack) > 0 {
					opStack, op = stack.Pop(opStack)
					if op.Type == types.O_FLOW {
						flowFound = true
						break
					} else if op.Type == types.O_DECLFUNC {
						var fn assignment
						declarationStack, fn = stack.Pop(declarationStack)
						bodySeq := seq.ExtractSubsequence(fn.sequencePos, seq.Len())
						seq.AppendValue(types.Value{
							Type:       types.V_CONST,
							Value:      bodySeq,
							TokenIndex: -1,
						})
					}
					seq.AppendOperator(op)
				}
				if !flowFound {
					return nil, NewCompileError(t, ti, "missing opening parenthesis")
				}
				lastOp, ok := seq.GetOperator(seq.Len() - 1)
				if ok && lastOp.Type == types.O_CALLFUNC {
					var fn assignment
					functionStack, fn = stack.Pop(functionStack)
					funcName, _ := seq.GetValue(fn.sequencePos)

					if ti < len(tokens)-1 && tokens[ti+1].Type == types.T_OP && types.ParseOperator(tokens[ti+1].Literal) == types.O_DECLFUNC {
						// function declaration
						argsSeq := seq.ExtractSubsequence(fn.sequencePos, seq.Len())
						seq.AppendValue(types.Value{
							Type:       types.V_FUNCARG,
							Value:      argsSeq,
							TokenIndex: -1,
						})
					} else {
						// function call
						if !funcName.Type.IsFunc() {
							return nil, NewCompileError(tokens[funcName.TokenIndex], funcName.TokenIndex, "'%v' is not a function", funcName.Value)
						}
						// detect function call with empty args
						if fn.tokenPos+2 == ti {
							seq.InsertValue(seq.Len()-1, types.Value{
								Type:       types.V_CONST,
								Value:      0,
								TokenIndex: -1,
							})
						}
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
				Type:       types.V_CONST,
				Value:      mantisse * math.Pow(10, order),
				TokenIndex: ti,
			})
		case types.T_NUM_DEC:
			var val float64
			val, err = strconv.ParseFloat(strings.ReplaceAll(t.Literal, "_", ""), 64)
			if err != nil {
				return nil, NewCompileError(t, ti, "unable to parse literal '%s' as number", t.Literal)
			}
			seq.AppendValue(types.Value{
				Type:       types.V_CONST,
				Value:      val,
				TokenIndex: ti,
			})
		case types.T_NUM_HEX:
			var val uint64
			val, err = strconv.ParseUint(strings.ReplaceAll(t.Literal, "_", ""), 16, 64)
			if err != nil {
				return nil, NewCompileError(t, ti, "unable to parse literal '%s' as hex number", t.Literal)
			}
			seq.AppendValue(types.Value{
				Type:       types.V_CONST,
				Value:      val,
				TokenIndex: ti,
			})
		case types.T_NUM_BIN:
			var val uint64
			val, err = strconv.ParseUint(strings.ReplaceAll(t.Literal, "_", ""), 2, 64)
			if err != nil {
				return nil, NewCompileError(t, ti, "unable to parse literal '%s' as bin number", t.Literal)
			}
			seq.AppendValue(types.Value{
				Type:       types.V_CONST,
				Value:      val,
				TokenIndex: ti,
			})
		case types.T_STR:
			seq.AppendValue(types.Value{
				Type:       types.V_CONST,
				Value:      t.Literal,
				TokenIndex: ti,
			})
		case types.T_UNIT:
			var valType types.ValueType

			// Try to find variable
			if ti < len(tokens)-1 && tokens[ti+1].Type == types.T_CTL && tokens[ti+1].Literal == "(" {
				if seq.HasUserFunction(t.Literal) || ctx.User.HasFunction(t.Literal) {
					valType = types.V_USERFUNC
				} else if ctx.Builtin.HasFunction(t.Literal) {
					valType = types.V_BUILTINFUNC
				} else {
					valType = types.V_FUNCNAME
				}
			} else {
				if seq.HasLocalVariable(t.Literal) {
					valType = types.V_LOCALVAR
				} else if seq.HasUserVariable(t.Literal) || ctx.User.HasVariable(t.Literal) {
					valType = types.V_USERVAR
				} else if ctx.Builtin.HasConstant(t.Literal) {
					valType = types.V_BUILTINCONST
				} else {
					valType = types.V_VARNAME
				}
			}

			if valType.IsFunc() || valType == types.V_FUNCNAME {
				functionStack = stack.Push(functionStack, assignment{
					tokenPos:    ti,
					sequencePos: seq.Len(),
				})
			}

			seq.AppendValue(types.Value{
				Type:       valType,
				Value:      t.Literal,
				TokenIndex: ti,
			})
		default:
			return nil, NewCompileError(t, ti, "unknown token #%d", t.Type)
		}
	}

	for len(opStack) > 0 {
		opStack, op = stack.Pop(opStack)
		if op.Type == types.O_DECLFUNC {
			var fn assignment
			declarationStack, fn = stack.Pop(declarationStack)
			bodySeq := seq.ExtractSubsequence(fn.sequencePos, seq.Len())
			seq.AppendValue(types.Value{
				Type:       types.V_FUNCBODY,
				Value:      bodySeq,
				TokenIndex: -1,
			})
		}
		seq.AppendOperator(op)
	}

	return seq, nil
}
