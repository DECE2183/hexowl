//go:build !nohighlight
// +build !nohighlight

package syntax

import (
	"github.com/dece2183/hexowl/input/ansi"
	"github.com/dece2183/hexowl/utils"
)

const (
	// Normal text color
	C_NORMAL = utils.W_STR + 1 + iota
	// Input prediction text color
	C_PREDICTION
	// Error text color
	C_ERROR
)

var colors = map[utils.WordType]string{
	utils.W_NONE:    ansi.CreateCS(ansi.SGR, 38, 5, 244),
	utils.W_NUM_SCI: ansi.CreateCS(ansi.SGR, 38, 5, 32),
	utils.W_NUM_DEC: ansi.CreateCS(ansi.SGR, 38, 5, 32),
	utils.W_NUM_HEX: ansi.CreateCS(ansi.SGR, 38, 5, 32),
	utils.W_NUM_BIN: ansi.CreateCS(ansi.SGR, 38, 5, 32),
	utils.W_UNIT:    ansi.CreateCS(ansi.SGR, 38, 5, 209),
	utils.W_OP:      ansi.CreateCS(ansi.SGR, 37),
	utils.W_CTL:     ansi.CreateCS(ansi.SGR, 37),
	utils.W_FUNC:    ansi.CreateCS(ansi.SGR, 38, 5, 230),
	utils.W_STR:     ansi.CreateCS(ansi.SGR, 38, 5, 71),

	C_NORMAL:     ansi.CreateCS(ansi.SGR, 37),
	C_PREDICTION: ansi.CreateCS(ansi.SGR, 38, 5, 244),
	C_ERROR:      ansi.CreateCS(ansi.SGR, 38, 2, 209, 84, 84),
}
