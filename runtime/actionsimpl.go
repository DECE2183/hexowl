package runtime

import (
	"fmt"
	"math"
	"math/bits"

	"github.com/dece2183/hexowl/v2/types"
	"github.com/dece2183/hexowl/v2/utils"
)

func implNONE(rn *Runtime, opLeft, opRight types.Value) (interface{}, error) {
	return nil, nil
}

func implSEQUENCE(rn *Runtime, opLeft, opRight types.Value) (interface{}, error) {
	return rn.obtainVariable(opRight)
}

func implDECLFUNC(rn *Runtime, opLeft, opRight types.Value) (interface{}, error) {
	args := opLeft.Value.(*types.ExecutionSequence)
	funcName, _ := args.GetValue(0)
	name := funcName.Value.(string)
	args = args.ExtractSubsequence(1, args.Len()-1)
	body := opRight.Value.(*types.ExecutionSequence)
	rn.ctx.User.SetFunctionVariant(name, types.UserFunctionVariant{
		ArgsSequence: args,
		BodySequence: body,
	})
	userFunc, _ := rn.ctx.User.GetFunction(name)
	return userFunc, nil
}

func implASSIGN(rn *Runtime, opLeft, opRight types.Value) (interface{}, error) {
	val, err := rn.obtainVariable(opRight)
	if err != nil {
		return nil, err
	}
	return val, rn.assignValue(opLeft, val)
}

func implASSIGNLOCAL(rn *Runtime, opLeft, opRight types.Value) (interface{}, error) {
	val, err := rn.obtainVariable(opRight)
	if err != nil {
		return nil, err
	}
	return val, rn.assignLocalValue(opLeft, val)
}

func implASSIGNMINUS(rn *Runtime, opLeft, opRight types.Value) (interface{}, error) {
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
}

func implASSIGNPLUS(rn *Runtime, opLeft, opRight types.Value) (interface{}, error) {
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
}

func implASSIGNMUL(rn *Runtime, opLeft, opRight types.Value) (interface{}, error) {
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
}

func implASSIGNDIV(rn *Runtime, opLeft, opRight types.Value) (interface{}, error) {
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
}

func implASSIGNBITAND(rn *Runtime, opLeft, opRight types.Value) (interface{}, error) {
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
}

func implASSIGNBITOR(rn *Runtime, opLeft, opRight types.Value) (interface{}, error) {
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
}

func implLOGICNOT(rn *Runtime, opLeft, opRight types.Value) (interface{}, error) {
	valRight, err := rn.obtainVariable(opRight)
	if err != nil {
		return nil, err
	}
	res := !utils.ToBool(valRight)
	return res, nil
}

func implLOGICOR(rn *Runtime, opLeft, opRight types.Value) (interface{}, error) {
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
}

func implLOGICAND(rn *Runtime, opLeft, opRight types.Value) (interface{}, error) {
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
}

func implEQUALITY(rn *Runtime, opLeft, opRight types.Value) (interface{}, error) {
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
}

func implNOTEQ(rn *Runtime, opLeft, opRight types.Value) (interface{}, error) {
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
}

func implMORE(rn *Runtime, opLeft, opRight types.Value) (interface{}, error) {
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
}

func implLESS(rn *Runtime, opLeft, opRight types.Value) (interface{}, error) {
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
}

func implMOREEQ(rn *Runtime, opLeft, opRight types.Value) (interface{}, error) {
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
}

func implLESSEQ(rn *Runtime, opLeft, opRight types.Value) (interface{}, error) {
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
}

func implMINUS(rn *Runtime, opLeft, opRight types.Value) (interface{}, error) {
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
}

func implPLUS(rn *Runtime, opLeft, opRight types.Value) (interface{}, error) {
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
}

func implMULTIPLY(rn *Runtime, opLeft, opRight types.Value) (interface{}, error) {
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
}

func implDIVIDE(rn *Runtime, opLeft, opRight types.Value) (interface{}, error) {
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
}

func implMODULO(rn *Runtime, opLeft, opRight types.Value) (interface{}, error) {
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
}

func implPOWER(rn *Runtime, opLeft, opRight types.Value) (interface{}, error) {
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
}

func implLEFTSHIFT(rn *Runtime, opLeft, opRight types.Value) (interface{}, error) {
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
}

func implRIGHTSHIFT(rn *Runtime, opLeft, opRight types.Value) (interface{}, error) {
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
}

func implBITOR(rn *Runtime, opLeft, opRight types.Value) (interface{}, error) {
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
}

func implBITAND(rn *Runtime, opLeft, opRight types.Value) (interface{}, error) {
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
}

func implBITXOR(rn *Runtime, opLeft, opRight types.Value) (interface{}, error) {
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
}

func implBITCLEAR(rn *Runtime, opLeft, opRight types.Value) (interface{}, error) {
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
}

func implPOPCNT(rn *Runtime, opLeft, opRight types.Value) (interface{}, error) {
	valRight, err := rn.obtainVariable(opRight)
	if err != nil {
		return nil, err
	}
	res := uint64(bits.OnesCount64(utils.ToNumber[uint64](valRight)))
	return res, nil
}

func implBITINVERSE(rn *Runtime, opLeft, opRight types.Value) (interface{}, error) {
	valRight, err := rn.obtainVariable(opRight)
	if err != nil {
		return nil, err
	}
	res := 0xFFFFFFFFFFFFFFFF ^ utils.ToNumber[uint64](valRight)
	return res, nil
}

func implENUMERATE(rn *Runtime, opLeft, opRight types.Value) (interface{}, error) {
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
}

func implCALLFUNC(rn *Runtime, opLeft, opRight types.Value) (interface{}, error) {
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
		return rn.ExecuteUserFunction(fn, args.([]interface{}))
	case types.BuiltinFunction:
		return fn.Exec(rn.ctx, args.([]interface{}))
	}

	return nil, fmt.Errorf("'%v' is not a function", opLeft.Value)
}
