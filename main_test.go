package main

// import (
// 	"reflect"
// 	"testing"

// 	"github.com/dece2183/hexowl/operators"
// 	"github.com/dece2183/hexowl/utils"
// )

// const testExpr = "a = pow(2,6) + (1)"
// const testExprRes = 65

// var testParseRes = []utils.Word{
// 	{Type: utils.W_UNIT, Literal: "a"},
// 	{Type: utils.W_OP, Literal: "="},
// 	{Type: utils.W_UNIT, Literal: "pow"},
// 	{Type: utils.W_CTL, Literal: "("},
// 	{Type: utils.W_NUM_DEC, Literal: "2"},
// 	{Type: utils.W_OP, Literal: ","},
// 	{Type: utils.W_NUM_DEC, Literal: "6"},
// 	{Type: utils.W_CTL, Literal: ")"},
// 	{Type: utils.W_OP, Literal: "+"},
// 	{Type: utils.W_CTL, Literal: "("},
// 	{Type: utils.W_NUM_DEC, Literal: "1"},
// 	{Type: utils.W_CTL, Literal: ")"},
// }

// var testParsedExpr []utils.Word
// var testGeneratedOps *operators.Operator

// func TestParsePrompt(t *testing.T) {
// 	testParsedExpr = utils.ParsePrompt(testExpr)
// 	if testParsedExpr == nil {
// 		t.Error("failed to parse prompt")
// 		return
// 	}

// 	for i, w := range testParsedExpr {
// 		if w != testParseRes[i] {
// 			t.Errorf("failed to parse prompt at word #%d:\r\n\texpected: %+v\r\n\tresult:   %+v\r\n", i, testParseRes[i], w)
// 			return
// 		}
// 	}
// }

// func TestGenerateOp(t *testing.T) {
// 	var err error

// 	if testParsedExpr == nil {
// 		t.Log("parsing prompt...")
// 		TestParsePrompt(t)
// 	}

// 	testGeneratedOps, err = operators.Generate(testParsedExpr, make(map[string]interface{}))
// 	if err != nil {
// 		t.Errorf("failed to generate operators: %s", err)
// 	}
// }

// func TestCalculateExpr(t *testing.T) {
// 	if testGeneratedOps == nil {
// 		t.Log("generating operators...")
// 		TestGenerateOp(t)
// 	}

// 	res, err := operators.Calculate(testGeneratedOps, make(map[string]interface{}))
// 	if err != nil {
// 		t.Errorf("failed to calculate operators: %s", err)
// 	}

// 	resNum, ok := res.(float64)
// 	if !ok {
// 		resType := reflect.TypeOf(res)
// 		t.Errorf("failed to calculate operators, result types mismatch:\r\n\texpected: float64\r\n\tresult:   %+v\r\n", resType)
// 	}

// 	if resNum != testExprRes {
// 		t.Errorf("failed to calculate operators, wrong result:\r\n\texpected: %d\r\n\tresult:   %f\r\n", testExprRes, resNum)
// 	}
// }
