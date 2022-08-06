package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math"
	"math/bits"
	"os"
	"strconv"
	"strings"
	"time"
)

type number interface {
	int64 | uint64 | float64
}

const (
	stringLiterals   = "@QWERTYUIOPASDFGHJKLZXCVBNMqwertyuiopasdfghjklzxcvbnm"
	decLiterals      = "0123456789."
	hexLiterals      = "0123456789ABCDEFabcdef"
	binLiterals      = "01"
	controlLiterals  = "()"
	operatorLiterals = "#?=-+*/%^!&|~<>,"
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
	W_FUNC
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
	OP_NOTEQUALITY
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
	OP_BITXOR
	OP_BITCLEAR
	OP_BITINVERSE
	OP_POPCNT

	OP_FUNCARGSEP
	OP_USERFUNC
	OP_BUILTINFUNC
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

type userFunc map[string]string
type builtinFunc func(args ...interface{}) (interface{}, error)

type saveStruct struct {
	UserVars  map[string]interface{}
	UserFuncs map[string]userFunc
}

var (
	userVars     = map[string]interface{}{}
	userFuncs    = map[string]userFunc{}
	builtinFuncs = map[string]builtinFunc{
		"sin": func(args ...interface{}) (interface{}, error) {
			if len(args) < 1 {
				return nil, fmt.Errorf("not enough arguments")
			}
			return math.Sin(toNumber[float64](args[0])), nil
		},
		"cos": func(args ...interface{}) (interface{}, error) {
			if len(args) < 1 {
				return nil, fmt.Errorf("not enough arguments")
			}
			return math.Cos(toNumber[float64](args[0])), nil
		},
		"pow": func(args ...interface{}) (interface{}, error) {
			if len(args) < 2 {
				return nil, fmt.Errorf("not enough arguments")
			}
			return math.Pow(toNumber[float64](args[0]), toNumber[float64](args[1])), nil
		},
		"exit": func(args ...interface{}) (interface{}, error) {
			exitCode := toNumber[int64](args[0])
			os.Exit(int(exitCode))
			return exitCode, nil
		},
		"vars": func(args ...interface{}) (interface{}, error) {
			varsCount := uint64(len(userVars))
			if varsCount > 0 {
				fmt.Printf("\n\tUser variables:\n")
				for key, value := range userVars {
					fmt.Printf("\t\t[%s] = %v\n", key, value)
				}
			} else {
				fmt.Printf("\n\tThere is no user defined variables.\n")
			}
			if len(builtinConstants) > 0 {
				fmt.Printf("\n\tBuiltin constants:\n")
				for key, value := range builtinConstants {
					fmt.Printf("\t\t[%s] = %v\n", key, value)
				}
			} else {
				fmt.Printf("\n\tThere is no builtin constants.\n")
			}
			return varsCount, nil
		},
		"save": func(args ...interface{}) (interface{}, error) {
			envID := toNumber[uint64](args[0])
			userDir, err := os.UserHomeDir()
			if err != nil {
				return nil, fmt.Errorf("unable to get user home directory")
			}
			saveDir := fmt.Sprintf("%s/.pcalc/environment", userDir)
			err = os.MkdirAll(saveDir, 0666)
			if err != nil {
				return nil, fmt.Errorf("unable to create save directory")
			}
			savePath := fmt.Sprintf("%s/0x%016X.json", saveDir, envID)
			saveData := saveStruct{
				UserVars:  userVars,
				UserFuncs: userFuncs,
			}
			saveJson, err := json.Marshal(saveData)
			if err != nil {
				return nil, fmt.Errorf("unable to create data")
			}
			err = os.WriteFile(savePath, saveJson, 0666)
			if err != nil {
				return nil, fmt.Errorf("unable to write data to file")
			}
			fmt.Printf("\n\tSaving environment as 0x%016X\n", envID)
			return 0, nil
		},
		"load": func(args ...interface{}) (interface{}, error) {
			envID := toNumber[uint64](args[0])
			userDir, err := os.UserHomeDir()
			if err != nil {
				return nil, fmt.Errorf("unable to get user home directory")
			}
			loadPath := fmt.Sprintf("%s/.pcalc/environment/0x%016X.json", userDir, envID)
			loadData := saveStruct{}
			loadBuffer, err := os.ReadFile(loadPath)
			if err != nil {
				return nil, fmt.Errorf("environment doesn't exists")
			}
			err = json.Unmarshal(loadBuffer, &loadData)
			if err != nil {
				return nil, fmt.Errorf("unable to parse environment data")
			}
			userVars = loadData.UserVars
			userFuncs = loadData.UserFuncs
			fmt.Printf("\n\tEnvironment 0x%016X loaded\n", envID)
			return 0, nil
		},
	}
	builtinConstants = map[string]interface{}{
		"pi":    math.Pi,
		"true":  true,
		"false": false,
	}
)

var (
	operatorsPriorityList = [...]string{
		"=", "-=", "+=", "*=", "/=", ",", "==", "!=", "||", "&&", "-", "+", "*", "**", "/", "%", "<<", ">>", "|", "&", "^", "&^", "~", "#", "!",
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
				fmt.Printf("\n\tTime:\t%d ms\r\n\n", calcTime.Milliseconds())
			}
		}
	}
}

