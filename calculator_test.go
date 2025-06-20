package hexowl

import (
	"fmt"
	"reflect"
	"testing"
)

func initTestCalc() *Calculator {
	return NewCalculator(DefaultSystem())
}

func testExpr[T comparable](expr string, result T) error {
	calc := initTestCalc()
	res, err := calc.Eval(expr)
	if err != nil {
		return fmt.Errorf("Error occured: %w", err)
	}
	resTyped, ok := res.(T)
	if !ok {
		return fmt.Errorf("Result type mismatch (expected: %s, received: %s)", reflect.TypeOf(result), reflect.TypeOf(res))
	}
	if resTyped != result {
		return fmt.Errorf("Result value mismatch (expected: %v, received: %v)", result, res)
	}
	return nil
}

func TestBasic(t *testing.T) {
	const result = 7.0
	const expr = `-1+2+3*2`

	err := testExpr(expr, result)
	if err != nil {
		t.Fatal(err)
	}
}

func TestBrackets(t *testing.T) {
	const result = 4.0
	const expr = `(1+2) - (3-4)`

	err := testExpr(expr, result)
	if err != nil {
		t.Fatal(err)
	}
}

func TestOrder(t *testing.T) {
	const result = 20.0
	const expr = `1*(2+3)*4`

	err := testExpr(expr, result)
	if err != nil {
		t.Fatal(err)
	}
}

func TestVariables(t *testing.T) {
	const result = 28.0
	const expr = `a := 2; b = 14; a*b`

	err := testExpr(expr, result)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFunctions(t *testing.T) {
	const result = 120.0
	const expr = `
fac(x>=1)->x*fac(x-1)
fac(x)->1
fac(5)`

	err := testExpr(expr, result)
	if err != nil {
		t.Fatal(err)
	}
}
