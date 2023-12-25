package ansi

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

const (
	CUU = 'A' // Cursor Up
	CUD = 'B' // Cursor Down
	CUF = 'C' // Cursor Forward
	CUB = 'D' // Cursor Back
	CNL = 'E' // Cursor Next Line
	CPL = 'F' // Cursor Previous Line
	CHA = 'G' // Cursor Horizontal Absolute
	CUP = 'H' // Cursor Position
	ED  = 'J' // Erase in Display
	EL  = 'K' // Erase in Line
	SU  = 'S' // Scroll Up
	SD  = 'T' // Scroll Down
	HVP = 'f' // Horizontal Vertical Position
	SGR = 'm' // Select Graphic Rendition
	VT  = '~'
)

var (
	cmdMap = map[rune]rune{
		'A': CUU,
		'B': CUD,
		'C': CUF,
		'D': CUB,
		'E': CNL,
		'F': CPL,
		'G': CHA,
		'H': CUP,
		'J': ED,
		'K': EL,
		'S': SU,
		'T': SD,
		'f': HVP,
		'm': SGR,
		'~': VT,
	}
)

func IsESC(c rune) bool {
	return c == '\u001B'
}

func IsCSI(c rune) bool {
	return c == '['
}

func ReadCS(reader *bufio.Reader) (rune, []int) {
	str := strings.Builder{}

	c, _, err := reader.ReadRune()
	if err != nil || !IsCSI(c) {
		return 0, nil
	}

	for {
		c, _, err = reader.ReadRune()
		if cmdMap[c] != 0 || err != nil {
			stringArgs := strings.Split(str.String(), ";")
			args := make([]int, len(stringArgs))

			for i, a := range stringArgs {
				arg, _ := strconv.ParseInt(a, 10, 32)
				args[i] = int(arg)
			}

			return c, args
		}
		str.WriteRune(c)
	}
}

func WriteCS(writer io.Writer, cmd rune, args ...int64) {
	fmt.Fprint(writer, CreateCS(cmd, args...))
}

func CreateCS(cmd rune, args ...int64) string {
	var argstr string

	if len(args) > 0 {
		for _, arg := range args {
			argstr += fmt.Sprintf("%d;", arg)
		}
		argstr = argstr[:len(argstr)-1]
	}

	return fmt.Sprintf("\033[%s%c", argstr, cmd)
}
