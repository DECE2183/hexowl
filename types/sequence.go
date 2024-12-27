package types

import (
	"slices"
)

type ExecutionSequence struct {
	sequence  []interface{}
	localVars map[string]int
	userVars  map[string]bool
	userFuncs map[string]bool
}

func NewExecutionSequence() *ExecutionSequence {
	return &ExecutionSequence{
		sequence:  make([]interface{}, 0),
		localVars: make(map[string]int),
		userVars:  make(map[string]bool),
		userFuncs: make(map[string]bool),
	}
}

func (seq *ExecutionSequence) ExtractSubsequence(startPos, endPos int) *ExecutionSequence {
	s := &ExecutionSequence{
		sequence:  make([]interface{}, endPos-startPos),
		localVars: make(map[string]int),
		userVars:  make(map[string]bool),
		userFuncs: make(map[string]bool),
	}

	copy(s.sequence, seq.sequence[startPos:endPos])
	for i := range s.sequence {
		val, ok := s.sequence[i].(Value)
		if !ok {
			continue
		}
		switch val.Type {
		case V_LOCALVAR, V_VARNAME:
			switch varName := val.Value.(type) {
			case int:
				for key, index := range seq.localVars {
					if index == varName {
						s.localVars[key] = index
						break
					}
				}
			case string:
				varIndex, exists := s.localVars[varName]
				if !exists {
					varIndex = len(s.localVars)
					s.localVars[varName] = varIndex
					seq.localVars[varName] = varIndex
				}
				val.Value = varIndex
			}

			if val.Type == V_VARNAME {
				val.Type = V_LOCALVAR
				s.sequence[i] = val
			}
		case V_USERVAR:
			s.userVars[val.Value.(string)] = true
		case V_USERFUNC, V_FUNCNAME:
			s.userFuncs[val.Value.(string)] = true
			if val.Type == V_FUNCNAME {
				val.Type = V_USERFUNC
				s.sequence[i] = val
			}
		}
	}

	seq.sequence = slices.Delete(seq.sequence, startPos, endPos)
	return s
}

func (seq *ExecutionSequence) Len() int {
	return len(seq.sequence)
}

func (seq *ExecutionSequence) LocalVariablesLen() int {
	return len(seq.localVars)
}

func (seq *ExecutionSequence) GetSequence() []interface{} {
	return seq.sequence
}

func (seq *ExecutionSequence) GetValue(idx int) (Value, bool) {
	val, ok := seq.sequence[idx].(Value)
	return val, ok
}

func (seq *ExecutionSequence) SetValue(idx int, v Value) {
	seq.prepareVal(&v)
	seq.sequence[idx] = v
}

func (seq *ExecutionSequence) InsertValue(idx int, v Value) {
	seq.prepareVal(&v)
	seq.sequence = slices.Insert(seq.sequence, idx, interface{}(v))
}

func (seq *ExecutionSequence) AppendValue(v Value) {
	seq.prepareVal(&v)
	seq.sequence = append(seq.sequence, v)
}

func (seq *ExecutionSequence) GetOperator(idx int) (Operator, bool) {
	op, ok := seq.sequence[idx].(Operator)
	return op, ok
}

func (seq *ExecutionSequence) AppendOperator(op Operator) {
	seq.sequence = append(seq.sequence, op)
}

func (seq *ExecutionSequence) SetOperator(idx int, op Operator) {
	seq.sequence[idx] = op
}

func (seq *ExecutionSequence) InsertOperator(idx int, op Operator) {
	seq.sequence = slices.Insert(seq.sequence, idx, interface{}(op))
}

func (seq *ExecutionSequence) HasLocalVariable(name string) bool {
	_, ok := seq.localVars[name]
	return ok
}

func (seq *ExecutionSequence) GetLocalVariableIndex(name string) (int, bool) {
	index, ok := seq.localVars[name]
	return index, ok
}

func (seq *ExecutionSequence) HasUserVariable(name string) bool {
	_, ok := seq.userVars[name]
	return ok
}

func (seq *ExecutionSequence) HasUserFunction(name string) bool {
	_, ok := seq.userFuncs[name]
	return ok
}

func (seq *ExecutionSequence) prepareVal(v *Value) {
	if !v.Type.IsAssignable() {
		return
	}

	valName := v.Value.(string)

	switch v.Type {
	case V_LOCALVAR:
		var varIndex int
		if index, ok := seq.localVars[valName]; ok {
			varIndex = index
		} else {
			varIndex = len(seq.localVars)
			seq.localVars[valName] = varIndex
		}
		v.Value = varIndex
	case V_USERVAR:
		seq.userVars[valName] = true
	case V_USERFUNC:
		seq.userFuncs[valName] = true
	}
}
