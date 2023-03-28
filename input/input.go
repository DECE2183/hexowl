package input

import (
	"bufio"
	"fmt"
	"io"
	"unicode"

	"github.com/dece2183/hexowl/builtin"
	"github.com/dece2183/hexowl/input/ansi"
	"github.com/dece2183/hexowl/user"
)

var (
	history    []string = make([]string, 1)
	historyIdx int      = 0
)

func GetPrediction(word string) string {
	var prediction string

	if prediction = user.PredictVariable(word); len(prediction) > 0 {
		return prediction
	}

	if prediction = user.PredictFunction(word); len(prediction) > 0 {
		return prediction
	}

	if prediction = builtin.PredictConstant(word); len(prediction) > 0 {
		return prediction
	}

	if prediction = builtin.PredictFunction(word); len(prediction) > 0 {
		return prediction
	}

	return ""
}

func Prompt(writer io.Writer, reader *bufio.Reader) (string, error) {
	rewriteInputLine(writer, "", 1)

	var (
		cursorPos int
		readRune  rune
		readErr   error
	)

	for {
		readRune, _, readErr = reader.ReadRune()
		if readErr != nil && readErr != io.EOF {
			return "", readErr
		}

		if unicode.IsPrint(readRune) {
			if historyIdx > 0 {
				history[0] = history[historyIdx]
				historyIdx = 0
			}

			if cursorPos < len(history[0]) {
				history[0] = fmt.Sprintf("%s%c%s", history[0][:cursorPos], readRune, history[0][cursorPos:])
			} else {
				history[0] = fmt.Sprintf("%s%c", history[0], readRune)
			}

			cursorPos++
			rewriteInputLine(writer, history[0], cursorPos)
		} else {
			switch readRune {
			case '\n', '\r':
				// new line
				fmt.Fprint(writer, "\n")
				if historyIdx != 0 {
					beg := history[1:historyIdx]
					end := history[historyIdx+1:]
					line := history[historyIdx]
					history = append([]string{"", line}, append(beg, end...)...)
					historyIdx = 0
					return line, nil
				} else {
					history = append([]string{""}, history...)
					return history[1], nil
				}
			case '\u007F':
				// backspace
				if historyIdx > 0 {
					history[0] = history[historyIdx]
					historyIdx = 0
				}

				if len(history[0]) > 0 && cursorPos > 0 {
					if cursorPos < len(history[0]) {
						history[0] = history[0][:cursorPos-1] + history[0][cursorPos:]
					} else {
						history[0] = history[0][:cursorPos-1]
					}

					cursorPos--
					rewriteInputLine(writer, history[0], cursorPos)
				}
			case '\u001B':
				// ANSI ESC command sequence
				var (
					cmd  rune
					args []int
				)

				cmd, args = ansi.ReadCS(reader)
				if cmd == 0 {
					continue
				}

				switch cmd {
				case ansi.CUB:
					// cursor left
					if args[0] > 0 {
						cursorPos -= int(args[0])
					} else {
						cursorPos -= 1
					}
					if cursorPos < 0 {
						cursorPos = 0
					}
					rewriteInputLine(writer, history[historyIdx], cursorPos)
				case ansi.CUF:
					// cursor right
					if cursorPos < len(history[historyIdx]) {
						cursorPos++
						rewriteInputLine(writer, history[historyIdx], cursorPos)
					}
				case ansi.CUU:
					// cursor up
					if historyIdx < len(history)-1 {
						historyIdx++
						cursorPos = len(history[historyIdx])
						rewriteInputLine(writer, history[historyIdx], cursorPos)
					}
				case ansi.CUD:
					// cursor down
					if historyIdx > 0 {
						historyIdx--
						cursorPos = len(history[historyIdx])
						rewriteInputLine(writer, history[historyIdx], cursorPos)
					}
				case ansi.CPL:
					// end
					cursorPos = len(history[historyIdx])
					rewriteInputLine(writer, history[historyIdx], cursorPos)
				case ansi.CUP:
					// home
					cursorPos = 0
					rewriteInputLine(writer, history[historyIdx], cursorPos)
				case ansi.VT:
					switch args[0] {
					case 1, 7:
						// home
						cursorPos = 0
						rewriteInputLine(writer, history[historyIdx], cursorPos)
					case 4, 8:
						// end
						cursorPos = len(history[historyIdx])
						rewriteInputLine(writer, history[historyIdx], cursorPos)
					case 3:
						// forward delete
						if historyIdx > 0 {
							history[0] = history[historyIdx]
							historyIdx = 0
						}

						if cursorPos < len(history[0]) {
							history[0] = history[0][:cursorPos] + history[0][cursorPos+1:]
							rewriteInputLine(writer, history[0], cursorPos)
						}
					}
				}
			}
		}
	}
}

func rewriteInputLine(writer io.Writer, str string, cpos int) {
	fmt.Fprintf(writer, "\n%s%s>: %s", ansi.CreateCS(ansi.CUU, 1), ansi.CreateCS(ansi.EL), str)

	coffset := len(str) - cpos
	if coffset > 0 {
		ansi.WriteCS(writer, ansi.CUB, int64(coffset))
	}
}
