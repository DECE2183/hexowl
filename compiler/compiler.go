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
					declarationStack, err = declFuncBody(seq, tokens, declarationStack)
					if err != nil {
						return nil, err
					}
				}
				seq.AppendOperator(sop)
			}
			if opType.IsAssign() {
				lastVal, ok := seq.GetValue(seq.Len() - 1)
				if !ok {
					return nil, NewCompileError(t, ti, "there is no variable for assignment")
				}
				if !lastVal.Type.IsAssignable() || lastVal.Value.(string) == types.K_VARGS {
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
					Type:       types.O_FLOWBEG,
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
					if op.Type == types.O_FLOWBEG {
						flowFound = true
						break
					} else if op.Type == types.O_DECLFUNC {
						declarationStack, err = declFuncBody(seq, tokens, declarationStack)
						if err != nil {
							return nil, err
						}
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
						// mark this function as user defined
						funcName.Type = types.V_USERFUNC
						seq.SetValue(fn.sequencePos, funcName)

						// detect arguments declaration
						for i := fn.sequencePos + 1; i < seq.Len(); i++ {
							v, isVal := seq.GetValue(i)
							if !isVal {
								continue
							}

							if v.Type == types.V_VARNAME {
								t = types.Token{
									Type:    types.T_UNIT,
									Literal: v.Value.(string),
								}
								return nil, NewCompileError(t, v.TokenIndex, "illegal definition in argument list '%s'", t.Literal)
							} else if v.Type == types.V_UNKNOWN {
								nextOp, isOp := seq.GetOperator(i + 1)
								if isOp && !nextOp.Type.IsUnary() && nextOp.Type != types.O_ENUMERATE && nextOp.Type != types.O_CALLFUNC && !seq.HasLocalVariable(v.Value.(string)) {
									t = types.Token{
										Type:    types.T_UNIT,
										Literal: v.Value.(string),
									}
									return nil, NewCompileError(t, v.TokenIndex, "unknown variable '%s'", t.Literal)
								}
								v.Type = types.V_LOCALVAR
							}

							if v.Type != types.V_LOCALVAR {
								continue
							}

							v.Type = types.V_LOCALVAR
							seq.SetValue(i, v)
						}
						// function declaration
						argsSeq, err := seq.ExtractSubsequence(fn.sequencePos, seq.Len())
						if err != nil {
							err := err.(*types.BadSequence)
							t := types.Token{
								Type:    types.T_UNIT,
								Literal: err.Value.Value.(string),
							}
							return nil, NewCompileError(t, err.Value.TokenIndex, "unknown variable '%s'", t.Literal)
						}
						seq.AppendValue(types.Value{
							Type: types.V_FUNCARG,
							Value: types.UserFunctionPart{
								Definition: tokens[fn.tokenPos : ti+1],
								Sequence:   argsSeq,
							},
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

			var nextOp types.OperatorType
			if ti < len(tokens)-1 {
				nextOp = types.ParseOperator(tokens[ti+1].Literal)
			}

			if nextOp == types.O_FLOWBEG {
				// Try to find function
				if seq.HasUserFunction(t.Literal) || ctx.User.HasFunction(t.Literal) {
					valType = types.V_USERFUNC
				} else if ctx.Builtin.HasFunction(t.Literal) {
					valType = types.V_BUILTINFUNC
				} else {
					valType = types.V_FUNCNAME
				}

				functionStack = stack.Push(functionStack, assignment{
					tokenPos:    ti,
					sequencePos: seq.Len(),
				})
			} else {
				// Try to find variable
				if seq.HasLocalVariable(t.Literal) {
					valType = types.V_LOCALVAR
				} else if seq.HasUserVariable(t.Literal) || ctx.User.HasVariable(t.Literal) {
					valType = types.V_USERVAR
				} else if ctx.Builtin.HasConstant(t.Literal) {
					valType = types.V_BUILTINCONST
				} else if nextOp == types.O_ASSIGN || nextOp == types.O_ASSIGNLOCAL {
					valType = types.V_VARNAME
				} else if seq.HasUserFunction(t.Literal) || ctx.User.HasFunction(t.Literal) {
					valType = types.V_USERFUNC
				} else if ctx.Builtin.HasFunction(t.Literal) {
					valType = types.V_BUILTINFUNC
				} else {
					valType = types.V_UNKNOWN
					// return nil, NewCompileError(t, ti, "unknown variable '%s'", t.Literal)
				}
			}

			seq.AppendValue(types.Value{
				Type:       valType,
				Value:      t.Literal,
				TokenIndex: ti,
			})
		default:
			return nil, NewCompileError(t, ti, "unknown token #%d '%s'", t.Type, t.Literal)
		}
	}

	for len(opStack) > 0 {
		opStack, op = stack.Pop(opStack)
		if op.Type == types.O_DECLFUNC {
			_, err = declFuncBody(seq, tokens, declarationStack)
			if err != nil {
				return nil, err
			}
		}
		seq.AppendOperator(op)
	}

	return seq, nil
}

func declFuncBody(seq *types.ExecutionSequence, tokens []types.Token, declStack []assignment) ([]assignment, error) {
	var fn assignment
	declStack, fn = stack.Pop(declStack)
	bodySeq, err := seq.ExtractSubsequence(fn.sequencePos, seq.Len())
	if err != nil {
		err := err.(*types.BadSequence)
		t := types.Token{
			Type:    types.T_UNIT,
			Literal: err.Value.Value.(string),
		}
		return nil, NewCompileError(t, err.Value.TokenIndex, "unknown variable '%s'", t.Literal)
	}
	seq.AppendValue(types.Value{
		Type: types.V_FUNCBODY,
		Value: types.UserFunctionPart{
			Definition: tokens[fn.tokenPos+1:],
			Sequence:   bodySeq,
		},
		TokenIndex: -1,
	})
	return declStack, nil
}
