package runtime

import (
	"fmt"
	"math"
	"math/bits"

	"github.com/dece2183/hexowl/v2/calculator/types"
	"github.com/dece2183/hexowl/v2/utils"
)

type actionHandler func(rn *Runtime, opLeft, opRight types.Value) (interface{}, error)

var actionHandlerMap = map[types.OperatorType]actionHandler{
	types.O_NONE: func(rn *Runtime, opLeft, opRight types.Value) (interface{}, error) {
		return nil, nil
	},

	types.O_SEQUENCE: func(rn *Runtime, opLeft, opRight types.Value) (interface{}, error) {
		return rn.obtainVariable(opRight)
	},

	types.O_DECLFUNC: func(rn *Runtime, opLeft, opRight types.Value) (interface{}, error) {
		//TODO: implement
		//TODO: return func ptr special type
		return nil, nil
	},

	types.O_ASSIGN: func(rn *Runtime, opLeft, opRight types.Value) (interface{}, error) {
		val, err := rn.obtainVariable(opRight)
		if err != nil {
			return nil, err
		}
		if opLeft.Type != types.V_LOCALVAR {
			opLeft.Type = types.V_USERVAR
		}
		return val, rn.assignValue(opLeft, val)
	},

	types.O_ASSIGNLOCAL: func(rn *Runtime, opLeft, opRight types.Value) (interface{}, error) {
		val, err := rn.obtainVariable(opRight)
		if err != nil {
			return nil, err
		}
		opLeft.Type = types.V_LOCALVAR
		return val, rn.assignValue(opLeft, val)
	},

	types.O_ASSIGNMINUS: func(rn *Runtime, opLeft, opRight types.Value) (interface{}, error) {
		valLeft, err := rn.obtainVariable(opLeft)
		if err != nil {
			return nil, err
		}
		valRight, err := rn.obtainVariable(opRight)
		if err != nil {
			return nil, err
		}
		res := utils.ToNumber[float64](valLeft) - utils.ToNumber[float64](valRight)
		return res, rn.assignValue(opLeft, res)
	},

	types.O_ASSIGNPLUS: func(rn *Runtime, opLeft, opRight types.Value) (interface{}, error) {
		valLeft, err := rn.obtainVariable(opLeft)
		if err != nil {
			return nil, err
		}
		valRight, err := rn.obtainVariable(opRight)
		if err != nil {
			return nil, err
		}
		res := utils.ToNumber[float64](valLeft) + utils.ToNumber[float64](valRight)
		return res, rn.assignValue(opLeft, res)
	},

	types.O_ASSIGNMUL: func(rn *Runtime, opLeft, opRight types.Value) (interface{}, error) {
		valLeft, err := rn.obtainVariable(opLeft)
		if err != nil {
			return nil, err
		}
		valRight, err := rn.obtainVariable(opRight)
		if err != nil {
			return nil, err
		}
		res := utils.ToNumber[float64](valLeft) * utils.ToNumber[float64](valRight)
		return res, rn.assignValue(opLeft, res)
	},

	types.O_ASSIGNDIV: func(rn *Runtime, opLeft, opRight types.Value) (interface{}, error) {
		valLeft, err := rn.obtainVariable(opLeft)
		if err != nil {
			return nil, err
		}
		valRight, err := rn.obtainVariable(opRight)
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
		return res, rn.assignValue(opLeft, res)
	},

	types.O_ASSIGNBITAND: func(rn *Runtime, opLeft, opRight types.Value) (interface{}, error) {
		valLeft, err := rn.obtainVariable(opLeft)
		if err != nil {
			return nil, err
		}
		valRight, err := rn.obtainVariable(opRight)
		if err != nil {
			return nil, err
		}
		res := utils.ToNumber[uint64](valLeft) & utils.ToNumber[uint64](valRight)
		return res, rn.assignValue(opLeft, res)
	},

	types.O_ASSIGNBITOR: func(rn *Runtime, opLeft, opRight types.Value) (interface{}, error) {
		valLeft, err := rn.obtainVariable(opLeft)
		if err != nil {
			return nil, err
		}
		valRight, err := rn.obtainVariable(opRight)
		if err != nil {
			return nil, err
		}
		res := utils.ToNumber[uint64](valLeft) | utils.ToNumber[uint64](valRight)
		return res, rn.assignValue(opLeft, res)
	},

	types.O_LOGICNOT: func(rn *Runtime, opLeft, opRight types.Value) (interface{}, error) {
		valRight, err := rn.obtainVariable(opRight)
		if err != nil {
			return nil, err
		}
		res := !utils.ToBool(valRight)
		return res, nil
	},

	types.O_LOGICOR: func(rn *Runtime, opLeft, opRight types.Value) (interface{}, error) {
		valLeft, err := rn.obtainVariable(opLeft)
		if err != nil {
			return nil, err
		}
		valRight, err := rn.obtainVariable(opRight)
		if err != nil {
			return nil, err
		}
		res := utils.ToBool(valLeft) || utils.ToBool(valRight)
		return res, nil
	},

	types.O_LOGICAND: func(rn *Runtime, opLeft, opRight types.Value) (interface{}, error) {
		valLeft, err := rn.obtainVariable(opLeft)
		if err != nil {
			return nil, err
		}
		valRight, err := rn.obtainVariable(opRight)
		if err != nil {
			return nil, err
		}
		res := utils.ToBool(valLeft) || utils.ToBool(valRight)
		return res, nil
	},

	types.O_EQUALITY: func(rn *Runtime, opLeft, opRight types.Value) (interface{}, error) {
		valLeft, err := rn.obtainVariable(opLeft)
		if err != nil {
			return nil, err
		}
		valRight, err := rn.obtainVariable(opRight)
		if err != nil {
			return nil, err
		}
		res := utils.ToNumber[float64](valLeft) == utils.ToNumber[float64](valRight)
		return res, nil
	},

	types.O_NOTEQ: func(rn *Runtime, opLeft, opRight types.Value) (interface{}, error) {
		valLeft, err := rn.obtainVariable(opLeft)
		if err != nil {
			return nil, err
		}
		valRight, err := rn.obtainVariable(opRight)
		if err != nil {
			return nil, err
		}
		res := utils.ToNumber[float64](valLeft) != utils.ToNumber[float64](valRight)
		return res, nil
	},

	types.O_MORE: func(rn *Runtime, opLeft, opRight types.Value) (interface{}, error) {
		valLeft, err := rn.obtainVariable(opLeft)
		if err != nil {
			return nil, err
		}
		valRight, err := rn.obtainVariable(opRight)
		if err != nil {
			return nil, err
		}
		res := utils.ToNumber[float64](valLeft) > utils.ToNumber[float64](valRight)
		return res, nil
	},

	types.O_LESS: func(rn *Runtime, opLeft, opRight types.Value) (interface{}, error) {
		valLeft, err := rn.obtainVariable(opLeft)
		if err != nil {
			return nil, err
		}
		valRight, err := rn.obtainVariable(opRight)
		if err != nil {
			return nil, err
		}
		res := utils.ToNumber[float64](valLeft) < utils.ToNumber[float64](valRight)
		return res, nil
	},

	types.O_MOREEQ: func(rn *Runtime, opLeft, opRight types.Value) (interface{}, error) {
		valLeft, err := rn.obtainVariable(opLeft)
		if err != nil {
			return nil, err
		}
		valRight, err := rn.obtainVariable(opRight)
		if err != nil {
			return nil, err
		}
		res := utils.ToNumber[float64](valLeft) >= utils.ToNumber[float64](valRight)
		return res, nil
	},

	types.O_LESSEQ: func(rn *Runtime, opLeft, opRight types.Value) (interface{}, error) {
		valLeft, err := rn.obtainVariable(opLeft)
		if err != nil {
			return nil, err
		}
		valRight, err := rn.obtainVariable(opRight)
		if err != nil {
			return nil, err
		}
		res := utils.ToNumber[float64](valLeft) <= utils.ToNumber[float64](valRight)
		return res, nil
	},

	types.O_MINUS: func(rn *Runtime, opLeft, opRight types.Value) (interface{}, error) {
		valLeft, err := rn.obtainVariable(opLeft)
		if err != nil {
			return nil, err
		}
		valRight, err := rn.obtainVariable(opRight)
		if err != nil {
			return nil, err
		}
		res := utils.ToNumber[float64](valLeft) - utils.ToNumber[float64](valRight)
		return res, nil
	},

	types.O_PLUS: func(rn *Runtime, opLeft, opRight types.Value) (interface{}, error) {
		valLeft, err := rn.obtainVariable(opLeft)
		if err != nil {
			return nil, err
		}
		valRight, err := rn.obtainVariable(opRight)
		if err != nil {
			return nil, err
		}
		res := utils.ToNumber[float64](valLeft) + utils.ToNumber[float64](valRight)
		return res, nil
	},

	types.O_MULTIPLY: func(rn *Runtime, opLeft, opRight types.Value) (interface{}, error) {
		valLeft, err := rn.obtainVariable(opLeft)
		if err != nil {
			return nil, err
		}
		valRight, err := rn.obtainVariable(opRight)
		if err != nil {
			return nil, err
		}
		res := utils.ToNumber[float64](valLeft) * utils.ToNumber[float64](valRight)
		return res, nil
	},

	types.O_DIVIDE: func(rn *Runtime, opLeft, opRight types.Value) (interface{}, error) {
		valLeft, err := rn.obtainVariable(opLeft)
		if err != nil {
			return nil, err
		}
		valRight, err := rn.obtainVariable(opRight)
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

	types.O_MODULO: func(rn *Runtime, opLeft, opRight types.Value) (interface{}, error) {
		valLeft, err := rn.obtainVariable(opLeft)
		if err != nil {
			return nil, err
		}
		valRight, err := rn.obtainVariable(opRight)
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

	types.O_POWER: func(rn *Runtime, opLeft, opRight types.Value) (interface{}, error) {
		valLeft, err := rn.obtainVariable(opLeft)
		if err != nil {
			return nil, err
		}
		valRight, err := rn.obtainVariable(opRight)
		if err != nil {
			return nil, err
		}
		res := math.Pow(utils.ToNumber[float64](valLeft), utils.ToNumber[float64](valRight))
		return res, nil
	},

	types.O_LEFTSHIFT: func(rn *Runtime, opLeft, opRight types.Value) (interface{}, error) {
		valLeft, err := rn.obtainVariable(opLeft)
		if err != nil {
			return nil, err
		}
		valRight, err := rn.obtainVariable(opRight)
		if err != nil {
			return nil, err
		}
		res := utils.ToNumber[uint64](valLeft) << utils.ToNumber[uint64](valRight)
		return res, nil
	},

	types.O_RIGHTSHIFT: func(rn *Runtime, opLeft, opRight types.Value) (interface{}, error) {
		valLeft, err := rn.obtainVariable(opLeft)
		if err != nil {
			return nil, err
		}
		valRight, err := rn.obtainVariable(opRight)
		if err != nil {
			return nil, err
		}
		res := utils.ToNumber[uint64](valLeft) >> utils.ToNumber[uint64](valRight)
		return res, nil
	},

	types.O_BITOR: func(rn *Runtime, opLeft, opRight types.Value) (interface{}, error) {
		valLeft, err := rn.obtainVariable(opLeft)
		if err != nil {
			return nil, err
		}
		valRight, err := rn.obtainVariable(opRight)
		if err != nil {
			return nil, err
		}
		res := utils.ToNumber[uint64](valLeft) | utils.ToNumber[uint64](valRight)
		return res, nil
	},

	types.O_BITAND: func(rn *Runtime, opLeft, opRight types.Value) (interface{}, error) {
		valLeft, err := rn.obtainVariable(opLeft)
		if err != nil {
			return nil, err
		}
		valRight, err := rn.obtainVariable(opRight)
		if err != nil {
			return nil, err
		}
		res := utils.ToNumber[uint64](valLeft) & utils.ToNumber[uint64](valRight)
		return res, nil
	},

	types.O_BITXOR: func(rn *Runtime, opLeft, opRight types.Value) (interface{}, error) {
		valLeft, err := rn.obtainVariable(opLeft)
		if err != nil {
			return nil, err
		}
		valRight, err := rn.obtainVariable(opRight)
		if err != nil {
			return nil, err
		}
		res := utils.ToNumber[uint64](valLeft) ^ utils.ToNumber[uint64](valRight)
		return res, nil
	},

	types.O_BITCLEAR: func(rn *Runtime, opLeft, opRight types.Value) (interface{}, error) {
		valLeft, err := rn.obtainVariable(opLeft)
		if err != nil {
			return nil, err
		}
		valRight, err := rn.obtainVariable(opRight)
		if err != nil {
			return nil, err
		}
		res := utils.ToNumber[uint64](valLeft) &^ utils.ToNumber[uint64](valRight)
		return res, nil
	},

	types.O_POPCNT: func(rn *Runtime, opLeft, opRight types.Value) (interface{}, error) {
		valRight, err := rn.obtainVariable(opRight)
		if err != nil {
			return nil, err
		}
		res := uint64(bits.OnesCount64(utils.ToNumber[uint64](valRight)))
		return res, nil
	},

	types.O_BITINVERSE: func(rn *Runtime, opLeft, opRight types.Value) (interface{}, error) {
		valRight, err := rn.obtainVariable(opRight)
		if err != nil {
			return nil, err
		}
		res := 0xFFFFFFFFFFFFFFFF ^ utils.ToNumber[uint64](valRight)
		return res, nil
	},

	types.O_ENUMERATE: func(rn *Runtime, opLeft, opRight types.Value) (interface{}, error) {
		valLeft, err := rn.obtainVariable(opLeft)
		if err != nil {
			return nil, err
		}
		valRight, err := rn.obtainVariable(opRight)
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

	types.O_CALLFUNC: func(rn *Runtime, opLeft, opRight types.Value) (interface{}, error) {
		fn, err := rn.obtainVariable(opLeft)
		if err != nil {
			return nil, err
		}
		args, err := rn.obtainVariable(opRight)
		if err != nil {
			return nil, err
		}

		if _, ok := args.([]interface{}); !ok {
			args = []interface{}{args}
		}

		switch fn := fn.(type) {
		case types.UserFunction:
			_ = fn
		case types.BuiltinFunction:
			return fn.Exec(rn.ctx, args.([]interface{}))
		}

		return nil, fmt.Errorf("'%v' is not a function", opLeft.Value)
	},

	// types.O_BUILTINFUNC: func(rn *Runtime, opLeft, opRight types.Value) (interface{}, error) {
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

	// types.O_USERFUNC: func(rn *Runtime, opLeft, opRight types.Value) (interface{}, error) {
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
