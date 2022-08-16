package operators

import (
	"fmt"
	"math"
	"math/bits"
	"strconv"

	"github.com/dece2183/hexowl/builtin"
	"github.com/dece2183/hexowl/user"
	"github.com/dece2183/hexowl/utils"
)

type operatorType int

// Operator types
const (
	OP_NONE operatorType = iota

	OP_DECLFUNC operatorType = iota

	OP_ASSIGN    operatorType = iota
	OP_DECREMENT operatorType = iota
	OP_INCREMENT operatorType = iota
	OP_ASSIGNMUL operatorType = iota
	OP_ASSIGNDIV operatorType = iota

	OP_LOGICNOT operatorType = iota
	OP_LOGICOR  operatorType = iota
	OP_LOGICAND operatorType = iota
	OP_EQUALITY operatorType = iota
	OP_NOTEQ    operatorType = iota
	OP_MORE     operatorType = iota
	OP_LESS     operatorType = iota
	OP_MOREEQ   operatorType = iota
	OP_LESSEQ   operatorType = iota

	OP_MINUS    operatorType = iota
	OP_PLUS     operatorType = iota
	OP_MULTIPLY operatorType = iota
	OP_DIVIDE   operatorType = iota
	OP_MODULO   operatorType = iota
	OP_POWER    operatorType = iota

	OP_LEFTSHIFT  operatorType = iota
	OP_RIGHTSHIFT operatorType = iota
	OP_BITOR      operatorType = iota
	OP_BITAND     operatorType = iota
	OP_BITXOR     operatorType = iota
	OP_BITCLEAR   operatorType = iota
	OP_BITINVERSE operatorType = iota
	OP_POPCNT     operatorType = iota

	OP_FUNCARGSEP  operatorType = iota
	OP_USERFUNC    operatorType = iota
	OP_BUILTINFUNC operatorType = iota
)

type Operator struct {
	Type     operatorType
	OperandA *Operator
	OperandB *Operator
	Result   interface{}
}

var (
	operatorsPriorityList = [...]string{
		"->", "=", "-=", "+=", "*=", "/=", ",", "==", "!=", ">", "<", ">=", "<=", "||", "&&", "!", "+", "-", "*", "**", "/", "%", "<<", ">>", "|", "&", "^", "&^", "&~", "~", "#",
	}
)

func getVariable(literal string) (val interface{}, found bool) {
	val, found = user.GetVariable(literal)
	if !found {
		val, found = builtin.GetConstant(literal)
	}
	return
}

