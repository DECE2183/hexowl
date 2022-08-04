package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

type number interface {
	int64 | uint64 | float64
}

const (
	wordDividers      = " -=+-/*%!&|$#?.,:;\"'~`()[]{}<>\n"
	stringLiterals    = "@QWERTYUIOPASDFGHJKLZXCVBNMqwertyuiopasdfghjklzxcvbnm"
	decLiterals       = "0123456789"
	hexLiterals       = "0123456789ABCDEFabcdef"
	binLiterals       = "01"
	operatorLiterals  = "?=-+*/%^!&|~<>()"
	operatorsPriority = "= -= += *= /= ?= || && - + * / % ^ << >> | & ~ ! ("
)

// Word types
const (
	W_NONE = iota
	W_NUM_DEC
	W_NUM_HEX
	W_NUM_BIN
	W_STR
	W_OP
)

// Operator types
const (
	OP_NONE = iota
	OP_ASSIGN
	OP_DECREMENT
	OP_INCREMENT
	OP_ASSIGNMUL
	OP_ASSIGNDIV
	OP_LOGICNOT
	OP_LOGICOR
	OP_LOGICAND
	OP_EQUALITY
	OP_MINUS
	OP_PLUS
	OP_MULTIPLY
	OP_DIVIDE
	OP_MODULO
	OP_POWER
	OP_LEFTSHIFT
	OP_RIGHTSHIFT
	OP_BITOR
	OP_BITAND
	OP_BITINVERSE
	OP_BRACKET
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
	userVars         = map[string]interface{}{}
	userFuncs        = map[string]CalcFunc{}
	builtinFuncs     = map[string]CalcFunc{}
	builtinConstants = map[string]interface{}{
		"pi": float64(3.14159265358979323846),
		"PI": float64(3.14159265358979323846),
	}
)

func main() {
	var words []Word
	stdreader := bufio.NewReader(os.Stdin)

	for {
		words = promt(stdreader)
		if len(words) > 0 {
			calcBeginTime := time.Now()
			err := calculate(words)
			calcTime := time.Since(calcBeginTime)

			if err != nil {
				fmt.Printf("\n\tError occurred: %s\n\n", err)
			} else {
				fmt.Printf("\n\tTime:   %d ms\r\n\n", calcTime.Milliseconds())
			}
		}
	}
}

