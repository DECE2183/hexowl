package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/dece2183/hexowl/builtin"
	"github.com/dece2183/hexowl/input"
	"github.com/dece2183/hexowl/operators"
	"github.com/dece2183/hexowl/utils"

	"net/http"
	netpprof "net/http/pprof"
)

var profileEnabled bool

func main() {
	builtin.FuncsInit(os.Stdout)

	if len(os.Args) > 1 {
		var expr string

		for i := 1; i < len(os.Args); i++ {
			if os.Args[i][0] == '-' {
				switch os.Args[i] {
				case "-prof", "--prof":
					mux := http.NewServeMux()
					mux.HandleFunc("/hxl-pprof/", netpprof.Index)
					mux.HandleFunc("/hxl-pprof/cmdline", netpprof.Cmdline)
					mux.HandleFunc("/hxl-pprof/profile", netpprof.Profile)
					mux.HandleFunc("/hxl-pprof/symbol", netpprof.Symbol)
					mux.HandleFunc("/hxl-pprof/trace", netpprof.Trace)
					go http.ListenAndServe(":8080", mux)
					fmt.Print("\nprofiler enabled at http://localhost:8080/hxl-pprof/profile\n\n")
					profileEnabled = true
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
	}

	var words []utils.Word
	input.EnableRawMode()
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

func prompt(reader *bufio.Reader) []utils.Word {
	var inputString string
	inputString, _ = input.Prompt(os.Stdout, reader)
	return utils.ParsePrompt(inputString)
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
		case float32, float64:
			fmt.Printf("\n\tResult:\t%f\r\n", val)
			fmt.Printf("\t\t0x%X\r\n", utils.ToNumber[uint64](val))
			fmt.Printf("\t\t0b%b\r\n", utils.ToNumber[uint64](val))
		case int64, uint64:
			fmt.Printf("\n\tResult:\t%d\r\n", val)
			fmt.Printf("\t\t0x%X\r\n", utils.ToNumber[uint64](val))
			fmt.Printf("\t\t0b%b\r\n", utils.ToNumber[uint64](val))
		case []interface{}:
			fmt.Printf("\n\tResult:\t%v\r\n", val)
			if len(v) > 0 {
				var hstr, bstr string
				switch v[0].(type) {
				case float32, float64, int64, uint64:
					for _, el := range v {
						hstr += fmt.Sprintf("0x%X ", utils.ToNumber[uint64](el))
						bstr += fmt.Sprintf("0b%b ", utils.ToNumber[uint64](el))
					}
					fmt.Printf("\t\t[%s]\r\n", hstr[:len(hstr)-1])
					fmt.Printf("\t\t[%s]\r\n", bstr[:len(bstr)-1])
				}
			}
		default:
			fmt.Printf("\n\tResult:\t%v\r\n", val)
		}
	}

	if profileEnabled {
		var mem runtime.MemStats
		runtime.ReadMemStats(&mem)
		fmt.Printf("\n\tStack usage: %0.2f kB\r\n", float32(mem.StackSys)/1024)
		fmt.Printf("\tHeap usage:  %0.2f kB\r\n", float32(mem.HeapInuse)/1024)
	}

	return nil
}
