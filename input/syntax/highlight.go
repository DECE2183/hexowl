//go:build !nohighlight
// +build !nohighlight

package syntax

import (
	"strings"

	"github.com/dece2183/hexowl/v2/lexer"
	"github.com/dece2183/hexowl/v2/types"
)

func Highlight(str string) (out string) {
	tokens := lexer.Parse(str)

	var clr string
	var ok bool

	for _, t := range tokens {
		switch t.Type {
		case types.T_NUM_HEX:
			t.Literal = "0x" + t.Literal
		case types.T_NUM_BIN:
			t.Literal = "0b" + t.Literal
		}

		pos := strings.Index(str, t.Literal)
		if pos > 0 {
			out += colors[types.T_OP]
			out += str[:pos]
			str = str[pos:]
		}

		clr, ok = colors[t.Type]
		if !ok {
			out += colors[C_NORMAL]
		} else {
			out += clr
		}

		out += t.Literal
		if len(str) > 0 {
			str = str[len(t.Literal):]
		}
	}

	if len(clr) > 0 {
		out += colors[C_NORMAL]
	}

	out += str
	return
}

func Colorize(word string, wordType types.TokenType) string {
	clr, ok := colors[wordType]
	if !ok {
		return word
	}
	return clr + word + colors[C_NORMAL]
}