func promt(reader *bufio.Reader) []Word {
	var input string

	fmt.Printf(">: ")
	input, _ = reader.ReadString('\n')

	return parse(input)
	// return parse("~0b1100")
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

			if strings.Contains(stringLiterals, string(c)) {
				wordType = W_STR
			} else if strings.Contains(decLiterals, string(c)) {
				wordType = W_NUM_DEC
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

func calculate(words []Word) error {
	operator, err := generateOperators(words)
	if err != nil {
		return err
	}

	val := calcOperator(operator)

	fmt.Printf("\n\tResult: %v\r\n", toNumber[float64](val))
	fmt.Printf("\t        0x%X\r\n", toNumber[uint64](val))
	fmt.Printf("\t        0b%b\r\n", toNumber[uint64](val))

	return nil
}

func calcOperator(op *Operator) interface{} {
	if op.OpType == OP_NONE {
		return op.Result
	} else if op.OperandA == nil || op.OperandB == nil {
		return 0
	}

	op.OperandA.Result = calcOperator(op.OperandA)
	op.OperandB.Result = calcOperator(op.OperandB)

	switch op.OpType {
	case OP_ASSIGN:
		userVars[op.OperandA.Result.(string)] = op.OperandB.Result
		op.Result = op.OperandB.Result
	case OP_DECREMENT:
		op.Result = toNumber[float64](userVars[op.OperandA.Result.(string)]) - toNumber[float64](op.OperandB.Result)
		userVars[op.OperandA.Result.(string)] = op.Result
	case OP_INCREMENT:
		op.Result = toNumber[float64](userVars[op.OperandA.Result.(string)]) + toNumber[float64](op.OperandB.Result)
		userVars[op.OperandA.Result.(string)] = op.Result
	case OP_ASSIGNMUL:
		op.Result = toNumber[float64](userVars[op.OperandA.Result.(string)]) * toNumber[float64](op.OperandB.Result)
		userVars[op.OperandA.Result.(string)] = op.Result
	case OP_ASSIGNDIV:
		op.Result = toNumber[float64](userVars[op.OperandA.Result.(string)]) / toNumber[float64](op.OperandB.Result)
		userVars[op.OperandA.Result.(string)] = op.Result
	case OP_LOGICNOT:
		op.Result = !toBool(op.OperandB.Result)
	case OP_LOGICOR:
		op.Result = toBool(op.OperandA.Result) || toBool(op.OperandB.Result)
	case OP_LOGICAND:
		op.Result = toBool(op.OperandA.Result) && toBool(op.OperandB.Result)
	case OP_EQUALITY:
		op.Result = toNumber[uint64](op.OperandA.Result) == toNumber[uint64](op.OperandB.Result)
	case OP_MINUS:
		op.Result = toNumber[float64](op.OperandA.Result) - toNumber[float64](op.OperandB.Result)
	case OP_PLUS:
		op.Result = toNumber[float64](op.OperandA.Result) + toNumber[float64](op.OperandB.Result)
	case OP_MULTIPLY:
		op.Result = toNumber[float64](op.OperandA.Result) * toNumber[float64](op.OperandB.Result)
	case OP_DIVIDE:
		op.Result = toNumber[float64](op.OperandA.Result) / toNumber[float64](op.OperandB.Result)
	case OP_MODULO:
		op.Result = toNumber[int64](op.OperandA.Result) % toNumber[int64](op.OperandB.Result)
	case OP_POWER:
		op.Result = math.Pow(toNumber[float64](op.OperandA.Result), toNumber[float64](op.OperandB.Result))
	case OP_LEFTSHIFT:
		op.Result = toNumber[uint64](op.OperandA.Result) << toNumber[uint64](op.OperandB.Result)
	case OP_RIGHTSHIFT:
		op.Result = toNumber[uint64](op.OperandA.Result) >> toNumber[uint64](op.OperandB.Result)
	case OP_BITOR:
		op.Result = toNumber[uint64](op.OperandA.Result) | toNumber[uint64](op.OperandB.Result)
	case OP_BITAND:
		op.Result = toNumber[uint64](op.OperandA.Result) & toNumber[uint64](op.OperandB.Result)
	case OP_BITINVERSE:
		op.Result = 0xFFFFFFFFFFFFFFFF ^ toNumber[uint64](op.OperandB.Result)
	case OP_BRACKET:
		op.Result = 0
	}

	return op.Result
}

func tryGetVar(literal string) (val interface{}, found bool) {
	val, found = userVars[literal]
	if !found {
		val, found = builtinConstants[literal]
	}
	return
}

func generateOperators(words []Word) (*Operator, error) {
	var err error
	newOp := &Operator{}

	if len(words) == 1 {
		w := words[0]
		switch w.WordType {
		case W_STR:
			// Try to find variable of function
			val, found := tryGetVar(w.Literal)
			if !found {
				val, found = userFuncs[w.Literal]
				if !found {
					val, found = builtinFuncs[w.Literal]
					if !found {
						return nil, fmt.Errorf("there is no variable or function named '%s'", w.Literal)
					}
				}
			}
			newOp.Result = val

		case W_NUM_DEC, W_NUM_HEX, W_NUM_BIN:
			switch w.WordType {
			case W_NUM_DEC:
				newOp.Result, err = strconv.ParseFloat(w.Literal, 64)
				if err != nil {
					return nil, fmt.Errorf("unable to parse literal '%s' as number", w.Literal)
				}
			case W_NUM_HEX:
				newOp.Result, err = strconv.ParseUint(w.Literal, 16, 64)
				if err != nil {
					return nil, fmt.Errorf("unable to parse literal '%s' as hex number", w.Literal)
				}
			case W_NUM_BIN:
				newOp.Result, err = strconv.ParseUint(w.Literal, 2, 64)
				if err != nil {
					return nil, fmt.Errorf("unable to parse literal '%s' as bin number", w.Literal)
				}
			}
		}

		return newOp, nil
	}

	minPriority := len(operatorsPriority)
	minPriorityIndex := 0
	var minPriorityWord *Word

	for i, w := range words {
		if w.WordType != W_OP {
			continue
		}

		prio := strings.Index(operatorsPriority, w.Literal)
		if prio < 0 {
			return nil, fmt.Errorf("unknown operator '%s'", w.Literal)
		}

		if prio < minPriority {
			minPriority = prio
			minPriorityIndex = i
			minPriorityWord = &words[i]
		}
	}

	if minPriorityWord.Literal == "(" {

	} else {
		switch minPriorityWord.Literal {
		case "=":
			newOp.OpType = OP_ASSIGN
		case "-=":
			newOp.OpType = OP_DECREMENT
		case "+=":
			newOp.OpType = OP_INCREMENT
		case "*=":
			newOp.OpType = OP_ASSIGNMUL
		case "/=":
			newOp.OpType = OP_ASSIGNDIV
		case "!":
			newOp.OpType = OP_LOGICNOT
		case "||":
			newOp.OpType = OP_LOGICOR
		case "&&":
			newOp.OpType = OP_LOGICAND
		case "?=":
			newOp.OpType = OP_EQUALITY
		case "-":
			newOp.OpType = OP_MINUS
		case "+":
			newOp.OpType = OP_PLUS
		case "*":
			newOp.OpType = OP_MULTIPLY
		case "/":
			newOp.OpType = OP_DIVIDE
		case "%":
			newOp.OpType = OP_MODULO
		case "^":
			newOp.OpType = OP_POWER
		case "<<":
			newOp.OpType = OP_LEFTSHIFT
		case ">>":
			newOp.OpType = OP_RIGHTSHIFT
		case "|":
			newOp.OpType = OP_BITOR
		case "&":
			newOp.OpType = OP_BITAND
		case "~":
			newOp.OpType = OP_BITINVERSE
		case "(":
			newOp.OpType = OP_BRACKET
		default:
			return nil, fmt.Errorf("unknown operator '%s'", minPriorityWord.Literal)
		}

		if newOp.OpType == OP_BITINVERSE || newOp.OpType == OP_LOGICNOT {
			// One side operators
			newOp.OperandA = &Operator{}
		} else if newOp.OpType >= OP_ASSIGN && newOp.OpType <= OP_ASSIGNDIV {
			// Assign operators
			if minPriorityIndex < 1 {
				return nil, fmt.Errorf("missing a variable or function declaration on left side of operator '%s'", minPriorityWord.Literal)
			}
			lit := words[minPriorityIndex-1].Literal
			if newOp.OpType > OP_ASSIGN {
				if userVars[lit] == nil {
					return nil, fmt.Errorf("there is no variable named '%s'", lit)
				}
			}
			newOp.OperandA = &Operator{
				Result: lit,
			}
		} else {
			newOp.OperandA, err = generateOperators(words[:minPriorityIndex])
			if err != nil {
				return nil, err
			}
		}
		newOp.OperandB, err = generateOperators(words[minPriorityIndex+1:])
		if err != nil {
			return nil, err
		}
	}

	return newOp, nil
}

func toNumber[T number](i interface{}) T {
	switch v := i.(type) {
	case bool:
		if v == true {
			return T(1)
		} else {
			return T(0)
		}
	case int64:
		return T(v)
	case uint64:
		return T(v)
	case float64:
		return T(v)
	}

	return T(0)
}

func toBool(i interface{}) bool {
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
