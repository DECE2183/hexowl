package operators

import (
	"fmt"
	"math"
	"math/bits"

	"github.com/dece2183/hexowl/builtin"
	"github.com/dece2183/hexowl/user"
	"github.com/dece2183/hexowl/utils"
)

type action func(op *Operator, localVars map[string]interface{}) (interface{}, error)

var opActionListP *map[operatorType]action

var opActionList = map[operatorType]action{

	OP_NONE: func(op *Operator, localVars map[string]interface{}) (interface{}, error) {
		user.SetVariable(op.OperandA.Result.(string), op.Result)
		return op.Result, nil
	},

	OP_LOCALVAR: func(op *Operator, localVars map[string]interface{}) (interface{}, error) {
		localVars[op.OperandA.OperandA.Result.(string)] = op.Result
		return op.Result, nil
	},

	OP_USERVAR: func(op *Operator, localVars map[string]interface{}) (interface{}, error) {
		user.SetVariable(op.OperandA.OperandA.Result.(string), op.Result)
		return op.Result, nil
	},

	OP_SEQUENCE: func(op *Operator, localVars map[string]interface{}) (interface{}, error) {
		if op.OperandB != nil {
			op.Result = op.OperandB.Result
			return op.Result, nil
		}
		return nil, nil
	},

	OP_DECLFUNC: func(op *Operator, localVars map[string]interface{}) (interface{}, error) {
		leftSideWords := op.OperandA.Result.([]utils.Word)
		rightSideWords := op.OperandB.Result.([]utils.Word)
		funcName := leftSideWords[0].Literal
		newFunc := user.FuncVariant{
			Args: leftSideWords[2 : len(leftSideWords)-1],
			Body: rightSideWords,
		}
		user.SetFunctionVariant(funcName, newFunc)
		return op.Result, nil
	},

	OP_ASSIGN: func(op *Operator, localVars map[string]interface{}) (interface{}, error) {
		op.Result = op.OperandB.Result
		return opActionAssign(op, localVars)
	},

	OP_LOCALASSIGN: func(op *Operator, localVars map[string]interface{}) (interface{}, error) {
		op.Result = op.OperandB.Result
		localVars[op.OperandA.Result.(string)] = op.Result
		return op.Result, nil
	},

	OP_DECREMENT: func(op *Operator, localVars map[string]interface{}) (interface{}, error) {
		op.Result = utils.ToNumber[float64](op.OperandA.Result) - utils.ToNumber[float64](op.OperandB.Result)
		return opActionAssign(op, localVars)
	},

	OP_INCREMENT: func(op *Operator, localVars map[string]interface{}) (interface{}, error) {
		op.Result = utils.ToNumber[float64](op.OperandA.Result) + utils.ToNumber[float64](op.OperandB.Result)
		return opActionAssign(op, localVars)
	},

	OP_ASSIGNMUL: func(op *Operator, localVars map[string]interface{}) (interface{}, error) {
		op.Result = utils.ToNumber[float64](op.OperandA.Result) * utils.ToNumber[float64](op.OperandB.Result)
		return opActionAssign(op, localVars)
	},

	OP_ASSIGNDIV: func(op *Operator, localVars map[string]interface{}) (interface{}, error) {
		opB := utils.ToNumber[float64](op.OperandB.Result)
		if opB == 0 {
			op.Result = math.Inf(int(utils.ToNumber[float64](op.OperandA.Result)))
		} else {
			op.Result = utils.ToNumber[float64](op.OperandA.Result) / opB
		}
		return opActionAssign(op, localVars)
	},

	OP_ASSIGNBITAND: func(op *Operator, localVars map[string]interface{}) (interface{}, error) {
		op.Result = utils.ToNumber[uint64](op.OperandA.Result) & utils.ToNumber[uint64](op.OperandB.Result)
		return opActionAssign(op, localVars)
	},

	OP_ASSIGNBITOR: func(op *Operator, localVars map[string]interface{}) (interface{}, error) {
		op.Result = utils.ToNumber[uint64](op.OperandA.Result) | utils.ToNumber[uint64](op.OperandB.Result)
		return opActionAssign(op, localVars)
	},

	OP_LOGICNOT: func(op *Operator, localVars map[string]interface{}) (interface{}, error) {
		op.Result = !utils.ToBool(op.OperandB.Result)
		return op.Result, nil
	},

	OP_LOGICOR: func(op *Operator, localVars map[string]interface{}) (interface{}, error) {
		op.Result = utils.ToBool(op.OperandA.Result) || utils.ToBool(op.OperandB.Result)
		return op.Result, nil
	},

	OP_LOGICAND: func(op *Operator, localVars map[string]interface{}) (interface{}, error) {
		op.Result = utils.ToBool(op.OperandA.Result) && utils.ToBool(op.OperandB.Result)
		return op.Result, nil
	},

	OP_EQUALITY: func(op *Operator, localVars map[string]interface{}) (interface{}, error) {
		op.Result = utils.ToNumber[float64](op.OperandA.Result) == utils.ToNumber[float64](op.OperandB.Result)
		return op.Result, nil
	},

	OP_NOTEQ: func(op *Operator, localVars map[string]interface{}) (interface{}, error) {
		op.Result = utils.ToNumber[float64](op.OperandA.Result) != utils.ToNumber[float64](op.OperandB.Result)
		return op.Result, nil
	},

	OP_MORE: func(op *Operator, localVars map[string]interface{}) (interface{}, error) {
		op.Result = utils.ToNumber[float64](op.OperandA.Result) > utils.ToNumber[float64](op.OperandB.Result)
		return op.Result, nil
	},

	OP_LESS: func(op *Operator, localVars map[string]interface{}) (interface{}, error) {
		op.Result = utils.ToNumber[float64](op.OperandA.Result) < utils.ToNumber[float64](op.OperandB.Result)
		return op.Result, nil
	},

	OP_MOREEQ: func(op *Operator, localVars map[string]interface{}) (interface{}, error) {
		op.Result = utils.ToNumber[float64](op.OperandA.Result) >= utils.ToNumber[float64](op.OperandB.Result)
		return op.Result, nil
	},

	OP_LESSEQ: func(op *Operator, localVars map[string]interface{}) (interface{}, error) {
		op.Result = utils.ToNumber[float64](op.OperandA.Result) <= utils.ToNumber[float64](op.OperandB.Result)
		return op.Result, nil
	},

	OP_MINUS: func(op *Operator, localVars map[string]interface{}) (interface{}, error) {
		op.Result = utils.ToNumber[float64](op.OperandA.Result) - utils.ToNumber[float64](op.OperandB.Result)
		return op.Result, nil
	},

	OP_PLUS: func(op *Operator, localVars map[string]interface{}) (interface{}, error) {
		op.Result = utils.ToNumber[float64](op.OperandA.Result) + utils.ToNumber[float64](op.OperandB.Result)
		return op.Result, nil
	},

	OP_MULTIPLY: func(op *Operator, localVars map[string]interface{}) (interface{}, error) {
		op.Result = utils.ToNumber[float64](op.OperandA.Result) * utils.ToNumber[float64](op.OperandB.Result)
		return op.Result, nil
	},

	OP_DIVIDE: func(op *Operator, localVars map[string]interface{}) (interface{}, error) {
		opB := utils.ToNumber[float64](op.OperandB.Result)
		if opB == 0 {
			op.Result = math.Inf(int(utils.ToNumber[float64](op.OperandA.Result)))
		} else {
			op.Result = utils.ToNumber[float64](op.OperandA.Result) / opB
		}
		return op.Result, nil
	},

	OP_MODULO: func(op *Operator, localVars map[string]interface{}) (interface{}, error) {
		opB := utils.ToNumber[int64](op.OperandB.Result)
		if opB == 0 {
			op.Result = math.Inf(1)
		} else {
			op.Result = utils.ToNumber[int64](op.OperandA.Result) % opB
		}
		return op.Result, nil
	},

	OP_POWER: func(op *Operator, localVars map[string]interface{}) (interface{}, error) {
		op.Result = math.Pow(utils.ToNumber[float64](op.OperandA.Result), utils.ToNumber[float64](op.OperandB.Result))
		return op.Result, nil
	},

	OP_LEFTSHIFT: func(op *Operator, localVars map[string]interface{}) (interface{}, error) {
		op.Result = utils.ToNumber[uint64](op.OperandA.Result) << utils.ToNumber[uint64](op.OperandB.Result)
		return op.Result, nil
	},

	OP_RIGHTSHIFT: func(op *Operator, localVars map[string]interface{}) (interface{}, error) {
		op.Result = utils.ToNumber[uint64](op.OperandA.Result) >> utils.ToNumber[uint64](op.OperandB.Result)
		return op.Result, nil
	},

	OP_BITOR: func(op *Operator, localVars map[string]interface{}) (interface{}, error) {
		op.Result = utils.ToNumber[uint64](op.OperandA.Result) | utils.ToNumber[uint64](op.OperandB.Result)
		return op.Result, nil
	},

	OP_BITAND: func(op *Operator, localVars map[string]interface{}) (interface{}, error) {
		op.Result = utils.ToNumber[uint64](op.OperandA.Result) & utils.ToNumber[uint64](op.OperandB.Result)
		return op.Result, nil
	},

	OP_BITXOR: func(op *Operator, localVars map[string]interface{}) (interface{}, error) {
		op.Result = utils.ToNumber[uint64](op.OperandA.Result) ^ utils.ToNumber[uint64](op.OperandB.Result)
		return op.Result, nil
	},

	OP_BITCLEAR: func(op *Operator, localVars map[string]interface{}) (interface{}, error) {
		op.Result = utils.ToNumber[uint64](op.OperandA.Result) &^ utils.ToNumber[uint64](op.OperandB.Result)
		return op.Result, nil
	},

	OP_POPCNT: func(op *Operator, localVars map[string]interface{}) (interface{}, error) {
		op.Result = uint64(bits.OnesCount64(utils.ToNumber[uint64](op.OperandB.Result)))
		return op.Result, nil
	},

	OP_BITINVERSE: func(op *Operator, localVars map[string]interface{}) (interface{}, error) {
		op.Result = 0xFFFFFFFFFFFFFFFF ^ utils.ToNumber[uint64](op.OperandB.Result)
		return op.Result, nil
	},

	OP_ENUMERATE: func(op *Operator, localVars map[string]interface{}) (interface{}, error) {
		switch op.OperandA.Result.(type) {
		case []interface{}:
			break
		case nil:
			op.OperandA.Result = make([]interface{}, 0)
		default:
			op.OperandA.Result = []interface{}{op.OperandA.Result}
		}
		switch op.OperandB.Result.(type) {
		case []interface{}, nil:
			break
		default:
			op.OperandB.Result = []interface{}{op.OperandB.Result}
		}
		if op.OperandB.Result == nil {
			op.Result = op.OperandA.Result
		} else {
			op.Result = append(op.OperandA.Result.([]interface{}), op.OperandB.Result.([]interface{})...)
		}
		return op.Result, nil
	},

	OP_BUILTINFUNC: func(op *Operator, localVars map[string]interface{}) (interface{}, error) {
		var err error
		f, _ := builtin.GetFunction(op.OperandA.Result.(string))

		switch op.OperandB.Result.(type) {
		case []interface{}:
			op.Result, err = f.Exec(op.OperandB.Result.([]interface{})...)
		default:
			op.Result, err = f.Exec(op.OperandB.Result)
		}

		return op.Result, err
	},

	OP_USERFUNC: func(op *Operator, localVars map[string]interface{}) (interface{}, error) {
		var err error
		var args []interface{}
		fname := op.OperandA.Result.(string)
		f, _ := user.GetFunction(fname)

		switch op.OperandB.Result.(type) {
		case []interface{}:
			args = op.OperandB.Result.([]interface{})
		default:
			args = []interface{}{op.OperandB.Result}
		}

		op.Result, err = execUserFunc(f, args)
		if err != nil {
			return nil, fmt.Errorf("unable to find proper '%s' function variation for arguments: %v; (%s)", fname, args, err)
		}

		return op.Result, err
	},
}

func init() {
	opActionListP = &opActionList
}

func opActionAssign(op *Operator, localVars map[string]interface{}) (interface{}, error) {
	action, ok := (*opActionListP)[op.OperandA.Type]
	if ok {
		return action(op, localVars)
	} else {
		return nil, fmt.Errorf("try to assign non user variable")
	}
}

func opDoAction(op *Operator, localVars map[string]interface{}) (interface{}, error) {
	action, ok := (*opActionListP)[op.Type]
	if ok {
		return action(op, localVars)
	} else {
		return nil, fmt.Errorf("unable to find suitable action for operator type 0x%X", op.Type)
	}
}
