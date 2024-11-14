package functionimpl

import (
	"fmt"
	"math"
	"math/bits"
	"math/rand"

	"github.com/dece2183/hexowl/builtin/types"
	"github.com/dece2183/hexowl/utils"
)

func Sin(desc *types.Descriptor, args ...interface{}) (interface{}, error) {
	return math.Sin(utils.ToNumber[float64](args[0])), nil
}

func Cos(desc *types.Descriptor, args ...interface{}) (interface{}, error) {
	return math.Cos(utils.ToNumber[float64](args[0])), nil
}

func Asin(desc *types.Descriptor, args ...interface{}) (interface{}, error) {
	return math.Asin(utils.ToNumber[float64](args[0])), nil
}

func Acos(desc *types.Descriptor, args ...interface{}) (interface{}, error) {
	return math.Acos(utils.ToNumber[float64](args[0])), nil
}

func Tan(desc *types.Descriptor, args ...interface{}) (interface{}, error) {
	return math.Tan(utils.ToNumber[float64](args[0])), nil
}

func Atan(desc *types.Descriptor, args ...interface{}) (interface{}, error) {
	return math.Atan(utils.ToNumber[float64](args[0])), nil
}

func Pow(desc *types.Descriptor, args ...interface{}) (interface{}, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("not enough arguments")
	}
	return math.Pow(utils.ToNumber[float64](args[0]), utils.ToNumber[float64](args[1])), nil
}

func Sqrt(desc *types.Descriptor, args ...interface{}) (interface{}, error) {
	return math.Sqrt(utils.ToNumber[float64](args[0])), nil
}

func Logn(desc *types.Descriptor, args ...interface{}) (interface{}, error) {
	return math.Log(utils.ToNumber[float64](args[0])), nil
}

func Log2(desc *types.Descriptor, args ...interface{}) (interface{}, error) {
	return math.Log2(utils.ToNumber[float64](args[0])), nil
}

func Log10(desc *types.Descriptor, args ...interface{}) (interface{}, error) {
	return math.Log10(utils.ToNumber[float64](args[0])), nil
}

func Exp(desc *types.Descriptor, args ...interface{}) (interface{}, error) {
	return math.Exp(utils.ToNumber[float64](args[0])), nil
}

func Round(desc *types.Descriptor, args ...interface{}) (interface{}, error) {
	return math.Round(utils.ToNumber[float64](args[0])), nil
}

func Ceil(desc *types.Descriptor, args ...interface{}) (interface{}, error) {
	return math.Ceil(utils.ToNumber[float64](args[0])), nil
}

func Floor(desc *types.Descriptor, args ...interface{}) (interface{}, error) {
	return math.Floor(utils.ToNumber[float64](args[0])), nil
}

func Random(desc *types.Descriptor, args ...interface{}) (interface{}, error) {
	argslen := len(args)
	if argslen == 0 || args[0] == nil {
		return rand.Float64(), nil
	} else {
		if argslen == 1 {
			a := utils.ToNumber[int64](args[0])
			if a < 0 {
				return 0, fmt.Errorf("the first argument must be positive")
			}
			return rand.Int63n(a), nil
		} else {
			a := utils.ToNumber[int64](args[0])
			b := utils.ToNumber[int64](args[1])
			if b < a {
				return 0, fmt.Errorf("the first argument must be greater")
			}
			return rand.Int63n(b-a) + a, nil
		}
	}
}

func Popcount(desc *types.Descriptor, args ...interface{}) (interface{}, error) {
	return uint64(bits.OnesCount64(utils.ToNumber[uint64](args[0]))), nil
}
