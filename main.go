package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type number interface {
	int64 | float64
}

const (
	wordDividers      = " -=+-/*%!&|$#?.,:;\"'~`()[]{}<>\n"
	stringLiterals    = "@QWERTYUIOPASDFGHJKLZXCVBNMqwertyuiopasdfghjklzxcvbnm"
	operatorLiterals  = "=-+*/%^!&|<>"
	controlLiterals   = "()"
	decLiterals       = "0123456789"
	hexLiterals       = "0123456789ABCDEFabcdef"
	binLiterals       = "01"
	operatorsPriority = "& | << - + % / * -= += ="
)

// Word types
const (
	W_NONE = iota
	W_NUM_DEC
	W_NUM_HEX
	W_NUM_BIN
	W_STR
	W_OP
	W_CTL
)

// Operator types
const (
	OP_NONE = iota
	OP_ASSIGN
)

type Word struct {
	WordType int
	Literal  string
}

type Operator struct {
	OpType   int
	OperandA *Operator
	OperandB *Operator
	Result   interface{}
}

type CalcFunc func(args ...interface{}) interface{}

var (
	userVars    = map[string]interface{}{}
	builtinVars = map[string]interface{}{
		"pi": float64(3.14159265358979323846),
		"PI": float64(3.14159265358979323846),
	}

	userFuncs    = map[string]CalcFunc{}
	builtinFuncs = map[string]CalcFunc{}
)

func main() {
	var words []Word
	stdreader := bufio.NewReader(os.Stdin)

	for {
		words = promt(stdreader)
		if len(words) > 0 {
			err := calculate(words)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func promt(reader *bufio.Reader) []Word {
	var input string

	fmt.Printf(">: ")
	input, _ = reader.ReadString('\n')

	return parse(input)
	// return parse("x = a + 23 * 10(2-3)c")
}

func parse(str string) []Word {
	words := make([]Word, 0)

	wordType := W_NUM_DEC
	wordDone := false

	wordBegin := -1
	for i, c := range str {
		if wordBegin > -1 {
			switch wordType {
			case W_STR:
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
			case W_OP:
				if !strings.Contains(operatorLiterals, string(c)) {
					wordDone = true
				}
			case W_CTL:
				if !strings.Contains(controlLiterals, string(c)) {
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

			if strings.Contains(stringLiterals, string(c)) {
				wordType = W_STR
			} else if strings.Contains(decLiterals, string(c)) {
				wordType = W_NUM_DEC
			} else if strings.Contains(operatorLiterals, string(c)) {
				wordType = W_OP
			} else if strings.Contains(controlLiterals, string(c)) {
				wordType = W_CTL
			} else {
				wordBegin = -1
				wordType = W_NONE
			}
		}
	}

	if wordBegin > -1 && wordType != W_NONE {
		words = append(words, Word{wordType, str[wordBegin:]})
	}

	fmt.Println(words)
	return words
}

func generateOperators(words []Word) (*Operator, error) {
	var err error
	newOp := &Operator{}

	if len(words) == 1 {
		w := words[0]
		switch w.WordType {

		}
	}

	innerBegin := -1
	maxPriority := 0

	var maxPriorityOperator *Word

	for i, w := range words {
		if innerBegin < 0 {
			switch w.WordType {
			case W_CTL:
				if w.WordType == W_CTL && w.Literal == "(" {
					innerBegin = i + 1
				}
			case W_OP:
				prio := strings.Index(operatorsPriority, w.Literal)
				if prio > maxPriority {
					maxPriority = prio
					maxPriorityOperator = &w
				}
			}
		} else {

		}
	}

	for i, w := range words {
		if innerBegin < 0 {
			switch w.WordType {
			case W_CTL:
				if w.WordType == W_CTL && w.Literal == "(" {
					innerBegin = i
				}
			case W_OP:

			case W_STR:
				if newOp.OpType == OP_NONE {

				} else {
					val, found := userVars[w.Literal]
					if !found {
						val, found = builtinVars[w.Literal]
					}
					if !found {
						val, found = userFuncs[w.Literal]
					}
					if !found {
						val, found = builtinFuncs[w.Literal]
					}
					if !found {
						return nil, fmt.Errorf("there is no variable or function named '%s'", w.Literal)
					}
					newOp.OperandB = &Operator{
						Result: val,
					}
				}
			case W_NUM_DEC, W_NUM_HEX, W_NUM_BIN:
				newNumber := &Operator{}

				switch w.WordType {
				case W_NUM_DEC:
					newNumber.Result, err = strconv.ParseFloat(w.Literal, 64)
					if err != nil {
						return nil, fmt.Errorf("unable to parse literal '%s' as number", w.Literal)
					}
				case W_NUM_HEX:
					newNumber.Result, err = strconv.ParseUint(w.Literal, 16, 64)
					if err != nil {
						return nil, fmt.Errorf("unable to parse literal '%s' as hex number", w.Literal)
					}
				case W_NUM_BIN:
					newNumber.Result, err = strconv.ParseUint(w.Literal, 2, 64)
					if err != nil {
						return nil, fmt.Errorf("unable to parse literal '%s' as bin number", w.Literal)
					}
				}

				if newOp.OpType == OP_NONE {
					newOp.OperandA = newNumber
				} else {
					newOp.OperandB = newNumber
				}
			}
		} else if w.WordType == W_CTL && w.Literal == ")" {
			if newOp.OpType == OP_NONE {
				newOp.OperandA, err = generateOperators(words[innerBegin+1 : i])
			} else {
				newOp.OperandB, err = generateOperators(words[innerBegin+1 : i])
			}
			if err != nil {
				return nil, err
			}
			innerBegin = -1
		}
	}

	return newOp, nil
}

func calculate(words []Word) error {
	operator, err := generateOperators(words)
	if err != nil {
		return err
	}

	return nil
}
