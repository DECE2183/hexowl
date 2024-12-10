package utils

import (
	"bytes"
	"encoding/binary"
	"math"
)

type number interface {
	int64 | uint64 | float64
}

// Try to convert any variable to number T (int64 | uint64 | float64).
//
// It doesn't convert slices, arrays and structs.
func ToNumber[T number](i interface{}) T {
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
		if float64(v)-math.Floor(float64(v)) > 0 {
			b := make([]byte, 4, 8)
			b = append(b, ToByteArray(v)...)
			return FromByteArray[T](b)
		}
		return T(v)
	case float64:
		if v-math.Floor(v) > 0 {
			b := ToByteArray(v)
			return FromByteArray[T](b)
		}
		return T(v)
	}

	return T(0)
}

// Convert any variable to bool.
//
// In most cases it returns true if i > 0.
func ToBool(i interface{}) bool {
	switch v := i.(type) {
	case bool:
		return v
	case string:
		return len(v) > 0
	case byte:
		return v > 0
	case int:
		return v > 0
	case uint:
		return v > 0
	case int64:
		return v > 0
	case uint64:
		return v > 0
	case float64:
		return v > 0
	}

	return false
}

// Helper function that converts byte slice to type T.
func FromByteArray[T any](b []byte) (s T) {
	buf := bytes.NewReader(b)
	binary.Read(buf, binary.LittleEndian, &s)
	return
}

// Helper function that converts type T to byte slice.
func ToByteArray[T any](s T) (b []byte) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, s)
	b = buf.Bytes()
	return
}
