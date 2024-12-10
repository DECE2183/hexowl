//go:build !nohighlight
// +build !nohighlight

package syntax

import (
	"strings"

	"github.com/dece2183/hexowl/v2/calculator/lexer"
	"github.com/dece2183/hexowl/v2/calculator/types"
)

func Highlight(str string) (out string) {
	tokens := lexer.Parse(str)

	for _, t := range tokens {
		pos := strings.Index(str, t.Literal)
		if pos > 0 {
			out += colors[types.T_OP]
			out += str[:pos]
			str = str[pos:]
		}

		clr, ok := colors[t.Type]
		if !ok {
			return str
		}

		out += clr + t.Literal

		if len(str) > 0 {
			str = str[len(t.Literal):]
		}
	}

	out += str + colors[C_NORMAL]
	return
}

func Colorize(word string, wordType types.TokenType) string {
	clr, ok := colors[wordType]
	if !ok {
		return word
	}
	return clr + word + colors[C_NORMAL]
}
