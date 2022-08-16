# **hexowl**

[![Go Report Card](https://goreportcard.com/badge/github.com/dece2183/hexowl)](https://goreportcard.com/report/github.com/dece2183/hexowl)
[![Release](https://img.shields.io/github/v/release/dece2183/hexowl)](https://github.com/dece2183/hexowl/releases)

**hexowl** is a Lightweight and flexible programmer's calculator written in Go.

The main purpose of hexowl is to perform operations on numbers regardless of their base. A single expression can contain decimal, hexadecimal, and binary numbers.

![Work demonstration](/.github/assets/demo.gif)

## Features
 - Support for operations on decimal, hexadecimal and binary numbers;
 - Bitwise operators;
 - Boolean operators;
 - User defined variables;
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
|Logical not equal  |`!=`       |
|Logical equal      |`==`       |
|Argument iterator  |`,`        |
|Divide and assign  |`/=`       |
|Mutiply and assign |`*=`       |
|Add and assign     |`+=`       |
|Subtract and assign|`-=`       |
|Assign             |`=`        |

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
|`ceil(x)`          |The least integer value >= x               |
|`floor(x)`         |The greatest integer value <= x            |
|`popcnt(x)`        |Positive bits count                        |
|`vars()`           |List built in and user defined variables   |
|`funcs()`          |List alailable functions                   |
|`save(id)`         |Save working environment with `id`         |
|`load(id)`         |Load working environment with `id`         |
|`clear()`          |Clear terminal                             |
|`exit(error_code)` |Exit with error code                       |
