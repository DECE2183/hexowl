package types

import "slices"

type ExecutionSequence struct {
	sequence  []interface{}
	localVars map[string]ValueType
}

func NewExecutionSequence() *ExecutionSequence {
	return &ExecutionSequence{
		sequence:  make([]interface{}, 0),
		localVars: make(map[string]ValueType),
	}
}

func (seq *ExecutionSequence) Len() int {
	return len(seq.sequence)
}

func (seq *ExecutionSequence) GetValue(idx int) (Value, bool) {
	val, ok := seq.sequence[idx].(Value)
	return val, ok
}

func (seq *ExecutionSequence) GetOperator(idx int) (OperatorType, bool) {
	op, ok := seq.sequence[idx].(OperatorType)
	return op, ok
}

func (seq *ExecutionSequence) GetSequence() []interface{} {
	return seq.sequence
}

func (seq *ExecutionSequence) AppendValue(v Value) {
	seq.sequence = append(seq.sequence, v)
}

func (seq *ExecutionSequence) InsertValue(idx int, v Value) {
	seq.sequence = slices.Insert(seq.sequence, idx, interface{}(v))
}

func (seq *ExecutionSequence) AppendOperator(op OperatorType) {
	seq.sequence = append(seq.sequence, op)
}

func (seq *ExecutionSequence) InsertOperator(idx int, op OperatorType) {
	seq.sequence = slices.Insert(seq.sequence, idx, interface{}(op))
}

func (seq *ExecutionSequence) SetLocalVariable(v Value) {
	seq.localVars[v.Value.(string)] = v.Type
}

func (seq *ExecutionSequence) GetLocalVariable(name string) (ValueType, bool) {
	t, ok := seq.localVars[name]
	return t, ok
}

func (seq *ExecutionSequence) HasLocalVariable(name string) bool {
	_, ok := seq.localVars[name]
	return ok
}
