package types

type ExecutionSequence struct {
	sequence  []interface{}
	localVars map[string]bool
}

func NewExecutionSequence() *ExecutionSequence {
	return &ExecutionSequence{
		sequence:  make([]interface{}, 0),
		localVars: make(map[string]bool),
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

func (seq *ExecutionSequence) AppendOperator(op OperatorType) {
	seq.sequence = append(seq.sequence, op)
}

func (seq *ExecutionSequence) SetLocalVariable(name string) {
	seq.localVars[name] = true
}

func (seq *ExecutionSequence) HasLocalVariable(name string) bool {
	return seq.localVars[name]
}
