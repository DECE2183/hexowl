package input

import (
	"bufio"
	"io"
	"strings"
	"unicode"

	"github.com/dece2183/hexowl/builtin"
	"github.com/dece2183/hexowl/input/ansi"
	"github.com/dece2183/hexowl/input/syntax"
	"github.com/dece2183/hexowl/user"
	"github.com/dece2183/hexowl/utils"
)

var (
	history    []string = make([]string, 1)
	historyIdx int      = 0
)

// Searches the list of built-in and user defined variables and function for a matching completion for the word and.
//
// Returns the completed string.
func Predict(word string) string {
	var prediction string

	if prediction = user.Predict(word); len(prediction) > 0 {
		return prediction
	}

	if prediction = builtin.Predict(word); len(prediction) > 0 {
		return prediction
	}

	return ""
}

// Reads input from reader awaiting new line character. Then returns resulting string.
//
// This function implements input control and history.
func Prompt(writer io.Writer, reader *bufio.Reader) (string, error) {
	rewriteInputLine(writer, "", "", 1)

	var (
		cursorPos  int
		prediction string
		readRune   rune
		readErr    error
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

			if readRune == '(' && cursorPos == len(history[0]) {
				history[0] = history[0] + "()"
			} else if cursorPos == len(history[0]) {
				history[0] = history[0] + string(readRune)
			} else if readRune != ')' || history[0][cursorPos] != ')' {
				history[0] = history[0][:cursorPos] + string(readRune) + history[0][cursorPos:]
			}

			cursorPos++
			prediction = getPrediction(history[0], cursorPos)
			rewriteInputLine(writer, history[0], prediction, cursorPos)
		} else {
			switch readRune {
			case '\t':
				// tab
				if len(prediction) > 0 {
					history[0] = history[0][:cursorPos] + prediction + history[0][cursorPos:]
					cursorPos += len(prediction)
					if prediction[len(prediction)-1] == ')' {
						cursorPos--
					}
				}
				prediction = ""
				rewriteInputLine(writer, history[0], "", cursorPos)
			case '\n', '\r':
				// new line
				writer.Write([]byte("\n"))
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
					rewriteInputLine(writer, history[0], "", cursorPos)
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
					rewriteInputLine(writer, history[historyIdx], "", cursorPos)
				case ansi.CUF:
					// cursor right
					if cursorPos < len(history[historyIdx]) {
						cursorPos++
						rewriteInputLine(writer, history[historyIdx], "", cursorPos)
					}
				case ansi.CUU:
					// cursor up
					if historyIdx < len(history)-1 {
						historyIdx++
						cursorPos = len(history[historyIdx])
						rewriteInputLine(writer, history[historyIdx], "", cursorPos)
					}
				case ansi.CUD:
					// cursor down
					if historyIdx > 0 {
						historyIdx--
						cursorPos = len(history[historyIdx])
						rewriteInputLine(writer, history[historyIdx], "", cursorPos)
					}
				case ansi.CPL:
					// end
					cursorPos = len(history[historyIdx])
					rewriteInputLine(writer, history[historyIdx], "", cursorPos)
				case ansi.CUP:
					// home
					cursorPos = 0
					rewriteInputLine(writer, history[historyIdx], "", cursorPos)
				case ansi.VT:
					switch args[0] {
					case 1, 7:
						// home
						cursorPos = 0
						rewriteInputLine(writer, history[historyIdx], "", cursorPos)
					case 4, 8:
						// end
						cursorPos = len(history[historyIdx])
						rewriteInputLine(writer, history[historyIdx], "", cursorPos)
					case 3:
						// forward delete
						if historyIdx > 0 {
							history[0] = history[historyIdx]
							historyIdx = 0
						}

						if cursorPos < len(history[0]) {
							history[0] = history[0][:cursorPos] + history[0][cursorPos+1:]
							rewriteInputLine(writer, history[0], "", cursorPos)
						}
					}
				}
			}
		}
	}
}

func getPrediction(str string, cpos int) string {
	var cnt int
	var wordUnderCursor utils.Word
	words := utils.ParsePrompt(str)

	for _, w := range words {
		pos := strings.Index(str, w.Literal)
		if pos > 0 {
			cnt += pos
			str = str[pos:]
		}

		cnt += len(w.Literal)
		if cnt == cpos {
			wordUnderCursor = w
			break
		} else if cnt > cpos {
			return ""
		}

		if len(str) > 0 {
			str = str[len(w.Literal):]
		}
	}

	if wordUnderCursor.Type != utils.W_UNIT {
		return ""
	}

	predicted := Predict(wordUnderCursor.Literal)
	if len(predicted) == 0 {
		return ""
	}

	return predicted[len(wordUnderCursor.Literal):]
}

func rewriteInputLine(writer io.Writer, str, prediction string, cpos int) {
	var highlighted string

	if len(prediction) > 0 {
		highlighted = syntax.Highlight(str[:cpos]) + syntax.Colorize(prediction, syntax.C_PREDICTION) + syntax.Highlight(str[cpos:])
	} else {
		highlighted = syntax.Highlight(str)
	}

	writer.Write([]byte("\n" + ansi.CreateCS(ansi.CUU, 1) + ansi.CreateCS(ansi.EL) + ">: " + highlighted))

	coffset := (len(str) + len(prediction)) - cpos
	if coffset > 0 {
		ansi.WriteCS(writer, ansi.CUB, int64(coffset))
	}
}
