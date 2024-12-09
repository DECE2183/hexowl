package main

import (
	"fmt"
	"os"
	"time"

	"github.com/dece2183/hexowl/v2/calculator"
	"github.com/dece2183/hexowl/v2/calculator/builtin"
	"github.com/dece2183/hexowl/v2/input"
	"github.com/dece2183/hexowl/v2/input/syntax"
	"github.com/dece2183/hexowl/v2/input/terminal"
	"github.com/dece2183/hexowl/v2/utils"
)

func main() {
	calc := calculator.NewCalculator(builtin.DefaultSystem())

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
			res, err := calc.Eval(expr)
			if err != nil {
				displayError(err)
				os.Exit(1)
			}

			displayResult(res)
			fmt.Println()
			os.Exit(0)
		}

	ignoreArgs:
	}

	var (
		calcBeginTime time.Time
		calcDuration  time.Duration
		inputStr      string
		result        interface{}
		err           error
	)

	terminal.EnableRawMode()
	console := input.NewConsole(os.Stdout, os.Stdin, []input.Predictable{
		calc.GetUserContainer(),
		calc.GetBuiltinContainer(),
	})

	for {
		inputStr, err = console.Prompt()
		if err != nil {
			goto errorOccured
		}

		calcBeginTime = time.Now()
		result, err = calc.Eval(inputStr)
		calcDuration = time.Since(calcBeginTime)
		if err != nil {
			goto errorOccured
		}

		displayResult(result)
		fmt.Printf("\n\tTime:\t%d ms\r\n\n", calcDuration.Milliseconds())
		continue

	errorOccured:
		displayError(err)
	}
}

func displayResult(result interface{}) {
	if result == nil {
		return
	}

	var resultStr string

	switch v := result.(type) {
	case string:
		fmt.Printf("\n\t%s\r\n", v)
		return
	case float32, float64:
		resultStr = fmt.Sprintf(
			"\t%f\r\n\t\t0x%X\r\n\t\t0b%b\r\n",
			v,
			utils.ToNumber[uint64](result),
			utils.ToNumber[uint64](result),
		)
	case int64, uint64:
		resultStr = fmt.Sprintf(
			"\t%d\r\n\t\t0x%X\r\n\t\t0b%b\r\n",
			v,
			utils.ToNumber[uint64](result),
			utils.ToNumber[uint64](result),
		)
	case []interface{}:
		resultStr = fmt.Sprintf("\t%v\r\n", v)
		if len(v) > 0 {
			var hstr, bstr string
			switch v[0].(type) {
			case float32, float64, int64, uint64:
				for _, el := range v {
					hstr += fmt.Sprintf("0x%X ", utils.ToNumber[uint64](el))
					bstr += fmt.Sprintf("0b%b ", utils.ToNumber[uint64](el))
				}
				resultStr += fmt.Sprintf("\t\t[%s]\r\n", hstr[:len(hstr)-1])
				resultStr += fmt.Sprintf("\t\t[%s]\r\n", bstr[:len(bstr)-1])
			}
		}
	default:
		resultStr = fmt.Sprintf("\t%v\r\n", v)
	}

	fmt.Print("\n\tResult:" + syntax.Highlight(resultStr))
}

func displayError(err error) {
	fmt.Print(syntax.Colorize("\n\tError occurred: ", syntax.C_ERROR), err, "\n\n")
}