func GetType(op string) operatorType {
	switch op {
	case "->":
		return OP_DECLFUNC
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
		return OP_NOTEQ
	case ">":
		return OP_MORE
	case "<":
		return OP_LESS
	case ">=":
		return OP_MOREEQ
	case "<=":
		return OP_LESSEQ
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

func Generate(words []utils.Word) (*Operator, error) {
	var err error
	newOp := &Operator{}

	if len(words) == 0 {
		newOp.Result = uint64(0)
		return newOp, nil
	} else if len(words) == 1 {
		w := words[0]
		switch w.Type {
		case utils.W_STR:
			// Try to find variable
			val, found := getVariable(w.Literal)
			if !found {
				return nil, fmt.Errorf("there is no variable named '%s'", w.Literal)
			}
			newOp.Result = val

		case utils.W_FUNC:
			// Try to find function
			if user.HasFunction(w.Literal) {
				newOp.Type = OP_USERFUNC
				newOp.Result = w.Literal
				break
			}
			if builtin.HasFunction(w.Literal) {
				newOp.Type = OP_BUILTINFUNC
				newOp.Result = w.Literal
				break
			}
			return nil, fmt.Errorf("there is no function named '%s'", w.Literal)

		case utils.W_NUM_DEC, utils.W_NUM_HEX, utils.W_NUM_BIN:
			switch w.Type {
			case utils.W_NUM_DEC:
				newOp.Result, err = strconv.ParseFloat(w.Literal, 64)
				if err != nil {
					return nil, fmt.Errorf("unable to parse literal '%s' as number", w.Literal)
				}
			case utils.W_NUM_HEX:
				newOp.Result, err = strconv.ParseUint(w.Literal, 16, 64)
				if err != nil {
					return nil, fmt.Errorf("unable to parse literal '%s' as hex number", w.Literal)
				}
			case utils.W_NUM_BIN:
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
	var minPriorityWord *utils.Word

	bracketsCount := 0

	if words[0].Type == utils.W_CTL && words[len(words)-1].Type == utils.W_CTL {
		if words[0].Literal != "(" {
			return nil, fmt.Errorf("missing opening bracket")
		}
		if words[len(words)-1].Literal != ")" {
			return nil, fmt.Errorf("missing closing bracket")
		}

		for i := 1; i < len(words)-1; i++ {
			if words[i].Type != utils.W_CTL {
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

		if w.Type == utils.W_CTL {
			if w.Literal == "(" {
				bracketsCount++
			} else {
				bracketsCount--
			}
			continue
		} else if w.Type == utils.W_STR {
			// Function call detect
			if i+1 < len(words)-1 && words[i+1].Type == utils.W_CTL && words[i+1].Literal == "(" {
				words[i].Type = utils.W_FUNC
			}
			continue
		}

		if w.Type != utils.W_OP || bracketsCount > 0 {
			continue
		}

		prio := -1

		if GetType(w.Literal) == OP_MINUS && (i == 0 || words[i-1].Type == utils.W_OP) {
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
		if words[0].Type == utils.W_FUNC {
			if len(words) < 3 {
				return nil, fmt.Errorf("missing function '%s' arguments", words[0].Literal)
			}
			newOp.OperandA, err = Generate(words[:1])
			if err != nil {
				return nil, err
			}
			newOp.Type = newOp.OperandA.Type
			newOp.OperandA.Type = OP_NONE
			if len(words) > 3 {
				newOp.OperandB, err = Generate(words[2 : len(words)-1])
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
		newOp.Type = GetType(minPriorityWord.Literal)
		if newOp.Type < 0 {
			return nil, fmt.Errorf("unknown operator '%s'", minPriorityWord.Literal)
		}

		if newOp.Type == OP_DECLFUNC {
			// Function declaration operator
			if minPriorityIndex < 1 {
				return nil, fmt.Errorf("missing a function declaration on left side of operator '%s'", minPriorityWord.Literal)
			} else if minPriorityIndex >= len(words)-1 {
				return nil, fmt.Errorf("missing a function bidy on right side of operator '%s'", minPriorityWord.Literal)
			}
			newOp.OperandA = &Operator{
				Result: words[:minPriorityIndex],
			}
			newOp.OperandB = &Operator{
				Result: words[minPriorityIndex+1:],
			}
			return newOp, nil
		} else if newOp.Type == OP_BITINVERSE || newOp.Type == OP_POPCNT || newOp.Type == OP_LOGICNOT {
			// One side operators
			newOp.OperandA = &Operator{}
		} else if newOp.Type >= OP_ASSIGN && newOp.Type <= OP_ASSIGNDIV {
			// Assign operators
			if minPriorityIndex < 1 {
				return nil, fmt.Errorf("missing a variable on left side of operator '%s'", minPriorityWord.Literal)
			}
			lit := words[minPriorityIndex-1].Literal
			if newOp.Type > OP_ASSIGN {
				if !user.HasVariable(lit) {
					return nil, fmt.Errorf("there is no variable named '%s'", lit)
				}
			}
			newOp.OperandA = &Operator{
				Result: lit,
			}
		} else {
			newOp.OperandA, err = Generate(words[:minPriorityIndex])
			if err != nil {
				return nil, err
			}
		}
		newOp.OperandB, err = Generate(words[minPriorityIndex+1:])
		if err != nil {
			return nil, err
		}
	}

	return newOp, nil
}

func Calculate(op *Operator) (interface{}, error) {
	var err error

	if op == nil {
		return 0, nil
	}

	if op.Type == OP_NONE {
		return op.Result, nil
	} else if op.OperandA == nil || op.OperandB == nil {
		return 0, nil
	}

	op.OperandA.Result, err = Calculate(op.OperandA)
	if err != nil {
		return nil, err
	}
	op.OperandB.Result, err = Calculate(op.OperandB)
	if err != nil {
		return nil, err
	}

	switch op.Type {
	case OP_DECLFUNC:
		leftSideWords := op.OperandA.Result.([]utils.Word)
		// rightSideWords := op.OperandB.Result.([]utils.Word)
		funcName := leftSideWords[0].Literal
		newFunc := user.FuncVariant{}
		user.SetFunctionVariant(funcName, newFunc)
	case OP_ASSIGN:
		user.SetVariable(op.OperandA.Result.(string), op.OperandB.Result)
		op.Result = op.OperandB.Result
	case OP_DECREMENT:
		op.Result, _ = user.GetVariable(op.OperandA.Result.(string))
		op.Result = utils.ToNumber[float64](op.Result) - utils.ToNumber[float64](op.OperandB.Result)
		user.SetVariable(op.OperandA.Result.(string), op.Result)
	case OP_INCREMENT:
		op.Result, _ = user.GetVariable(op.OperandA.Result.(string))
		op.Result = utils.ToNumber[float64](op.Result) + utils.ToNumber[float64](op.OperandB.Result)
		user.SetVariable(op.OperandA.Result.(string), op.Result)
	case OP_ASSIGNMUL:
		op.Result, _ = user.GetVariable(op.OperandA.Result.(string))
		op.Result = utils.ToNumber[float64](op.Result) * utils.ToNumber[float64](op.OperandB.Result)
		user.SetVariable(op.OperandA.Result.(string), op.Result)
	case OP_ASSIGNDIV:
		op.Result, _ = user.GetVariable(op.OperandA.Result.(string))
		op.Result = utils.ToNumber[float64](op.Result) / utils.ToNumber[float64](op.OperandB.Result)
		user.SetVariable(op.OperandA.Result.(string), op.Result)
	case OP_LOGICNOT:
		op.Result = !utils.ToBool(op.OperandB.Result)
	case OP_LOGICOR:
		op.Result = utils.ToBool(op.OperandA.Result) || utils.ToBool(op.OperandB.Result)
	case OP_LOGICAND:
		op.Result = utils.ToBool(op.OperandA.Result) && utils.ToBool(op.OperandB.Result)
	case OP_EQUALITY:
		op.Result = utils.ToNumber[uint64](op.OperandA.Result) == utils.ToNumber[uint64](op.OperandB.Result)
	case OP_NOTEQ:
		op.Result = utils.ToNumber[uint64](op.OperandA.Result) != utils.ToNumber[uint64](op.OperandB.Result)
	case OP_MORE:
		op.Result = utils.ToNumber[uint64](op.OperandA.Result) > utils.ToNumber[uint64](op.OperandB.Result)
	case OP_LESS:
		op.Result = utils.ToNumber[uint64](op.OperandA.Result) < utils.ToNumber[uint64](op.OperandB.Result)
	case OP_MOREEQ:
		op.Result = utils.ToNumber[uint64](op.OperandA.Result) >= utils.ToNumber[uint64](op.OperandB.Result)
	case OP_LESSEQ:
		op.Result = utils.ToNumber[uint64](op.OperandA.Result) <= utils.ToNumber[uint64](op.OperandB.Result)
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
			f, _ := builtin.GetFunction(op.OperandA.Result.(string))
			op.Result, err = f.Exec(op.OperandB.Result.([]interface{})...)
		default:
			f, _ := builtin.GetFunction(op.OperandA.Result.(string))
			op.Result, err = f.Exec(op.OperandB.Result)
		}
	case OP_USERFUNC:
		op.Result = 0
	}

	return op.Result, err
}
