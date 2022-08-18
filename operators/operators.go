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

	OP_FUNCARGSEP operatorType = iota
	OP_SEQUENCE   operatorType = iota

	OP_LOCALVAR    operatorType = iota
	OP_USERVAR     operatorType = iota
	OP_CONSTANT    operatorType = iota
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
		"->", ";", "=", "-=", "+=", "*=", "/=", ",", "||", "&&", "==", "!=", "!", ">", "<", ">=", "<=", "+", "-", "*", "**", "/", "%", "<<", ">>", "|", "&", "^", "&^", "&~", "~", "#",
	}
)

func getLocalVariable(localVars map[string]interface{}, literal string) (val interface{}, found bool) {
	if localVars == nil {
		return nil, false
	}
	val, found = localVars[literal]
	return
}

func copyLocalVars(dest, localVars map[string]interface{}) {
	for k, v := range localVars {
		switch f := v.(type) {
		case []interface{}:
			var newSlice []interface{}
			copy(newSlice, f)
			dest[k] = newSlice
		case map[string]interface{}:
			newMap := make(map[string]interface{})
			copyLocalVars(newMap, f)
			dest[k] = newMap
		default:
			dest[k] = v
		}
	}
}

func execUserFunc(f user.Func, args []interface{}, localVars map[string]interface{}) (result interface{}, err error) {
	for _, variant := range f.Variants {
		argsLen := len(args)
		argNames := variant.ArgNames()
		argNamesLen := len(argNames)
		if argNamesLen != argsLen {
			if argNamesLen == 0 {
				if argsLen > 1 || args[0] != nil {
					continue
				}
			} else if argsLen < argNamesLen || argNames[argNamesLen-1] != "@" {
				continue
			}
		}

		argMap := make(map[string]interface{})
		copyLocalVars(argMap, localVars)

		for pos, name := range argNames {
			if name == "@" {
				argMap[name] = args[pos:]
				break
			} else {
				argMap[name] = args[pos]
			}
		}
		argOperators, err := Generate(variant.Args, argMap)
		if err != nil {
			continue
		}
		result, err := Calculate(argOperators, argMap)
		if err != nil {
			continue
		}

		argsCompatible := true

		switch r := result.(type) {
		case []interface{}:
			for _, val := range r {
				switch v := val.(type) {
				case bool:
					if !v {
						argsCompatible = false
						break
					}
				}

				if !argsCompatible {
					break
				}
			}
		default:
			switch val := r.(type) {
			case bool:
				if !val {
					argsCompatible = false
				}
			}
		}

		if !argsCompatible {
			continue
		}
		bodyOperators, err := Generate(variant.Body, argMap)
		if err != nil {
			result = nil
			return result, err
		}

		return Calculate(bodyOperators, argMap)
	}

	result = nil
	err = fmt.Errorf("variation not found")

	return
}

