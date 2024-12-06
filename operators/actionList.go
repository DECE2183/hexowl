package operators

import (
	"fmt"
	"math"
	"math/bits"

	"github.com/dece2183/hexowl/builtin"
	"github.com/dece2183/hexowl/builtin/types"
	"github.com/dece2183/hexowl/user"
	"github.com/dece2183/hexowl/utils"
)

type actionHandler func(opLeft, opRight value, localVars map[string]interface{}) (interface{}, error)

var actionHandlerMap = map[operatorType]actionHandler{
	OP_NONE: func(opLeft, opRight value, localVars map[string]interface{}) (interface{}, error) {
		return nil, nil
	},

	OP_LOCALVAR: func(opLeft, opRight value, localVars map[string]interface{}) (interface{}, error) {
		localVars[opLeft.v.(string)] = opRight.v
		return opRight.v, nil
	},

	OP_USERVAR: func(opLeft, opRight value, localVars map[string]interface{}) (interface{}, error) {
		user.SetVariable(opLeft.v.(string), opRight.v)
		return opRight.v, nil
	},

	OP_SEQUENCE: func(opLeft, opRight value, localVars map[string]interface{}) (interface{}, error) {
		return actionObtainVariable(opRight, localVars)
	},

	OP_DECLFUNC: func(opLeft, opRight value, localVars map[string]interface{}) (interface{}, error) {
		leftSideWords := opLeft.v.([]utils.Word)
		rightSideWords := opRight.v.([]utils.Word)
		funcName := leftSideWords[0].Literal
		newFunc := user.FuncVariant{
			Args: leftSideWords[2 : len(leftSideWords)-1],
			Body: rightSideWords,
		}
		user.SetFunctionVariant(funcName, newFunc)
		//TODO: return func ptr special type
		return nil, nil
	},

	OP_ASSIGN: func(opLeft, opRight value, localVars map[string]interface{}) (interface{}, error) {
		val, err := actionObtainVariable(opRight, localVars)
		if err != nil {
			return nil, err
		}
		if opLeft.t == _VAL_CONST {
			opLeft.t = _VAL_USERVAR
		}
		return val, actionAssign(opLeft, val, localVars)
	},

	OP_ASSIGNLOCAL: func(opLeft, opRight value, localVars map[string]interface{}) (interface{}, error) {
		val, err := actionObtainVariable(opRight, localVars)
		if err != nil {
			return nil, err
		}
		opLeft.t = _VAL_LOCALVAR
		return val, actionAssign(opLeft, val, localVars)
	},

	OP_ASSIGNMINUS: func(opLeft, opRight value, localVars map[string]interface{}) (interface{}, error) {
		valLeft, err := actionObtainVariable(opLeft, localVars)
		if err != nil {
			return nil, err
		}
		valRight, err := actionObtainVariable(opRight, localVars)
		if err != nil {
			return nil, err
		}
		res := utils.ToNumber[float64](valLeft) - utils.ToNumber[float64](valRight)
		return res, actionAssign(opLeft, res, localVars)
	},

	OP_ASSIGNPLUS: func(opLeft, opRight value, localVars map[string]interface{}) (interface{}, error) {
		valLeft, err := actionObtainVariable(opLeft, localVars)
		if err != nil {
			return nil, err
		}
		valRight, err := actionObtainVariable(opRight, localVars)
		if err != nil {
			return nil, err
		}
		res := utils.ToNumber[float64](valLeft) + utils.ToNumber[float64](valRight)
		return res, actionAssign(opLeft, res, localVars)
	},

	OP_ASSIGNMUL: func(opLeft, opRight value, localVars map[string]interface{}) (interface{}, error) {
		valLeft, err := actionObtainVariable(opLeft, localVars)
		if err != nil {
			return nil, err
		}
		valRight, err := actionObtainVariable(opRight, localVars)
		if err != nil {
			return nil, err
		}
		res := utils.ToNumber[float64](valLeft) * utils.ToNumber[float64](valRight)
		return res, actionAssign(opLeft, res, localVars)
	},

	OP_ASSIGNDIV: func(opLeft, opRight value, localVars map[string]interface{}) (interface{}, error) {
		valLeft, err := actionObtainVariable(opLeft, localVars)
		if err != nil {
			return nil, err
		}
		valRight, err := actionObtainVariable(opRight, localVars)
		if err != nil {
			return nil, err
		}
		var res interface{}
		valLeftNum := utils.ToNumber[float64](valLeft)
		valRightNum := utils.ToNumber[float64](valRight)
		if valRightNum == 0 {
			res = math.Inf(int(valLeftNum))
		} else {
			res = valLeftNum / valRightNum
		}
		return res, actionAssign(opLeft, res, localVars)
	},

	OP_ASSIGNBITAND: func(opLeft, opRight value, localVars map[string]interface{}) (interface{}, error) {
		valLeft, err := actionObtainVariable(opLeft, localVars)
		if err != nil {
			return nil, err
		}
		valRight, err := actionObtainVariable(opRight, localVars)
		if err != nil {
			return nil, err
		}
		res := utils.ToNumber[uint64](valLeft) & utils.ToNumber[uint64](valRight)
		return res, actionAssign(opLeft, res, localVars)
	},

	OP_ASSIGNBITOR: func(opLeft, opRight value, localVars map[string]interface{}) (interface{}, error) {
		valLeft, err := actionObtainVariable(opLeft, localVars)
		if err != nil {
			return nil, err
		}
		valRight, err := actionObtainVariable(opRight, localVars)
		if err != nil {
			return nil, err
		}
		res := utils.ToNumber[uint64](valLeft) | utils.ToNumber[uint64](valRight)
		return res, actionAssign(opLeft, res, localVars)
	},

	OP_LOGICNOT: func(opLeft, opRight value, localVars map[string]interface{}) (interface{}, error) {
		valRight, err := actionObtainVariable(opRight, localVars)
		if err != nil {
			return nil, err
		}
		res := !utils.ToBool(valRight)
		return res, nil
	},

	OP_LOGICOR: func(opLeft, opRight value, localVars map[string]interface{}) (interface{}, error) {
		valLeft, err := actionObtainVariable(opLeft, localVars)
		if err != nil {
			return nil, err
		}
		valRight, err := actionObtainVariable(opRight, localVars)
		if err != nil {
			return nil, err
		}
		res := utils.ToBool(valLeft) || utils.ToBool(valRight)
		return res, nil
	},

	OP_LOGICAND: func(opLeft, opRight value, localVars map[string]interface{}) (interface{}, error) {
		valLeft, err := actionObtainVariable(opLeft, localVars)
		if err != nil {
			return nil, err
		}
		valRight, err := actionObtainVariable(opRight, localVars)
		if err != nil {
			return nil, err
		}
		res := utils.ToBool(valLeft) || utils.ToBool(valRight)
		return res, nil
	},

	OP_EQUALITY: func(opLeft, opRight value, localVars map[string]interface{}) (interface{}, error) {
		valLeft, err := actionObtainVariable(opLeft, localVars)
		if err != nil {
			return nil, err
		}
		valRight, err := actionObtainVariable(opRight, localVars)
		if err != nil {
			return nil, err
		}
		res := utils.ToNumber[float64](valLeft) == utils.ToNumber[float64](valRight)
		return res, nil
	},

	OP_NOTEQ: func(opLeft, opRight value, localVars map[string]interface{}) (interface{}, error) {
		valLeft, err := actionObtainVariable(opLeft, localVars)
		if err != nil {
			return nil, err
		}
		valRight, err := actionObtainVariable(opRight, localVars)
		if err != nil {
			return nil, err
		}
		res := utils.ToNumber[float64](valLeft) != utils.ToNumber[float64](valRight)
		return res, nil
	},

	OP_MORE: func(opLeft, opRight value, localVars map[string]interface{}) (interface{}, error) {
		valLeft, err := actionObtainVariable(opLeft, localVars)
		if err != nil {
			return nil, err
		}
		valRight, err := actionObtainVariable(opRight, localVars)
		if err != nil {
			return nil, err
		}
		res := utils.ToNumber[float64](valLeft) > utils.ToNumber[float64](valRight)
		return res, nil
	},

	OP_LESS: func(opLeft, opRight value, localVars map[string]interface{}) (interface{}, error) {
		valLeft, err := actionObtainVariable(opLeft, localVars)
		if err != nil {
			return nil, err
		}
		valRight, err := actionObtainVariable(opRight, localVars)
		if err != nil {
			return nil, err
		}
		res := utils.ToNumber[float64](valLeft) < utils.ToNumber[float64](valRight)
		return res, nil
	},

	OP_MOREEQ: func(opLeft, opRight value, localVars map[string]interface{}) (interface{}, error) {
		valLeft, err := actionObtainVariable(opLeft, localVars)
		if err != nil {
			return nil, err
		}
		valRight, err := actionObtainVariable(opRight, localVars)
		if err != nil {
			return nil, err
		}
		res := utils.ToNumber[float64](valLeft) >= utils.ToNumber[float64](valRight)
		return res, nil
	},

	OP_LESSEQ: func(opLeft, opRight value, localVars map[string]interface{}) (interface{}, error) {
		valLeft, err := actionObtainVariable(opLeft, localVars)
		if err != nil {
			return nil, err
		}
		valRight, err := actionObtainVariable(opRight, localVars)
		if err != nil {
			return nil, err
		}
		res := utils.ToNumber[float64](valLeft) <= utils.ToNumber[float64](valRight)
		return res, nil
	},

	OP_MINUS: func(opLeft, opRight value, localVars map[string]interface{}) (interface{}, error) {
		valLeft, err := actionObtainVariable(opLeft, localVars)
		if err != nil {
			return nil, err
		}
		valRight, err := actionObtainVariable(opRight, localVars)
		if err != nil {
			return nil, err
		}
		res := utils.ToNumber[float64](valLeft) - utils.ToNumber[float64](valRight)
		return res, nil
	},

	OP_PLUS: func(opLeft, opRight value, localVars map[string]interface{}) (interface{}, error) {
		valLeft, err := actionObtainVariable(opLeft, localVars)
		if err != nil {
			return nil, err
		}
		valRight, err := actionObtainVariable(opRight, localVars)
		if err != nil {
			return nil, err
		}
		res := utils.ToNumber[float64](valLeft) + utils.ToNumber[float64](valRight)
		return res, nil
	},

	OP_MULTIPLY: func(opLeft, opRight value, localVars map[string]interface{}) (interface{}, error) {
		valLeft, err := actionObtainVariable(opLeft, localVars)
		if err != nil {
			return nil, err
		}
		valRight, err := actionObtainVariable(opRight, localVars)
		if err != nil {
			return nil, err
		}
		res := utils.ToNumber[float64](valLeft) * utils.ToNumber[float64](valRight)
		return res, nil
	},

	OP_DIVIDE: func(opLeft, opRight value, localVars map[string]interface{}) (interface{}, error) {
		valLeft, err := actionObtainVariable(opLeft, localVars)
		if err != nil {
			return nil, err
		}
		valRight, err := actionObtainVariable(opRight, localVars)
		if err != nil {
			return nil, err
		}
		var res interface{}
		valLeftNum := utils.ToNumber[float64](valLeft)
		valRightNum := utils.ToNumber[float64](valRight)
		if valRightNum == 0 {
			res = math.Inf(int(valLeftNum))
		} else {
			res = valLeftNum / valRightNum
		}
		return res, nil
	},

	OP_MODULO: func(opLeft, opRight value, localVars map[string]interface{}) (interface{}, error) {
		valLeft, err := actionObtainVariable(opLeft, localVars)
		if err != nil {
			return nil, err
		}
		valRight, err := actionObtainVariable(opRight, localVars)
		if err != nil {
			return nil, err
		}
		var res interface{}
		valLeftNum := utils.ToNumber[int64](valLeft)
		valRightNum := utils.ToNumber[int64](valRight)
		if valRightNum == 0 {
			res = math.Inf(1)
		} else {
			res = valLeftNum % valRightNum
		}
		return res, nil
	},

	OP_POWER: func(opLeft, opRight value, localVars map[string]interface{}) (interface{}, error) {
		valLeft, err := actionObtainVariable(opLeft, localVars)
		if err != nil {
			return nil, err
		}
		valRight, err := actionObtainVariable(opRight, localVars)
		if err != nil {
			return nil, err
		}
		res := math.Pow(utils.ToNumber[float64](valLeft), utils.ToNumber[float64](valRight))
		return res, nil
	},

	OP_LEFTSHIFT: func(opLeft, opRight value, localVars map[string]interface{}) (interface{}, error) {
		valLeft, err := actionObtainVariable(opLeft, localVars)
		if err != nil {
			return nil, err
		}
		valRight, err := actionObtainVariable(opRight, localVars)
		if err != nil {
			return nil, err
		}
		res := utils.ToNumber[uint64](valLeft) << utils.ToNumber[uint64](valRight)
		return res, nil
	},

	OP_RIGHTSHIFT: func(opLeft, opRight value, localVars map[string]interface{}) (interface{}, error) {
		valLeft, err := actionObtainVariable(opLeft, localVars)
		if err != nil {
			return nil, err
		}
		valRight, err := actionObtainVariable(opRight, localVars)
		if err != nil {
			return nil, err
		}
		res := utils.ToNumber[uint64](valLeft) >> utils.ToNumber[uint64](valRight)
		return res, nil
	},

	OP_BITOR: func(opLeft, opRight value, localVars map[string]interface{}) (interface{}, error) {
		valLeft, err := actionObtainVariable(opLeft, localVars)
		if err != nil {
			return nil, err
		}
		valRight, err := actionObtainVariable(opRight, localVars)
		if err != nil {
			return nil, err
		}
		res := utils.ToNumber[uint64](valLeft) | utils.ToNumber[uint64](valRight)
		return res, nil
	},

	OP_BITAND: func(opLeft, opRight value, localVars map[string]interface{}) (interface{}, error) {
		valLeft, err := actionObtainVariable(opLeft, localVars)
		if err != nil {
			return nil, err
		}
		valRight, err := actionObtainVariable(opRight, localVars)
		if err != nil {
			return nil, err
		}
		res := utils.ToNumber[uint64](valLeft) & utils.ToNumber[uint64](valRight)
		return res, nil
	},

	OP_BITXOR: func(opLeft, opRight value, localVars map[string]interface{}) (interface{}, error) {
		valLeft, err := actionObtainVariable(opLeft, localVars)
		if err != nil {
			return nil, err
		}
		valRight, err := actionObtainVariable(opRight, localVars)
		if err != nil {
			return nil, err
		}
		res := utils.ToNumber[uint64](valLeft) ^ utils.ToNumber[uint64](valRight)
		return res, nil
	},

	OP_BITCLEAR: func(opLeft, opRight value, localVars map[string]interface{}) (interface{}, error) {
		valLeft, err := actionObtainVariable(opLeft, localVars)
		if err != nil {
			return nil, err
		}
		valRight, err := actionObtainVariable(opRight, localVars)
		if err != nil {
			return nil, err
		}
		res := utils.ToNumber[uint64](valLeft) &^ utils.ToNumber[uint64](valRight)
		return res, nil
	},

	OP_POPCNT: func(opLeft, opRight value, localVars map[string]interface{}) (interface{}, error) {
		valRight, err := actionObtainVariable(opRight, localVars)
		if err != nil {
			return nil, err
		}
		res := uint64(bits.OnesCount64(utils.ToNumber[uint64](valRight)))
		return res, nil
	},

	OP_BITINVERSE: func(opLeft, opRight value, localVars map[string]interface{}) (interface{}, error) {
		valRight, err := actionObtainVariable(opRight, localVars)
		if err != nil {
			return nil, err
		}
		res := 0xFFFFFFFFFFFFFFFF ^ utils.ToNumber[uint64](valRight)
		return res, nil
	},

	OP_ENUMERATE: func(opLeft, opRight value, localVars map[string]interface{}) (interface{}, error) {
		valLeft, err := actionObtainVariable(opLeft, localVars)
		if err != nil {
			return nil, err
		}
		valRight, err := actionObtainVariable(opRight, localVars)
		if err != nil {
			return nil, err
		}

		switch valLeft.(type) {
		case []interface{}:
			break
		case nil:
			valLeft = make([]interface{}, 0)
		default:
			valLeft = []interface{}{valLeft}
		}

		switch valRight.(type) {
		case []interface{}:
			break
		case nil:
			valRight = []interface{}{}
		default:
			valRight = []interface{}{valRight}
		}

		res := append(valLeft.([]interface{}), valRight.([]interface{})...)
		return res, nil
	},

	OP_CALLFUNC: func(opLeft, opRight value, localVars map[string]interface{}) (interface{}, error) {
		fn, err := actionObtainVariable(opLeft, localVars)
		if err != nil {
			return nil, err
		}
		args, err := actionObtainVariable(opRight, localVars)
		if err != nil {
			return nil, err
		}

		if _, ok := args.([]interface{}); !ok {
			args = []interface{}{args}
		}

		switch fn := fn.(type) {
		case user.Func:
			return execUserFunc(fn, args.([]interface{}))
		case types.Func:
			return builtin.Exec(fn, args.([]interface{})...)
		}

		return nil, fmt.Errorf("'%v' is not a function", opLeft.v)
	},

	// OP_BUILTINFUNC: func(opLeft, opRight value, localVars map[string]interface{}) (interface{}, error) {
	// 	var err error
	// 	f, _ := builtin.GetFunction(op.OperandA.Result.(string))

	// 	switch op.OperandB.Result.(type) {
	// 	case []interface{}:
	// 		op.Result, err = builtin.Exec(f, op.OperandB.Result.([]interface{})...)
	// 	default:
	// 		op.Result, err = builtin.Exec(f, op.OperandB.Result)
	// 	}

	// 	return op.Result, err
	// },

	// OP_USERFUNC: func(opLeft, opRight value, localVars map[string]interface{}) (interface{}, error) {
	// 	var err error
	// 	var args []interface{}
	// 	fname := op.OperandA.Result.(string)
	// 	f, _ := user.GetFunction(fname)

	// 	switch op.OperandB.Result.(type) {
	// 	case []interface{}:
	// 		args = op.OperandB.Result.([]interface{})
	// 	default:
	// 		args = []interface{}{op.OperandB.Result}
	// 	}

	// 	op.Result, err = execUserFunc(f, args)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("unable to find proper '%s' function variation for arguments: %v; (%s)", fname, args, err)
	// 	}

	// 	return op.Result, err
	// },
}

