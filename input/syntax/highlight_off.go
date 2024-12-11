//go:build nohighlight
// +build nohighlight

package syntax

import "github.com/dece2183/hexowl/v2/types"

const (
	// Normal text color
	C_NORMAL = types.T_STR + 1 + iota
	// Input prediction text color
	C_PREDICTION
	// Error text color
	C_ERROR
)

func Highlight(str string) (out string) {
	return str
}

func Colorize(word string, wordType types.TokenType) string {
	return word
}
