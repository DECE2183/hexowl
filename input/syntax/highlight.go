//go:build !nohighlight
// +build !nohighlight

package syntax

import (
	"strings"

	"github.com/dece2183/hexowl/utils"
)

func Highlight(str string) (out string) {
	words := utils.ParsePrompt(str)

	for _, w := range words {
		pos := strings.Index(str, w.Literal)
		if pos > 0 {
			out += colors[utils.W_OP]
			out += str[:pos]
			str = str[pos:]
		}

		clr, ok := colors[w.Type]
		if !ok {
			return str
		}

		out += clr + w.Literal

		if len(str) > 0 {
			str = str[len(w.Literal):]
		}
	}

	out += str + colors[C_NORMAL]
	return
}

func Colorize(word string, wordType utils.WordType) string {
	clr, ok := colors[wordType]
	if !ok {
		return word
	}
	return clr + word + colors[C_NORMAL]
}
