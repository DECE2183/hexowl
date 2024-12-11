package input

import (
	"bufio"
	"io"
	"strings"
	"unicode"

	"github.com/dece2183/hexowl/v2/input/ansi"
	"github.com/dece2183/hexowl/v2/input/syntax"
	"github.com/dece2183/hexowl/v2/lexer"
	"github.com/dece2183/hexowl/v2/types"
)

type Predictable interface {
	Predict(word string) string
}

type Console struct {
	writer       io.Writer
	reader       *bufio.Reader
	history      []string
	historyIdx   int
	predictables []Predictable
}

func NewConsole(writer io.Writer, reader io.Reader, predictionList []Predictable) *Console {
	return &Console{
		writer:       writer,
		reader:       bufio.NewReader(reader),
		history:      make([]string, 1),
		historyIdx:   0,
		predictables: predictionList,
	}
}

// Searches the list of built-in and user defined variables and function for a matching completion for the word and.
//
// Returns the completed string.
func (inp *Console) Predict(word string) string {
	for _, p := range inp.predictables {
		prediction := p.Predict(word)
		if len(prediction) > 0 {
			return prediction
		}
	}

	return ""
}

// Reads input from reader awaiting new line character. Then returns resulting string.
//
// This function implements input control and history.
func (inp *Console) Prompt() (string, error) {
	inp.rewriteInputLine("", "", 1)

	var (
		cursorPos  int
		prediction string
		readRune   rune
		readErr    error
	)

	for {
		readRune, _, readErr = inp.reader.ReadRune()
		if readErr != nil && readErr != io.EOF {
			return "", readErr
		}

		if unicode.IsPrint(readRune) {
			if inp.historyIdx > 0 {
				inp.history[0] = inp.history[inp.historyIdx]
				inp.historyIdx = 0
			}

			if readRune == '(' && cursorPos == len(inp.history[0]) {
				inp.history[0] = inp.history[0] + "()"
			} else if cursorPos == len(inp.history[0]) {
				inp.history[0] = inp.history[0] + string(readRune)
			} else if readRune != ')' || inp.history[0][cursorPos] != ')' {
				inp.history[0] = inp.history[0][:cursorPos] + string(readRune) + inp.history[0][cursorPos:]
			}

			cursorPos++
			prediction = inp.getPrediction(inp.history[0], cursorPos)
			inp.rewriteInputLine(inp.history[0], prediction, cursorPos)
		} else {
			switch readRune {
			case '\t':
				// tab
				if len(prediction) > 0 {
					inp.history[0] = inp.history[0][:cursorPos] + prediction + inp.history[0][cursorPos:]
					cursorPos += len(prediction)
					if prediction[len(prediction)-1] == ')' {
						cursorPos--
					}
				}
				prediction = ""
				inp.rewriteInputLine(inp.history[0], "", cursorPos)
			case '\n', '\r':
				// new line
				inp.writer.Write([]byte("\n"))
				if inp.historyIdx != 0 {
					beg := inp.history[1:inp.historyIdx]
					end := inp.history[inp.historyIdx+1:]
					line := inp.history[inp.historyIdx]
					inp.history = append([]string{"", line}, append(beg, end...)...)
					inp.historyIdx = 0
					return line, nil
				} else {
					inp.history = append([]string{""}, inp.history...)
					return inp.history[1], nil
				}
			case '\u007F':
				// backspace
				if inp.historyIdx > 0 {
					inp.history[0] = inp.history[inp.historyIdx]
					inp.historyIdx = 0
				}

				if len(inp.history[0]) > 0 && cursorPos > 0 {
					if cursorPos < len(inp.history[0]) {
						inp.history[0] = inp.history[0][:cursorPos-1] + inp.history[0][cursorPos:]
					} else {
						inp.history[0] = inp.history[0][:cursorPos-1]
					}

					cursorPos--
					inp.rewriteInputLine(inp.history[0], "", cursorPos)
				}
			case '\u001B':
				// ANSI ESC command sequence
				var (
					cmd  rune
					args []int
				)

				cmd, args = ansi.ReadCS(inp.reader)
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
					inp.rewriteInputLine(inp.history[inp.historyIdx], "", cursorPos)
				case ansi.CUF:
					// cursor right
					if cursorPos < len(inp.history[inp.historyIdx]) {
						cursorPos++
						inp.rewriteInputLine(inp.history[inp.historyIdx], "", cursorPos)
					}
				case ansi.CUU:
					// cursor up
					if inp.historyIdx < len(inp.history)-1 {
						inp.historyIdx++
						cursorPos = len(inp.history[inp.historyIdx])
						inp.rewriteInputLine(inp.history[inp.historyIdx], "", cursorPos)
					}
				case ansi.CUD:
					// cursor down
					if inp.historyIdx > 0 {
						inp.historyIdx--
						cursorPos = len(inp.history[inp.historyIdx])
						inp.rewriteInputLine(inp.history[inp.historyIdx], "", cursorPos)
					}
				case ansi.CPL:
					// end
					cursorPos = len(inp.history[inp.historyIdx])
					inp.rewriteInputLine(inp.history[inp.historyIdx], "", cursorPos)
				case ansi.CUP:
					// home
					cursorPos = 0
					inp.rewriteInputLine(inp.history[inp.historyIdx], "", cursorPos)
				case ansi.VT:
					switch args[0] {
					case 1, 7:
						// home
						cursorPos = 0
						inp.rewriteInputLine(inp.history[inp.historyIdx], "", cursorPos)
					case 4, 8:
						// end
						cursorPos = len(inp.history[inp.historyIdx])
						inp.rewriteInputLine(inp.history[inp.historyIdx], "", cursorPos)
					case 3:
						// forward delete
						if inp.historyIdx > 0 {
							inp.history[0] = inp.history[inp.historyIdx]
							inp.historyIdx = 0
						}

						if cursorPos < len(inp.history[0]) {
							inp.history[0] = inp.history[0][:cursorPos] + inp.history[0][cursorPos+1:]
							inp.rewriteInputLine(inp.history[0], "", cursorPos)
						}
					}
				}
			}
		}
	}
}

func (inp *Console) getPrediction(str string, cpos int) string {
	var cnt int
	var wordUnderCursor types.Token
	words := lexer.Parse(str)

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

	if wordUnderCursor.Type != types.T_UNIT {
		return ""
	}

	predicted := inp.Predict(wordUnderCursor.Literal)
	if len(predicted) == 0 {
		return ""
	}

	return predicted[len(wordUnderCursor.Literal):]
}

func (inp *Console) rewriteInputLine(str, prediction string, cpos int) {
	var highlighted string

	if len(prediction) > 0 {
		highlighted = syntax.Highlight(str[:cpos]) + syntax.Colorize(prediction, syntax.C_PREDICTION) + syntax.Highlight(str[cpos:])
	} else {
		highlighted = syntax.Highlight(str)
	}

	inp.writer.Write([]byte("\n" + ansi.CreateCS(ansi.CUU, 1) + ansi.CreateCS(ansi.EL) + ">: " + highlighted))

	coffset := (len(str) + len(prediction)) - cpos
	if coffset > 0 {
		ansi.WriteCS(inp.writer, ansi.CUB, int64(coffset))
	}
}
