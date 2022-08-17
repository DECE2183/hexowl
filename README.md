# **hexowl**

[![Go Report Card](https://goreportcard.com/badge/github.com/dece2183/hexowl)](https://goreportcard.com/report/github.com/dece2183/hexowl)
[![Release](https://img.shields.io/github/v/release/dece2183/hexowl)](https://github.com/dece2183/hexowl/releases)

**hexowl** is a Lightweight and flexible programmer's calculator with user variables and functions support written in Go.

The main purpose of hexowl is to perform operations on numbers regardless of their base. A single expression can contain decimal, hexadecimal, and binary numbers.

![Work demonstration](/.github/assets/demo.gif)

## Features
 - Support for operations on decimal, hexadecimal and binary numbers;
 - Bitwise operators;
 - Boolean operators;
 - User defined variables;
 - User defined functions;
 - Ability to save and load the working environment.

# Installation

```bash
go install github.com/dece2183/hexowl@latest 
```

# Building

There are no dependencies, so you can simply type a build command in the cloned repository folder.

```bash
go build
```

# Reference

### Operators

|Operator           |Syntax     |
|-------------------|-----------|
|Positive bits count|`#`        |
|Bitwise NOT        |`~`        |
|Bitclear (AND NOT) |`&~` `&^`  |
|Bitwise XOR        |`^`        |
|Bitwise AND        |`&`        |
|Bitwise OR         |`\|`       |
|Right shift        |`>>`       |
|Left shift         |`<<`       |
|Modulo             |`%`        |
|Division           |`/`        |
|Exponentiation     |`**`       |
|Multiplication     |`*`        |
|Subtraction        |`-`        |
|Addition           |`+`        |
|Logical NOT        |`!`        |
|Logical AND        |`&&`       |
|Logical OR         |`\|\|`     |
|Less or equal      |`<=`       |
|More or equal      |`>=`       |
|Less               |`<`        |
|More               |`>`        |
|Not equal          |`!=`       |
|Equal              |`==`       |
|Enumerate          |`,`        |
|Divide and assign  |`/=`       |
|Mutiply and assign |`*=`       |
|Add and assign     |`+=`       |
|Subtract and assign|`-=`       |
|Assign             |`=`        |
|Declare function   |`->`       |

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

### Built in functions

|Function           |Meaning                                    |
|-------------------|-------------------------------------------|
|`sin(x)`           |Sine                                       |
|`cos(x)`           |Cosine                                     |
|`pow(x,y)`         |Exponentiation                             |
|`sqrt(x)`          |Square root                                |
|`exp(x)`           |The base-e exponential                     |
|`ceil(x)`          |The least integer value >= `x`             |
|`floor(x)`         |The greatest integer value <= `x`          |
|`popcnt(x)`        |Positive bits count                        |
|`vars()`           |List built in and user defined variables   |
|`clvars()`         |Delet user defined variables               |
|`funcs()`          |List alailable functions                   |
|`clfuncs()`        |Delet user defined functions               |
|`save(id)`         |Save working environment with `id`         |
|`load(id)`         |Load working environment with `id`         |
|`clear()`          |Clear terminal                             |
|`exit(error_code)` |Exit with error code                       |

### User functions

To declare a function, you must type its name, explain the arguments in `(` `)` and write the body of the function after `->` operator.

It should look like this:
```
>: mul(a,b) -> a * b
```

Once declared, you can call this function as a builtin:
```
>: mul(2,4)

    Result: 8
            0x8
            0b1000

    Time:   0 ms
```

### User function variations

You can also create variants of functions with expressions right in the explanation of the arguments.

Let's look at a simple example of declaring a factorial function:
```
>: f(x == 0) -> 1
>: f(x > 0) -> x * f(x-1)
```

When calling such a function, the interpreter tries to find a suitable variant depending on the arguments passed, and then calls it.

### Arrays and variadic arguments

You can define arrays with the enumerator operator `,`:
```
>: x = 1,2,3,4
```

All functions receive arguments as an array, so the expressions `foo(x)` and `foo(1,2,3,4)` are similar.

There is a single `@` keyword to handle such things. If it is specified as the last argument in a function declaration, it will receive an array of the arguments passed to it. The behavior is similar to the `...` and `__VA_ARGS__` preprocessor macros in C language.

An example of a function that calculates the sum of all elements of an array:
```
>: arrsum(a) -> a
>: arrsum(a, @) -> a+arrsum(@)
```

An example of a function that increments all elements of an array:
```
>: arrinc(v, a) -> a+v
>: arrinc(v, a, @) -> (a+v) , arrinc(v,@)
```
