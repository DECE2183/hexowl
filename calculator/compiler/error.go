package compiler

import (
	"fmt"

	"github.com/dece2183/hexowl/v2/calculator/types"
)

type CompileError struct {
	Token       types.Token
	Pos         int
	Description string
}

func NewCompileError(token types.Token, pos int, format string, a ...any) *CompileError {
	return &CompileError{
		Token:       token,
		Pos:         pos,
		Description: fmt.Sprintf(format, a...),
	}
}

func (err *CompileError) Error() string {
	return err.Description
}
