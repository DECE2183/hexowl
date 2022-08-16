package main

import (
	"bufio"
	"fmt"
	"math"
	"math/bits"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dece2183/hexowl/builtin"
	"github.com/dece2183/hexowl/user"
	"github.com/dece2183/hexowl/utils"
)

const (
	stringLiterals   = "@QWERTYUIOPASDFGHJKLZXCVBNMqwertyuiopasdfghjklzxcvbnm"
	decLiterals      = "0123456789."
	hexLiterals      = "0123456789ABCDEFabcdef"
	binLiterals      = "01"
	controlLiterals  = "()"
	operatorLiterals = "#?=-+*/%^!&|~<>,"
)

// Possible values:
// 	W_NONE
// 	W_NUM_DEC
// 	W_NUM_HEX
// 	W_NUM_BIN
// 	W_STR
// 	W_OP
// 	W_CTL
// 	W_FUNC
type WordType int

// Word types
const (
	W_NONE    WordType = iota
	W_NUM_DEC WordType = iota
	W_NUM_HEX WordType = iota
	W_NUM_BIN WordType = iota
	W_STR     WordType = iota
	W_OP      WordType = iota
	W_CTL     WordType = iota
	W_FUNC    WordType = iota
)

// Possible values:
//	OP_NONE
//	OP_ASSIGN
//	OP_DECREMENT
//	OP_INCREMENT
//	OP_ASSIGNMUL
//	OP_ASSIGNDIV
//	OP_LOGICNOT
//	OP_LOGICOR
//	OP_LOGICAND
//	OP_EQUALITY
//	OP_NOTEQUALITY
//	OP_MINUS
//	OP_PLUS
//	OP_MULTIPLY
//	OP_DIVIDE
//	OP_MODULO
//	OP_POWER
//	OP_LEFTSHIFT
//	OP_RIGHTSHIFT
//	OP_BITOR
//	OP_BITAND
//	OP_BITXOR
//	OP_BITCLEAR
//	OP_BITINVERSE
//	OP_POPCNT
//	OP_FUNCARGSEP
//	OP_USERFUNC
//	OP_BUILTINFUNC
type OperatorType int

// Operator types
const (
	OP_NONE OperatorType = iota

	OP_ASSIGN      OperatorType = iota
	OP_DECREMENT   OperatorType = iota
	OP_INCREMENT   OperatorType = iota
	OP_ASSIGNMUL   OperatorType = iota
	OP_ASSIGNDIV   OperatorType = iota
	OP_LOGICNOT    OperatorType = iota
	OP_LOGICOR     OperatorType = iota
	OP_LOGICAND    OperatorType = iota
	OP_EQUALITY    OperatorType = iota
	OP_NOTEQUALITY OperatorType = iota
	OP_MINUS       OperatorType = iota
	OP_PLUS        OperatorType = iota
	OP_MULTIPLY    OperatorType = iota
	OP_DIVIDE      OperatorType = iota
	OP_MODULO      OperatorType = iota
	OP_POWER       OperatorType = iota
	OP_LEFTSHIFT   OperatorType = iota
	OP_RIGHTSHIFT  OperatorType = iota
	OP_BITOR       OperatorType = iota
	OP_BITAND      OperatorType = iota
	OP_BITXOR      OperatorType = iota
	OP_BITCLEAR    OperatorType = iota
	OP_BITINVERSE  OperatorType = iota
	OP_POPCNT      OperatorType = iota

	OP_FUNCARGSEP  OperatorType = iota
	OP_USERFUNC    OperatorType = iota
	OP_BUILTINFUNC OperatorType = iota
)

type Word struct {
	Type    WordType
	Literal string
}

type Operator struct {
	Type     OperatorType
	OperandA *Operator
	OperandB *Operator
	Result   interface{}
}

var (
	operatorsPriorityList = [...]string{
		"=", "-=", "+=", "*=", "/=", ",", "==", "!=", "||", "&&", "!", "+", "-", "*", "**", "/", "%", "<<", ">>", "|", "&", "^", "&^", "&~", "~", "#",
	}
)

func main() {
	builtin.FuncsInit()

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
				fmt.Printf("\n\tTime:\t%d ms\r\n\n", calcTime.Milliseconds())
			}
		}
	}
}

