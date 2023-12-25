package utils

import (
	"bytes"
	"encoding/binary"
	"math"
	"strings"
)

const (
	stringLiterals   = "_@QWERTYUIOPASDFGHJKLZXCVBNMqwertyuiopasdfghjklzxcvbnm"
	sciLiterals      = "0123456789.eE-+_"
	decLiterals      = "0123456789._"
	hexLiterals      = "0123456789ABCDEFabcdef_"
	binLiterals      = "01_"
	controlLiterals  = "()"
	operatorLiterals = ";#?:=-+*/%^!&|~<>,"
)

type wordType int

// Word types
const (
	// Not word.
	W_NONE wordType = iota
	// Number in scientific notation.
	W_NUM_SCI wordType = iota
	// Number in decimal representation.
	W_NUM_DEC wordType = iota
	// Number in hexadecimal representation.
	W_NUM_HEX wordType = iota
	// Number in binary representation.
	W_NUM_BIN wordType = iota
	// Some variable, constant or function.
	W_UNIT wordType = iota
	// Operator.
	W_OP wordType = iota
	// Flow control (brackets).
	W_CTL wordType = iota
	// Detected function call.
	W_FUNC wordType = iota
	// String.
	W_STR wordType = iota
)

type Word struct {
	Type    wordType
	Literal string
}

// Is two words equal.
func WordsEqual(a, b []Word) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v.Literal != b[i].Literal {
			return false
		}
	}
	return true
}

// This functions splits provided string into words and determines its basic types.
func ParsePrompt(str string) []Word {
	words := make([]Word, 0)

	wordType := W_NUM_DEC
	wordDone := false

	wordBegin := -1
	for i, c := range str {
		if wordBegin > -1 {
			switch wordType {
			case W_UNIT:
				if !(strings.Contains(stringLiterals, string(c)) || strings.Contains(decLiterals, string(c))) {
					wordDone = true
				}
			case W_NUM_DEC, W_NUM_HEX, W_NUM_BIN, W_NUM_SCI:
				if (c == 'x' || c == 'b') && i-wordBegin == 1 && wordType == W_NUM_DEC {
					if c == 'x' {
						wordType = W_NUM_HEX
					} else {
						wordType = W_NUM_BIN
					}
					wordBegin += 2
				} else if (c == 'e' || c == 'E') && wordType == W_NUM_DEC {
					wordType = W_NUM_SCI
				} else {
					switch wordType {
					case W_NUM_SCI:
						if !strings.Contains(sciLiterals, string(c)) {
							wordDone = true
						}
					case W_NUM_DEC:
						if !strings.Contains(decLiterals, string(c)) {
							wordDone = true
						}
					case W_NUM_HEX:
						if !strings.Contains(hexLiterals, string(c)) {
							wordDone = true
						}
					case W_NUM_BIN:
						if !strings.Contains(binLiterals, string(c)) {
							wordDone = true
						}
					}
				}
			case W_STR:
				if c == '"' {
					wordDone = true
					words = append(words, Word{wordType, str[wordBegin:i]})
					wordBegin = -1
					continue
				}
			case W_CTL:
				wordDone = true
			case W_OP:
				if !strings.Contains(operatorLiterals, string(c)) {
					wordDone = true
				}
			}

			if wordDone && wordType != W_NONE {
				words = append(words, Word{wordType, str[wordBegin:i]})
				wordBegin = -1
			}
		}

		if wordBegin < 0 {
			wordBegin = i
			wordDone = false

			if c == '"' {
				wordType = W_STR
				wordBegin++
			} else if strings.Contains(stringLiterals, string(c)) {
				wordType = W_UNIT
			} else if strings.Contains(decLiterals, string(c)) {
				wordType = W_NUM_DEC
			} else if strings.Contains(controlLiterals, string(c)) {
				wordType = W_CTL
			} else if strings.Contains(operatorLiterals, string(c)) {
				wordType = W_OP
			} else {
				wordBegin = -1
				wordType = W_NONE
			}
		}
	}

	if wordBegin > -1 && wordType != W_NONE {
		words = append(words, Word{wordType, str[wordBegin:]})
	}

	return words
}

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