func promt(reader *bufio.Reader) []Word {
	// return parse("5 + sin(3)")
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

	fmt.Printf("\n\tResult:\t%v\r\n", toNumber[float64](val))
	fmt.Printf("\t\t0x%X\r\n", toNumber[uint64](val))
	fmt.Printf("\t\t0b%b\r\n", toNumber[uint64](val))

	return nil
}

func calcOperator(op *Operator) (interface{}, error) {
	var err error

	if op.OpType == OP_NONE {
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
	case OP_NOTEQUALITY:
		op.Result = toNumber[uint64](op.OperandA.Result) != toNumber[uint64](op.OperandB.Result)
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
	case OP_BITXOR:
		op.Result = toNumber[uint64](op.OperandA.Result) ^ toNumber[uint64](op.OperandB.Result)
	case OP_BITCLEAR:
		op.Result = toNumber[uint64](op.OperandA.Result) &^ toNumber[uint64](op.OperandB.Result)
	case OP_POPCNT:
		op.Result = uint64(bits.OnesCount64(toNumber[uint64](op.OperandB.Result)))
	case OP_BITINVERSE:
		op.Result = 0xFFFFFFFFFFFFFFFF ^ toNumber[uint64](op.OperandB.Result)

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
			op.Result, err = builtinFuncs[op.OperandA.Result.(string)](op.OperandB.Result.([]interface{})...)
		default:
			op.Result, err = builtinFuncs[op.OperandA.Result.(string)](op.OperandB.Result)
		}
	case OP_USERFUNC:
		op.Result = 0
	}

	return op.Result, err
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
			// Try to find variable
			val, found := tryGetVar(w.Literal)
			if !found {
				return nil, fmt.Errorf("there is no variable named '%s'", w.Literal)
			}
			newOp.Result = val

		case W_FUNC:
			// Try to find function
			_, found := userFuncs[w.Literal]
			if found {
				newOp.OpType = OP_USERFUNC
				newOp.Result = w.Literal
				break
			}
			_, found = builtinFuncs[w.Literal]
			if found {
				newOp.OpType = OP_BUILTINFUNC
				newOp.Result = w.Literal
				break
			}
			return nil, fmt.Errorf("there is no function named '%s'", w.Literal)

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

	minPriority := len(operatorsPriorityList)
	minPriorityIndex := 0
	var minPriorityWord *Word

	bracketsCount := 0

	if words[0].WordType == W_CTL && words[len(words)-1].WordType == W_CTL {
		if words[0].Literal != "(" {
			return nil, fmt.Errorf("missing opening bracket")
		}
		if words[len(words)-1].Literal != ")" {
			return nil, fmt.Errorf("missing closing bracket")
		}

		for i := 1; i < len(words)-1; i++ {
			if words[i].WordType != W_CTL {
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

		if w.WordType == W_CTL {
			if w.Literal == "(" {
				bracketsCount++
			} else {
				bracketsCount--
			}
			continue
		} else if w.WordType == W_STR {
			// Function call detect
			if i+1 < len(words)-1 && words[i+1].WordType == W_CTL && words[i+1].Literal == "(" {
				words[i].WordType = W_FUNC
			}
			continue
		}

		if w.WordType != W_OP || bracketsCount > 0 {
			continue
		}

		prio := -1
		for pr, lit := range operatorsPriorityList {
			if lit == w.Literal {
				prio = pr
				break
			}
		}

		if prio < 0 {
			return nil, fmt.Errorf("unknown operator '%s'", w.Literal)
		}

		if prio < minPriority {
			minPriority = prio
			minPriorityIndex = i
			minPriorityWord = &words[i]
		}
	}

	if minPriorityWord == nil {
		if words[0].WordType == W_FUNC {
			if len(words) < 3 {
				return nil, fmt.Errorf("missing function '%s' arguments", words[0].Literal)
			}
			newOp.OperandA, err = generateOperators(words[:1])
			if err != nil {
				return nil, err
			}
			newOp.OpType = newOp.OperandA.OpType
			newOp.OperandA.OpType = OP_NONE
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
		newOp.OpType = getOperatorType(minPriorityWord.Literal)
		if newOp.OpType < 0 {
			return nil, fmt.Errorf("unknown operator '%s'", minPriorityWord.Literal)
		}

		if newOp.OpType == OP_BITINVERSE || newOp.OpType == OP_POPCNT || newOp.OpType == OP_LOGICNOT {
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

func getOperatorType(op string) int {
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
	case "&^":
		return OP_BITCLEAR
	case "~":
		return OP_BITINVERSE
	case "#":
		return OP_POPCNT
	default:
		return -1
	}
}

func toNumber[T number](i interface{}) T {
	switch v := i.(type) {
	case bool:
		if v {
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