func actionAssign(variable value, val interface{}, localVars map[string]interface{}) error {
	switch variable.t {
	case _VAL_LOCALVAR:
		localVars[variable.v.(string)] = val
	case _VAL_USERVAR:
		user.SetVariable(variable.v.(string), val)
	default:
		return fmt.Errorf("%s '%s' not allowed for assignment", variable.v.(string), valStringRepresent[variable.t])
	}
	return nil
}

func actionObtainVariable(variable value, localVars map[string]interface{}) (interface{}, error) {
	if variable.t == _VAL_CONST {
		return variable.v, nil
	}

	var (
		val interface{}
		ok  bool
	)

	varName := variable.v.(string)

	switch variable.t {
	case _VAL_LOCALVAR:
		if val, ok = localVars[varName]; !ok {
			return nil, fmt.Errorf("'%s' is not a local variable", varName)
		}
	case _VAL_USERVAR:
		if val, ok = user.GetVariable(varName); !ok {
			return nil, fmt.Errorf("'%s' is not a user variable", varName)
		}
	case _VAL_CONSTANT:
		if val, ok = builtin.GetConstant(varName); !ok {
			return nil, fmt.Errorf("'%s' is not a built-in constant", varName)
		}
	case _VAL_USERFUNC:
		if val, ok = user.GetFunction(varName); !ok {
			return nil, fmt.Errorf("'%s' is not a user function", varName)
		}
	case _VAL_BUILTINFUNC:
		if val, ok = builtin.GetFunction(varName); !ok {
			return nil, fmt.Errorf("'%s' is not a built-in function", varName)
		}
	case _VAL_LOCALFUNCPTR:
	case _VAL_FUNCPTR:
	}

	return val, nil
}