func GetType(op string) operatorType {
	switch op {
	case ";":
		return OP_SEQUENCE
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

func Generate(words []utils.Word, localVars map[string]interface{}) (*Operator, error) {
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
			_, found := getLocalVariable(localVars, w.Literal)
			if found {
				newOp.Type = OP_LOCALVAR
				newOp.Result = w.Literal
			} else if user.HasVariable(w.Literal) {
				newOp.Type = OP_USERVAR
				newOp.Result = w.Literal
			} else if builtin.HasConstant(w.Literal) {
				newOp.Type = OP_CONSTANT
				newOp.Result = w.Literal
			} else {
				return nil, fmt.Errorf("there is no variable named '%s'", w.Literal)
			}

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
			newOp.OperandA, err = Generate(words[:1], localVars)
			if err != nil {
				return nil, err
			}
			newOp.Type = newOp.OperandA.Type
			newOp.OperandA.Type = OP_NONE
			if len(words) > 3 {
				newOp.OperandB, err = Generate(words[2:len(words)-1], localVars)
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
			if minPriorityIndex < 3 {
				return nil, fmt.Errorf("missing a function declaration on left side of operator '%s'", minPriorityWord.Literal)
			} else if minPriorityIndex >= len(words)-1 {
				return nil, fmt.Errorf("missing a function bidy on right side of operator '%s'", minPriorityWord.Literal)
			}

			// Find brackets to determine arguments
			bracketsCount = 0
			lastBracketIndex := -1
			if words[1].Literal != "(" {
				return nil, fmt.Errorf("wrong function declaration syntax, missing '('")
			} else {
				for i := 1; i < minPriorityIndex; i++ {
					if words[i].Type == utils.W_CTL {
						if words[i].Literal == "(" {
							bracketsCount++
						} else {
							lastBracketIndex = i
							bracketsCount--
						}
					} else {
						continue
					}
				}
				if bracketsCount > 0 || lastBracketIndex < 0 {
					return nil, fmt.Errorf("wrong function declaration syntax, missing ')'")
				}
			}

			newOp.OperandA = &Operator{
				Result: words[:lastBracketIndex+1],
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

			_, foundLocal := getLocalVariable(localVars, lit)
			if foundLocal {
				newOp.OperandA = &Operator{
					Type:   OP_LOCALVAR,
					Result: lit,
				}
			} else if user.HasVariable(lit) {
				newOp.OperandA = &Operator{
					Type:   OP_USERVAR,
					Result: lit,
				}
			} else {
				newOp.OperandA = &Operator{
					Result: lit,
				}
			}

			if newOp.Type > OP_ASSIGN && newOp.OperandA.Type == OP_NONE {
				return nil, fmt.Errorf("there is no user variable named '%s'", lit)
			}
		} else {
			newOp.OperandA, err = Generate(words[:minPriorityIndex], localVars)
			if err != nil {
				return nil, err
			}
		}
		newOp.OperandB, err = Generate(words[minPriorityIndex+1:], localVars)
		if err != nil {
			return nil, err
		}
	}

	return newOp, nil
}

func Calculate(op *Operator, localVars map[string]interface{}) (interface{}, error) {
	var err error

	if op == nil {
		return nil, nil
	}

	if op.OperandA == nil && op.OperandB == nil {
		switch op.Type {
		case OP_NONE:
			return op.Result, nil
		case OP_LOCALVAR:
			op.OperandA = &Operator{
				Result: op.Result.(string),
			}
			op.Result, _ = getLocalVariable(localVars, op.Result.(string))
		case OP_USERVAR:
			op.OperandA = &Operator{
				Result: op.Result.(string),
			}
			op.Result, _ = user.GetVariable(op.Result.(string))
		case OP_CONSTANT:
			op.OperandA = &Operator{
				Result: op.Result.(string),
			}
			op.Result, _ = builtin.GetConstant(op.Result.(string))
		case OP_USERFUNC:
			op.OperandA = &Operator{
				Result: op.Result.(string),
			}
			op.Result, _ = user.GetFunction(op.Result.(string))
		case OP_BUILTINFUNC:
			op.OperandA = &Operator{
				Result: op.Result.(string),
			}
			op.Result, _ = builtin.GetFunction(op.Result.(string))
		default:
			return nil, fmt.Errorf("missing operands")
		}

		return op.Result, nil
	}

	if op.OperandA != nil {
		op.OperandA.Result, err = Calculate(op.OperandA, localVars)
		if err != nil {
			return nil, err
		}
	}
	if op.OperandB != nil {
		op.OperandB.Result, err = Calculate(op.OperandB, localVars)
		if err != nil {
			return nil, err
		}
	}

	switch op.Type {
	case OP_SEQUENCE:
		if op.OperandB != nil {
			op.Result = op.OperandB.Result
		} else {
			return nil, nil
		}
	case OP_DECLFUNC:
		leftSideWords := op.OperandA.Result.([]utils.Word)
		rightSideWords := op.OperandB.Result.([]utils.Word)
		funcName := leftSideWords[0].Literal
		newFunc := user.FuncVariant{
			Args: leftSideWords[2 : len(leftSideWords)-1],
			Body: rightSideWords,
		}
		user.SetFunctionVariant(funcName, newFunc)
	case OP_ASSIGN, OP_DECREMENT, OP_INCREMENT, OP_ASSIGNMUL, OP_ASSIGNDIV:
		switch op.Type {
		case OP_ASSIGN:
			op.Result = op.OperandB.Result
		case OP_DECREMENT:
			op.Result = utils.ToNumber[float64](op.OperandA.Result) - utils.ToNumber[float64](op.OperandB.Result)
		case OP_INCREMENT:
			op.Result = utils.ToNumber[float64](op.OperandA.Result) + utils.ToNumber[float64](op.OperandB.Result)
		case OP_ASSIGNMUL:
			op.Result = utils.ToNumber[float64](op.OperandA.Result) * utils.ToNumber[float64](op.OperandB.Result)
		case OP_ASSIGNDIV:
			op.Result = utils.ToNumber[float64](op.OperandA.Result) / utils.ToNumber[float64](op.OperandB.Result)
		}
		switch op.OperandA.Type {
		case OP_NONE:
			user.SetVariable(op.OperandA.Result.(string), op.Result)
		case OP_USERVAR:
			user.SetVariable(op.OperandA.OperandA.Result.(string), op.Result)
		case OP_LOCALVAR:
			localVars[op.OperandA.OperandA.Result.(string)] = op.Result
		default:
			return nil, fmt.Errorf("tru to assign non user variable")
		}
	case OP_LOGICNOT:
		op.Result = !utils.ToBool(op.OperandB.Result)
	case OP_LOGICOR:
		op.Result = utils.ToBool(op.OperandA.Result) || utils.ToBool(op.OperandB.Result)
	case OP_LOGICAND:
		op.Result = utils.ToBool(op.OperandA.Result) && utils.ToBool(op.OperandB.Result)
	case OP_EQUALITY:
		op.Result = utils.ToNumber[float64](op.OperandA.Result) == utils.ToNumber[float64](op.OperandB.Result)
	case OP_NOTEQ:
		op.Result = utils.ToNumber[float64](op.OperandA.Result) != utils.ToNumber[float64](op.OperandB.Result)
	case OP_MORE:
		op.Result = utils.ToNumber[float64](op.OperandA.Result) > utils.ToNumber[float64](op.OperandB.Result)
	case OP_LESS:
		op.Result = utils.ToNumber[float64](op.OperandA.Result) < utils.ToNumber[float64](op.OperandB.Result)
	case OP_MOREEQ:
		op.Result = utils.ToNumber[float64](op.OperandA.Result) >= utils.ToNumber[float64](op.OperandB.Result)
	case OP_LESSEQ:
		op.Result = utils.ToNumber[float64](op.OperandA.Result) <= utils.ToNumber[float64](op.OperandB.Result)
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
		f, _ := builtin.GetFunction(op.OperandA.Result.(string))

		switch op.OperandB.Result.(type) {
		case []interface{}:
			op.Result, err = f.Exec(op.OperandB.Result.([]interface{})...)
		default:
			op.Result, err = f.Exec(op.OperandB.Result)
		}
	case OP_USERFUNC:
		var args []interface{}
		fname := op.OperandA.Result.(string)
		f, _ := user.GetFunction(fname)

		switch op.OperandB.Result.(type) {
		case []interface{}:
			args = op.OperandB.Result.([]interface{})
		default:
			args = []interface{}{op.OperandB.Result}
		}

		op.Result, err = execUserFunc(f, args, localVars)
		if err != nil {
			return nil, fmt.Errorf("unable to find proper '%s' function variation for argsuments: %v", fname, args)
		}
	}

	// fmt.Println(op.Result)
	return op.Result, err
}
