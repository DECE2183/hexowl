# ![Logo](/.assets/hexowl-logo-full.svg)

[![TryIt Card](https://img.shields.io/badge/try%20it-online-green)](https://dece2183.github.io/web-hexowl)
[![GitHub License](https://img.shields.io/github/license/dece2183/hexowl)](https://github.com/dece2183/hexowl/blob/main/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/dece2183/hexowl)](https://goreportcard.com/report/github.com/dece2183/hexowl)
[![Release](https://img.shields.io/github/v/release/dece2183/hexowl)](https://github.com/dece2183/hexowl/releases)
[![Go Reference](https://pkg.go.dev/badge/github.com/dece2183/hexow.svg)](https://pkg.go.dev/github.com/dece2183/hexowl)

**hexowl** is a Lightweight and flexible programmer's calculator with user variables and functions support written in Go.

The main purpose of hexowl is to perform operations on numbers regardless of their base. A single expression can contain decimal, hexadecimal, and binary numbers.

![Work demonstration](/.assets/demo.gif)

## Features
 - Support for operations on decimal, hexadecimal and binary numbers;
 - Bitwise operators;
 - Boolean operators;
 - User defined variables;
 - User defined functions;
 - Ability to save and load the working environment.

## Installation

```bash
go install github.com/dece2183/hexowl@latest 
```

## Building

There are no dependencies, so you can simply type a build command in the cloned repository folder.

```bash
go build
```

## Reference

### Operators

|Operator                |Syntax     |
|------------------------|-----------|
|Positive bits count     |`#`        |
|Bitwise NOT             |`~`        |
|Bitclear (AND NOT)      |`&~` `&^`  |
|Bitwise XOR             |`^`        |
|Bitwise AND             |`&`        |
|Bitwise OR              |`\|`       |
|Right shift             |`>>`       |
|Left shift              |`<<`       |
|Modulo                  |`%`        |
|Division                |`/`        |
|Exponentiation          |`**`       |
|Multiplication          |`*`        |
|Subtraction             |`-`        |
|Addition                |`+`        |
|Logical NOT             |`!`        |
|Less or equal           |`<=`       |
|More or equal           |`>=`       |
|Less                    |`<`        |
|More                    |`>`        |
|Not equal               |`!=`       |
|Equal                   |`==`       |
|Logical AND             |`&&`       |
|Logical OR              |`\|\|`     |
|Enumerate               |`,`        |
|Bitwise OR and assign   |`\|=`      |
|Bitwise AND and assign  |`&=`       |
|Divide and assign       |`/=`       |
|Mutiply and assign      |`*=`       |
|Add and assign          |`+=`       |
|Subtract and assign     |`-=`       |
|Local assign            |`:=`       |
|Assign                  |`=`        |
|Sequence                |`;`        |
|Declare function        |`->`       |

### Built in constants

|Constant           |Value              |
|-------------------|-------------------|
|`pi`               |`3.141592653589793`|
|`e`                |`2.718281828459045`|
|`true`             |`1`                |
|`false`            |`0`                |
|`inf`              |`+Inf`             |
|`nan`              |`NaN`              |
|`nil`              |`nil`              |
|`help`             |Help Message       |
|`version`          |hexowl version     |

### Built in functions

| Function    | Arguments        | Description
|-------------|------------------|--------------------------------------------------------------------------|
| `acos`      | (`x`)            | The arccosine of the radian argument `x`                                 |
| `asin`      | (`x`)            | The arcsine of the radian argument `x`                                   |
| `atan`      | (`x`)            | The arctangent of the radian argument `x`                                |
| `ceil`      | (`x`)            | The least integer value greater than or equal to `x`                     |
| `clear`     | ( )              | Clear screen                                                             |
| `clfuncs`   | ( )              | Delete user defined functions                                            |
| `clvars`    | ( )              | Delete user defined variables                                            |
| `cos`       | (`x`)            | The cosine of the radian argument `x`                                    |
| `envs`      | ( )              | List all available environments                                          |
| `exit`      | (`code`)         | Exit with error `code`                                                   |
| `exp`       | (`x`)            | The base-e exponential of `x`                                            |
| `floor`     | (`x`)            | The greatest integer value less than or equal to `x`                     |
| `funcs`     | ( )              | List alailable functions                                                 |
| `import`    | (`id`,`unit`)    | Import unit from the working environment with `id`                       |
| `load`      | (`id`)           | Load working environment with `id`                                       |
| `log10`     | (`x`)            | The decimal logarithm of `x`                                             |
| `log2`      | (`x`)            | The binary logarithm of `x`                                              |
| `logn`      | (`x`)            | The natural logarithm of `x`                                             |
| `popcnt`    | (`x`)            | The number of one bits ("population count") in `x`                       |
| `pow`       | (`x`,`y`)        | The base-`x` exponential of `y`                                          |
| `rand`      | (`a`,`b`)        | The random number in the range [a,b) or [0,1) if no arguments are passed |
| `rmfunc`    | (`name`)         | Delete user function with `name`                                         |
| `rmfuncvar` | (`name`,`varid`) | Delete user function `name` variation number `varid`                     |
| `rmvar`     | (`name`)         | Delete user variable with `name`                                         |
| `round`     | (`x`)            | The nearest integer, rounding half away from zero                        |
| `save`      | (`envname`)      | Save working environment with `envname`                                  |
| `sin`       | (`x`)            | The sine of the radian argument `x`                                      |
| `sqrt`      | (`x`)            | The square root of `x`                                                   |
| `tan`       | (`x`)            | The tangent of the radian argument `x`                                   |
| `vars`      | ( )              | List available variables                                                 |

### User functions

To declare a function, you must type its name, explain the arguments in `(` `)` and write the body of the function after `->` operator.

It should look like this:
```hexowl
>: mul(a,b) -> a * b
```

Once declared, you can call this function as a builtin:
```hexowl
>: mul(2,4)

    Result: 8
            0x8
            0b1000

    Time:   0 ms
```

### User function variations

You can also create variants of functions with expressions right in the explanation of the arguments.

Let's look at a simple example of declaring a factorial function:
```hexowl
>: f(x == 0) -> 1
>: f(x > 0) -> x * f(x-1)
```

When calling such a function, the interpreter tries to find a suitable variant depending on the arguments passed, and then calls it.

### Arrays and variadic arguments

You can define arrays with the enumerator operator `,`:
```hexowl
>: x = 1,2,3,4
```

All functions receive arguments as an array, so the expressions `foo(x)` and `foo(1,2,3,4)` are similar.

There is a single `@` keyword to handle such things. If it is specified as the last argument in a function declaration, it will receive an array of the arguments passed to it. The behavior is similar to the `...` and `__VA_ARGS__` preprocessor macros in C language.

An example of a function that calculates the sum of all elements of an array:
```hexowl
>: arrsum(a) -> a
>: arrsum(a, @) -> a+arrsum(@)
```

An example of a function that increments all elements of an array:
```hexowl
>: arrinc(v, a) -> a+v
>: arrinc(v, a, @) -> (a+v) , arrinc(v,@)
```

## Integration guide

Hexowl is specially designed for use as an embeddable calculator.

An example of a minimal setup is shown below:

```go
package main

import (
	"fmt"

	"github.com/dece2183/hexowl/operators"
	"github.com/dece2183/hexowl/utils"
)

const expresion = "2+2"

func main() {
    localVars := make(map[string]interface{})
    words := utils.ParsePrompt(expresion)

    operatorTree, err := operators.Generate(words, localVars))
	if err != nil {
		return err
	}

	result, err := operators.Calculate(operatorTree, localVars))
	if err != nil {
		return err
	}

    fmt.Printf("%s = %v", expresion, result);
}
```

For more specific designs, it is posible to provide an sdtout writer and callbacks for working with environment save files.

```go
package main

import (
	"bytes"
	"fmt"
	"io"

	"github.com/dece2183/hexowl/builtin"
)

type dummyCloser bytes.Buffer

var outbuff = &bytes.Buffer{}
var envFiles map[string]*dummyCloser

func (dc *dummyCloser) Close() error {
	return nil
}

func (dc *dummyCloser) Read(dest []byte) (int, error) {
	b := bytes.Buffer(*dc)
	return b.Read(dest)
}

func (dc *dummyCloser) Write(data []byte) (int, error) {
	b := bytes.Buffer(*dc)
	return b.Write(data)
}

func init() {
	sysDesc := builtin.System{
		Stdout: outbuff,
		ListEnvironments: func() ([]string, error) {
			return maps.Keys(envFiles), nil
		},
		WriteEnvironment: func(name string) (io.WriteCloser, error) {
			if _, ok := envFiles[name]; !ok {
				envFiles[name] = &dummyCloser{}
			}
			return envFiles[name], nil
		},
		ReadEnvironment: func(name string) (io.ReadCloser, error) {
			if _, ok := envFiles[name]; !ok {
				return nil, fmt.Errorf("not found")
			}
			return envFiles[name], nil
		},
	}

	// Now all the additional output will be printed in outbuff.
	// And environment files will be seved to and loaded from envFiles map.
	builtin.SystemInit(sysDesc)
}
```

There are also functions for registering and manage self-written built-in functions and constants. They are described in [`hexowl/builtin`](https://pkg.go.dev/github.com/dece2183/hexowl/builtin) package.
