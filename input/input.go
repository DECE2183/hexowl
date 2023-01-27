package input

import (
	"bufio"
	"fmt"
	"io"
	"unicode"

	"github.com/dece2183/hexowl/builtin"
	"github.com/dece2183/hexowl/user"
	"github.com/jcorbin/anansi/ansi"
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
			case '\n':
				// new line
				fmt.Fprint(writer, "\n")
				if historyIdx == 0 {
					history = append([]string{""}, history...)
					return history[1], nil
				} else {
					defer func() { historyIdx = 0 }()
					return history[historyIdx], nil
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
					b    byte
					bseq []byte
					esc  ansi.Escape
					// escArgs []byte
				)

				b, readErr = reader.ReadByte()
				if readErr != nil {
					continue
				}

				esc = ansi.ESC(b)
				if !esc.IsEscape() {
					continue
				}

				bseq = make([]byte, esc.Size()-1)

				_, readErr = reader.Read(bseq)
				if readErr != nil {
					continue
				}

				bseq = append([]byte{byte(readRune), b}, bseq...)

				esc, _, _ = ansi.DecodeEscape(bseq)
				switch esc {
				case ansi.CUB:
					// cursor left
					if cursorPos > 0 {
						cursorPos--
						rewriteInputLine(writer, history[historyIdx], cursorPos)
					}
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
				}
			}
		}
	}
}

func rewriteInputLine(writer io.Writer, str string, cpos int) {
	fmt.Fprintf(writer, "\n\033[1A\033[K>: %s", str)

	coffset := len(str) - cpos
	if coffset > 0 {
		fmt.Fprintf(writer, "\033[%dD", coffset)
	}
}
