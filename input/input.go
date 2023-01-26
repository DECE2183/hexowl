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
	history []string
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
		cursorPos  int
		readString string
		readRune   rune
		readErr    error
	)

	for {
		readRune, _, readErr = reader.ReadRune()
		if readErr != nil && readErr != io.EOF {
			return "", readErr
		}

		if unicode.IsPrint(readRune) {
			if cursorPos < len(readString) {
				readString = fmt.Sprintf("%s%c%s", readString[:cursorPos], readRune, readString[cursorPos:])
			} else {
				readString = fmt.Sprintf("%s%c", readString, readRune)
			}

			cursorPos++
			rewriteInputLine(writer, readString, cursorPos)
		} else {
			switch readRune {
			case '\n':
				// new line
				fmt.Fprint(writer, "\n")
				return readString, nil
			case '\u007F':
				// backspace
				if len(readString) > 0 && cursorPos > 0 {
					if cursorPos < len(readString) {
						readString = readString[:cursorPos-1] + readString[cursorPos:]
					} else {
						readString = readString[:cursorPos-1]
					}

					cursorPos--
					rewriteInputLine(writer, readString, cursorPos)
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
					if cursorPos > 0 {
						cursorPos--
						rewriteInputLine(writer, readString, cursorPos)
					}
				case ansi.CUF:
					if cursorPos < len(readString) {
						cursorPos++
						rewriteInputLine(writer, readString, cursorPos)
					}
				}
			}
		}
	}
}

func rewriteInputLine(writer io.Writer, str string, cpus int) {
	fmt.Fprintf(writer, "\n\033[1A\033[K>: %s", str)

	coffset := len(str) - cpus
	if coffset > 0 {
		fmt.Fprintf(writer, "\033[%dD", coffset)
	}
}
