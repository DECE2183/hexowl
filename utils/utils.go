package utils

import "math"

type Number interface {
	int64 | uint64 | float64
}

func ToNumber[T Number](i interface{}) T {
	switch v := i.(type) {
	case bool:
		if v {
			return T(1)
		} else {
			return T(0)
		}
	case string:
		arr := []byte(v)
		var val uint64
		for i := 0; i < int(math.Min(float64(len(arr)), 8)); i++ {
			val |= (uint64(arr[i]) << (i * 8))
		}
		return T(val)
	case byte:
		return T(v)
	case int:
		return T(v)
	case uint:
		return T(v)
	case int64:
		return T(v)
	case uint64:
		return T(v)
	case float32:
		return T(v)
	case float64:
		return T(v)
	}

	return T(0)
}

func ToBool(i interface{}) bool {
	switch v := i.(type) {
	case bool:
		return v
	case int64:
		if v > 0 {
			return true
		} else {
			return false
		}
	case uint64:
		if v > 0 {
			return true
		} else {
			return false
		}
	case float64:
		if v > 0 {
			return true
		} else {
			return false
		}
	}

	return false
}
