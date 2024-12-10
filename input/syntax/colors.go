//go:build !nohighlight
// +build !nohighlight

package syntax

import (
	"github.com/dece2183/hexowl/v2/calculator/types"
	"github.com/dece2183/hexowl/v2/input/ansi"
)

const (
	// Normal text color
	C_NORMAL = types.T_STR + 1 + iota
	// Input prediction text color
	C_PREDICTION
	// Error text color
	C_ERROR
)

var colors = map[types.TokenType]string{
	types.T_NONE:    ansi.CreateCS(ansi.SGR, 38, 5, 244),
	types.T_NUM_SCI: ansi.CreateCS(ansi.SGR, 38, 5, 32),
	types.T_NUM_DEC: ansi.CreateCS(ansi.SGR, 38, 5, 32),
	types.T_NUM_HEX: ansi.CreateCS(ansi.SGR, 38, 5, 32),
	types.T_NUM_BIN: ansi.CreateCS(ansi.SGR, 38, 5, 32),
	types.T_UNIT:    ansi.CreateCS(ansi.SGR, 38, 5, 209),
	types.T_OP:      ansi.CreateCS(ansi.SGR, 37),
	types.T_CTL:     ansi.CreateCS(ansi.SGR, 37),
	types.T_STR:     ansi.CreateCS(ansi.SGR, 38, 5, 71),

	C_NORMAL:     ansi.CreateCS(ansi.SGR, 37),
	C_PREDICTION: ansi.CreateCS(ansi.SGR, 38, 5, 244),
	C_ERROR:      ansi.CreateCS(ansi.SGR, 38, 2, 209, 84, 84),
}