func promt(reader *bufio.Reader) []Word {
	// return parse("help + 0")
	var input string

	fmt.Printf(">: ")
	input, _ = reader.ReadString('\n')

	return parse(input)
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

			if strings.Contains(stringLiterals, string(c)) {
				wordType = W_STR
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

func calculate(words []Word) error {
	operator, err := generateOperators(words)
	if err != nil {
		return err
	}

	val, err := calcOperator(operator)
	if err != nil {
		return err
	}

	if val != nil {
		switch v := val.(type) {
		case string:
			fmt.Printf("\n\t%s\r\n", v)
		case bool:
			fmt.Printf("\n\tResult:\t%v\r\n", v)
		default:
			fmt.Printf("\n\tResult:\t%v\r\n", utils.ToNumber[float64](val))
			fmt.Printf("\t\t0x%X\r\n", utils.ToNumber[uint64](val))
			fmt.Printf("\t\t0b%b\r\n", utils.ToNumber[uint64](val))
		}
	}

	return nil
}

func calcOperator(op *Operator) (interface{}, error) {
	var err error

	if op == nil {
		return 0, nil
	}

	if op.Type == OP_NONE {
		return op.Result, nil
	} else if op.OperandA == nil || op.OperandB == nil {
		return 0, nil
	}

	op.OperandA.Result, err = calcOperator(op.OperandA)
	if err != nil {
		return nil, err
	}
	op.OperandB.Result, err = calcOperator(op.OperandB)
	if err != nil {
		return nil, err
	}

	switch op.Type {
	case OP_ASSIGN:
		user.Variables[op.OperandA.Result.(string)] = op.OperandB.Result
		op.Result = op.OperandB.Result
	case OP_DECREMENT:
		op.Result = utils.ToNumber[float64](user.Variables[op.OperandA.Result.(string)]) - utils.ToNumber[float64](op.OperandB.Result)
		user.Variables[op.OperandA.Result.(string)] = op.Result
	case OP_INCREMENT:
		op.Result = utils.ToNumber[float64](user.Variables[op.OperandA.Result.(string)]) + utils.ToNumber[float64](op.OperandB.Result)
		user.Variables[op.OperandA.Result.(string)] = op.Result
	case OP_ASSIGNMUL:
		op.Result = utils.ToNumber[float64](user.Variables[op.OperandA.Result.(string)]) * utils.ToNumber[float64](op.OperandB.Result)
		user.Variables[op.OperandA.Result.(string)] = op.Result
	case OP_ASSIGNDIV:
		op.Result = utils.ToNumber[float64](user.Variables[op.OperandA.Result.(string)]) / utils.ToNumber[float64](op.OperandB.Result)
		user.Variables[op.OperandA.Result.(string)] = op.Result
	case OP_LOGICNOT:
		op.Result = !utils.ToBool(op.OperandB.Result)
	case OP_LOGICOR:
		op.Result = utils.ToBool(op.OperandA.Result) || utils.ToBool(op.OperandB.Result)
	case OP_LOGICAND:
		op.Result = utils.ToBool(op.OperandA.Result) && utils.ToBool(op.OperandB.Result)
	case OP_EQUALITY:
		op.Result = utils.ToNumber[uint64](op.OperandA.Result) == utils.ToNumber[uint64](op.OperandB.Result)
	case OP_NOTEQUALITY:
		op.Result = utils.ToNumber[uint64](op.OperandA.Result) != utils.ToNumber[uint64](op.OperandB.Result)
	case OP_MINUS:
		op.Result = utils.ToNumber[float64](op.OperandA.Result) - utils.ToNumber[float64](op.OperandB.Result)
	case OP_PLUS:
		op.Result = utils.ToNumber[float64](op.OperandA.Result) + utils.ToNumber[float64](op.OperandB.Result)
	case OP_MULTIPLY:
		op.Result = utils.ToNumber[float64](op.OperandA.Result) * utils.ToNumber[float64](op.OperandB.Result)
	case OP_DIVIDE:
		op.Result = utils.ToNumber[float64](op.OperandA.Result) / utils.ToNumber[float64](op.OperandB.Result)
	case OP_MODULO:
		op.Result = utils.ToNumber[int64](op.OperandA.Result) % utils.ToNumber[int64](op.OperandB.Result)
	case OP_POWER:
		op.Result = math.Pow(utils.ToNumber[float64](op.OperandA.Result), utils.ToNumber[float64](op.OperandB.Result))
	case OP_LEFTSHIFT:
		op.Result = utils.ToNumber[uint64](op.OperandA.Result) << utils.ToNumber[uint64](op.OperandB.Result)
	case OP_RIGHTSHIFT:
		op.Result = utils.ToNumber[uint64](op.OperandA.Result) >> utils.ToNumber[uint64](op.OperandB.Result)
	case OP_BITOR:
		op.Result = utils.ToNumber[uint64](op.OperandA.Result) | utils.ToNumber[uint64](op.OperandB.Result)
	case OP_BITAND:
		op.Result = utils.ToNumber[uint64](op.OperandA.Result) & utils.ToNumber[uint64](op.OperandB.Result)
	case OP_BITXOR:
		op.Result = utils.ToNumber[uint64](op.OperandA.Result) ^ utils.ToNumber[uint64](op.OperandB.Result)
	case OP_BITCLEAR:
		op.Result = utils.ToNumber[uint64](op.OperandA.Result) &^ utils.ToNumber[uint64](op.OperandB.Result)
	case OP_POPCNT:
		op.Result = uint64(bits.OnesCount64(utils.ToNumber[uint64](op.OperandB.Result)))
	case OP_BITINVERSE:
		op.Result = 0xFFFFFFFFFFFFFFFF ^ utils.ToNumber[uint64](op.OperandB.Result)

	case OP_FUNCARGSEP:
		switch op.OperandA.Result.(type) {
		case []interface{}:
			break
		default:
			op.OperandA.Result = []interface{}{op.OperandA.Result}
		}
		switch op.OperandB.Result.(type) {
		case []interface{}:
			break
		default:
			op.OperandB.Result = []interface{}{op.OperandB.Result}
		}
		op.Result = append(op.OperandA.Result.([]interface{}), op.OperandB.Result.([]interface{})...)
	case OP_BUILTINFUNC:
		switch op.OperandB.Result.(type) {
		case []interface{}:
			op.Result, err = builtin.Functions[op.OperandA.Result.(string)].Exec(op.OperandB.Result.([]interface{})...)
		default:
			op.Result, err = builtin.Functions[op.OperandA.Result.(string)].Exec(op.OperandB.Result)
		}
	case OP_USERFUNC:
		op.Result = 0
	}

	return op.Result, err
}

func getVariable(literal string) (val interface{}, found bool) {
	val, found = user.Variables[literal]
	if !found {
		val, found = builtin.Constants[literal]
	}
	return
}

func generateOperators(words []Word) (*Operator, error) {
	var err error
	newOp := &Operator{}

	if len(words) == 0 {
		newOp.Result = uint64(0)
		return newOp, nil
	} else if len(words) == 1 {
		w := words[0]
		switch w.Type {
		case W_STR:
			// Try to find variable
			val, found := getVariable(w.Literal)
			if !found {
				return nil, fmt.Errorf("there is no variable named '%s'", w.Literal)
			}
			newOp.Result = val

		case W_FUNC:
			// Try to find function
			_, found := user.Functions[w.Literal]
			if found {
				newOp.Type = OP_USERFUNC
				newOp.Result = w.Literal
				break
			}
			_, found = builtin.Functions[w.Literal]
			if found {
				newOp.Type = OP_BUILTINFUNC
				newOp.Result = w.Literal
				break
			}
			return nil, fmt.Errorf("there is no function named '%s'", w.Literal)

		case W_NUM_DEC, W_NUM_HEX, W_NUM_BIN:
			switch w.Type {
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

	minPriority := len(operatorsPriorityList)
	minPriorityIndex := 0
	var minPriorityWord *Word

	bracketsCount := 0

	if words[0].Type == W_CTL && words[len(words)-1].Type == W_CTL {
		if words[0].Literal != "(" {
			return nil, fmt.Errorf("missing opening bracket")
		}
		if words[len(words)-1].Literal != ")" {
			return nil, fmt.Errorf("missing closing bracket")
		}

		for i := 1; i < len(words)-1; i++ {
			if words[i].Type != W_CTL {
				continue
			}
			if words[i].Literal == "(" {
				bracketsCount++
			} else {
				bracketsCount--
				if bracketsCount < 0 {
					break
				}
			}
		}

		if bracketsCount >= 0 {
			words = words[1 : len(words)-1]
		}
	}

	bracketsCount = 0

	for i := 0; i < len(words); i++ {
		w := words[i]

		if w.Type == W_CTL {
			if w.Literal == "(" {
				bracketsCount++
			} else {
				bracketsCount--
			}
			continue
		} else if w.Type == W_STR {
			// Function call detect
			if i+1 < len(words)-1 && words[i+1].Type == W_CTL && words[i+1].Literal == "(" {
				words[i].Type = W_FUNC
			}
			continue
		}

		if w.Type != W_OP || bracketsCount > 0 {
			continue
		}

		prio := -1

		if getOperatorType(w.Literal) == OP_MINUS && (i == 0 || words[i-1].Type == W_OP) {
			// If it is a single minus operator give it the max prioriy
			prio = len(operatorsPriorityList)
		} else {
			for pr, lit := range operatorsPriorityList {
				if lit == w.Literal {
					prio = pr
					break
				}
			}
		}

		if prio < 0 {
			return nil, fmt.Errorf("unknown operator '%s'", w.Literal)
		}

		if prio <= minPriority {
			minPriority = prio
			minPriorityIndex = i
			minPriorityWord = &words[i]
		}
	}

	if minPriorityWord == nil {
		if words[0].Type == W_FUNC {
			if len(words) < 3 {
				return nil, fmt.Errorf("missing function '%s' arguments", words[0].Literal)
			}
			newOp.OperandA, err = generateOperators(words[:1])
			if err != nil {
				return nil, err
			}
			newOp.Type = newOp.OperandA.Type
			newOp.OperandA.Type = OP_NONE
			if len(words) > 3 {
				newOp.OperandB, err = generateOperators(words[2 : len(words)-1])
				if err != nil {
					return nil, err
				}
			} else {
				newOp.OperandB = &Operator{}
			}
		} else {
			return nil, fmt.Errorf("operators not found")
		}
	} else {
		newOp.Type = getOperatorType(minPriorityWord.Literal)
		if newOp.Type < 0 {
			return nil, fmt.Errorf("unknown operator '%s'", minPriorityWord.Literal)
		}

		if newOp.Type == OP_BITINVERSE || newOp.Type == OP_POPCNT || newOp.Type == OP_LOGICNOT {
			// One side operators
			newOp.OperandA = &Operator{}
		} else if newOp.Type >= OP_ASSIGN && newOp.Type <= OP_ASSIGNDIV {
			// Assign operators
			if minPriorityIndex < 1 {
				return nil, fmt.Errorf("missing a variable or function declaration on left side of operator '%s'", minPriorityWord.Literal)
			}
			lit := words[minPriorityIndex-1].Literal
			if newOp.Type > OP_ASSIGN {
				if user.Variables[lit] == nil {
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

func getOperatorType(op string) OperatorType {
	switch op {
	case "=":
		return OP_ASSIGN
	case "-=":
		return OP_DECREMENT
	case "+=":
		return OP_INCREMENT
	case "*=":
		return OP_ASSIGNMUL
	case "/=":
		return OP_ASSIGNDIV
	case ",":
		return OP_FUNCARGSEP
	case "!":
		return OP_LOGICNOT
	case "||":
		return OP_LOGICOR
	case "&&":
		return OP_LOGICAND
	case "==":
		return OP_EQUALITY
	case "!=":
		return OP_NOTEQUALITY
	case "-":
		return OP_MINUS
	case "+":
		return OP_PLUS
	case "*":
		return OP_MULTIPLY
	case "**":
		return OP_POWER
	case "/":
		return OP_DIVIDE
	case "%":
		return OP_MODULO
	case "<<":
		return OP_LEFTSHIFT
	case ">>":
		return OP_RIGHTSHIFT
	case "|":
		return OP_BITOR
	case "&":
		return OP_BITAND
	case "^":
		return OP_BITXOR
	case "&^", "&~":
		return OP_BITCLEAR
	case "~":
		return OP_BITINVERSE
	case "#":
		return OP_POPCNT
	default:
		return -1
	}
}
