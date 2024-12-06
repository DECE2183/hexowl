package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	_ "github.com/dece2183/hexowl/builtin/default_system"
	"github.com/dece2183/hexowl/input"
	"github.com/dece2183/hexowl/input/syntax"
	"github.com/dece2183/hexowl/input/terminal"
	"github.com/dece2183/hexowl/operators"
	"github.com/dece2183/hexowl/utils"
)

func main() {
	if len(os.Args) > 1 {
		var expr string

		for i := 1; i < len(os.Args); i++ {
			if os.Args[i][0] == '-' {
				switch os.Args[i] {
				case "-ignore", "--ignore":
					goto ignoreArgs
				}
			} else {
				expr += os.Args[i]
			}
		}

		if len(expr) > 0 {
			words := utils.ParsePrompt(expr)
			err := calculate(words)
			if err != nil {
				fmt.Printf("\n\tError occurred: %s\n\n", err)
				os.Exit(1)
			}

			fmt.Println()
			os.Exit(0)
		}

	ignoreArgs:
	}

	var words []utils.Word
	terminal.EnableRawMode()
	stdreader := bufio.NewReader(os.Stdin)

	for {
		words = prompt(stdreader)
		if len(words) > 0 {
			var outstr string

			calcBeginTime := time.Now()
			err := calculate(words)
			calcTime := time.Since(calcBeginTime)

			if err != nil {
				outstr = fmt.Sprintf("\n\tError occurred: %s\n\n", err)
				outstr = syntax.Colorize(outstr, syntax.C_ERROR)
			} else {
				outstr = fmt.Sprintf("\n\tTime:\t%d ms\r\n\n", calcTime.Milliseconds())
				outstr = syntax.Colorize(outstr, utils.W_NONE)
			}

			fmt.Print(outstr)
		}
	}
}

func prompt(reader *bufio.Reader) []utils.Word {
	var inputString string
	inputString, _ = input.Prompt(os.Stdout, reader)
	return utils.ParsePrompt(inputString)
}

func calculate(words []utils.Word) error {
	var tbeg time.Time
	var tend time.Duration

	tbeg = time.Now()
	seq, err := operators.Prepare(words)
	if err != nil {
		return err
	}
	tend = time.Since(tbeg)
	fmt.Printf("\n\tseq: %+v; t: %s\r\n", seq, tend)

	tbeg = time.Now()
	res, err := operators.Execute(seq)
	if err != nil {
		return err
	}
	tend = time.Since(tbeg)
	fmt.Printf("\n\tres: %+v; t: %s\r\n", res, tend)

	tbeg = time.Now()
	operator, err := operators.Generate(words, make(map[string]interface{}))
	if err != nil {
		return err
	}
	tend = time.Since(tbeg)
	fmt.Printf("\n\tops: %+v; t: %s\r\n", operator, tend)

	tbeg = time.Now()
	val, err := operators.Calculate(operator, make(map[string]interface{}))
	if err != nil {
		return err
	}
	tend = time.Since(tbeg)
	fmt.Printf("\n\tclc: %+v; t: %s\r\n", val, tend)

	// if val != nil {
	// 	var resultStr string

	// 	switch v := val.(type) {
	// 	case string:
	// 		fmt.Printf("\n\t%s\r\n", v)
	// 		return nil
	// 	case float32, float64:
	// 		resultStr = fmt.Sprintf(
	// 			"\t%f\r\n\t\t0x%X\r\n\t\t0b%b\r\n",
	// 			v,
	// 			utils.ToNumber[uint64](val),
	// 			utils.ToNumber[uint64](val),
	// 		)
	// 	case int64, uint64:
	// 		resultStr = fmt.Sprintf(
	// 			"\t%d\r\n\t\t0x%X\r\n\t\t0b%b\r\n",
	// 			v,
	// 			utils.ToNumber[uint64](val),
	// 			utils.ToNumber[uint64](val),
	// 		)
	// 	case []interface{}:
	// 		resultStr = fmt.Sprintf("\t%v\r\n", v)
	// 		if len(v) > 0 {
	// 			var hstr, bstr string
	// 			switch v[0].(type) {
	// 			case float32, float64, int64, uint64:
	// 				for _, el := range v {
	// 					hstr += fmt.Sprintf("0x%X ", utils.ToNumber[uint64](el))
	// 					bstr += fmt.Sprintf("0b%b ", utils.ToNumber[uint64](el))
	// 				}
	// 				resultStr += fmt.Sprintf("\t\t[%s]\r\n", hstr[:len(hstr)-1])
	// 				resultStr += fmt.Sprintf("\t\t[%s]\r\n", bstr[:len(bstr)-1])
	// 			}
	// 		}
	// 	default:
	// 		resultStr = fmt.Sprintf("\t%v\r\n", v)
	// 	}

	// 	fmt.Print("\n\tResult:" + syntax.Highlight(resultStr))
	// }

	return nil
}
