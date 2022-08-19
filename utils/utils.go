package utils

import (
	"math"
	"strings"
)

const (
	stringLiterals   = "_@QWERTYUIOPASDFGHJKLZXCVBNMqwertyuiopasdfghjklzxcvbnm"
	decLiterals      = "0123456789."
	hexLiterals      = "0123456789ABCDEFabcdef"
	binLiterals      = "01"
	controlLiterals  = "()"
	operatorLiterals = ";#?:=-+*/%^!&|~<>,"
)

type wordType int

// Word types
const (
	W_NONE    wordType = iota
	W_NUM_DEC wordType = iota
	W_NUM_HEX wordType = iota
	W_NUM_BIN wordType = iota
	W_UNIT    wordType = iota
	W_OP      wordType = iota
	W_CTL     wordType = iota
	W_FUNC    wordType = iota
	W_STR     wordType = iota
)

type Word struct {
	Type    wordType
	Literal string
}

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
			case W_NUM_DEC, W_NUM_HEX, W_NUM_BIN:
				if (c == 'x' || c == 'b') && i-wordBegin == 1 {
					if c == 'x' {
						wordType = W_NUM_HEX
					} else {
						wordType = W_NUM_BIN
					}
					wordBegin += 2
				} else {
					switch wordType {
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
