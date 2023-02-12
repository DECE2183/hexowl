package operators

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/dece2183/hexowl/builtin"
	"github.com/dece2183/hexowl/user"
	"github.com/dece2183/hexowl/utils"
)

type operatorType int

type Operator struct {
	Type     operatorType
	OperandA *Operator
	OperandB *Operator
	Result   interface{}
}

func getLocalVariable(localVars map[string]interface{}, literal string) (val interface{}, found bool) {
	if localVars == nil {
		return nil, false
	}
	val, found = localVars[literal]
	return
}

func execUserFunc(f user.Func, args []interface{}) (result interface{}, err error) {
	var lasterr error

	for vi, variant := range f.Variants {
		argsLen := len(args)
		argNames := variant.ArgNames()
		argNamesLen := len(argNames)
		if argNamesLen != argsLen {
			if argNamesLen == 0 {
				if argsLen > 1 || args[0] != nil {
					lasterr = fmt.Errorf("expected 0 args but got %d (#%d)", argsLen, vi)
					continue
				}
			} else if argsLen < argNamesLen || argNames[argNamesLen-1] != "@" {
				lasterr = fmt.Errorf("expected %d args but got %d (#%d)", argNamesLen, argsLen, vi)
				continue
			}
		}

		argMap := make(map[string]interface{})
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
			lasterr = fmt.Errorf("%s (#%d)", err, vi)
			continue
		}
		result, err := Calculate(argOperators, argMap)
		if err != nil {
			lasterr = fmt.Errorf("%s (#%d)", err, vi)
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
			lasterr = fmt.Errorf("args not compatible (#%d)", vi)
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
	err = lasterr

	return
}

func GetType(op string) operatorType {
	t, ok := opStringRepresent[op]
	if ok {
		return t
	} else {
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
		case utils.W_UNIT:
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
			} else if user.HasFunction(w.Literal) {
				newOp.Type = OP_USERFUNC
				newOp.Result = w.Literal
			} else if builtin.HasFunction(w.Literal) {
				newOp.Type = OP_BUILTINFUNC
				newOp.Result = w.Literal
			} else {
				return nil, fmt.Errorf("there is no variable named '%s'", w.Literal)
			}

		case utils.W_FUNC:
			// Try to find function
			v, found := getLocalVariable(localVars, w.Literal)
			if found || user.HasVariable(w.Literal) {
				if !found {
					v, _ = user.GetVariable(w.Literal)
				}
				switch fname := v.(type) {
				case string:
					if user.HasFunction(fname) {
						newOp.Type = OP_USERFUNC
						newOp.Result = fname
						return newOp, nil
					} else if builtin.HasFunction(fname) {
						newOp.Type = OP_BUILTINFUNC
						newOp.Result = fname
						return newOp, nil
					}
				}
			}
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

		case utils.W_NUM_DEC, utils.W_NUM_HEX, utils.W_NUM_BIN, utils.W_NUM_SCI:
			switch w.Type {
			case utils.W_NUM_SCI:
				num := strings.Split(w.Literal, "e")
				var mantisse, order float64
				mantisse, err = strconv.ParseFloat(strings.ReplaceAll(num[0], "_", ""), 64)
				if err != nil {
					return nil, fmt.Errorf("unable to parse mantisse part of literal '%s'", w.Literal)
				}
				order, err = strconv.ParseFloat(strings.ReplaceAll(num[1], "_", ""), 64)
				if err != nil {
					return nil, fmt.Errorf("unable to parse order part of literal '%s'", w.Literal)
				}
				newOp.Result = mantisse * math.Pow(10, order)
			case utils.W_NUM_DEC:
				newOp.Result, err = strconv.ParseFloat(strings.ReplaceAll(w.Literal, "_", ""), 64)
				if err != nil {
					return nil, fmt.Errorf("unable to parse literal '%s' as number", w.Literal)
				}
			case utils.W_NUM_HEX:
				newOp.Result, err = strconv.ParseUint(strings.ReplaceAll(w.Literal, "_", ""), 16, 64)
				if err != nil {
					return nil, fmt.Errorf("unable to parse literal '%s' as hex number", w.Literal)
				}
			case utils.W_NUM_BIN:
				newOp.Result, err = strconv.ParseUint(strings.ReplaceAll(w.Literal, "_", ""), 2, 64)
				if err != nil {
					return nil, fmt.Errorf("unable to parse literal '%s' as bin number", w.Literal)
				}
			}

		case utils.W_STR:
			newOp.Result = w.Literal
		}

		return newOp, nil
	}

	minPriority := OP_COUNT
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
		} else if w.Type == utils.W_UNIT {
			// Function call detect
			if i+1 < len(words)-1 && words[i+1].Type == utils.W_CTL && words[i+1].Literal == "(" {
				words[i].Type = utils.W_FUNC
			}
			continue
		}

		if w.Type != utils.W_OP || bracketsCount > 0 {
			continue
		}

		prio := GetType(w.Literal)

		if prio == OP_MINUS && (i == 0 || words[i-1].Type == utils.W_OP) {
			// If it is a single minus operator give it the max prioriy
			prio = OP_COUNT
		}

		if prio <= 0 {
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

			if newOp.OperandA.Type == OP_NONE {
				if newOp.Type > OP_LOCALASSIGN {
					return nil, fmt.Errorf("there is no user variable named '%s'", lit)
				} else {
					localVars[lit] = nil
				}
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
				Result: true,
			}
			op.Result = op.Result.(string)
		case OP_BUILTINFUNC:
			op.OperandA = &Operator{
				Result: true,
			}
			op.Result = op.Result.(string)
		default:
			return nil, fmt.Errorf("missing operands")
		}
		return op.Result, nil
	} else {
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
	}

	return opDoAction(op, localVars)
}
