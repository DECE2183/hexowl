package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/dece2183/hexowl/builtin"
	"github.com/dece2183/hexowl/operators"
	"github.com/dece2183/hexowl/utils"
)

func main() {
	builtin.FuncsInit()

	if len(os.Args) > 1 {
		var input string

		for i := 1; i < len(os.Args); i++ {
			input += os.Args[i]
		}

		words := utils.ParsePrompt(input)
		err := calculate(words)
		if err != nil {
			fmt.Printf("\n\tError occurred: %s\n\n", err)
			os.Exit(1)
		}

		fmt.Println()
		os.Exit(0)
	} else {
		var words []utils.Word
		stdreader := bufio.NewReader(os.Stdin)

		for {
			words = prompt(stdreader)
			if len(words) > 0 {
				calcBeginTime := time.Now()
				err := calculate(words)
				calcTime := time.Since(calcBeginTime)

				if err != nil {
					fmt.Printf("\n\tError occurred: %s\n\n", err)
				} else {
					fmt.Printf("\n\tTime:\t%d ms\r\n\n", calcTime.Milliseconds())
				}
			}
		}
	}
}

func prompt(reader *bufio.Reader) []utils.Word {
	var input string

	fmt.Printf(">: ")
	input, _ = reader.ReadString('\n')

	return utils.ParsePrompt(input)
}

func calculate(words []utils.Word) error {
	operator, err := operators.Generate(words, make(map[string]interface{}))
	if err != nil {
		return err
	}

	val, err := operators.Calculate(operator, make(map[string]interface{}))
	if err != nil {
		return err
	}

	if val != nil {
		switch v := val.(type) {
		case string:
			fmt.Printf("\n\t%s\r\n", v)
		case bool:
			fmt.Printf("\n\tResult:\t%v\r\n", v)
		default:
			fmt.Printf("\n\tResult:\t%v\r\n", val)
			fmt.Printf("\t\t0x%X\r\n", utils.ToNumber[uint64](val))
			fmt.Printf("\t\t0b%b\r\n", utils.ToNumber[uint64](val))
		}
	}

	return nil
}
