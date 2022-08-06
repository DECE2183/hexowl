# **hexowl**

**hexowl** is a Lightweight and flexible programmer's calculator written in Go.

The main purpose of hexowl is to perform operations on numbers regardless of their base. A single expression can contain decimal, hexadecimal, and binary numbers.

![Work demonstration](/.github/assets/demo.gif)

## Features
 - Support for operations on decimal, hexadecimal and binary numbers;
 - Bitwise operators;
 - Boolean operators;
 - User defined variables;
 - Ability to save and load the working environment.

# Building

There are no dependencies, so you can simply type a build command in the cloned repository folder.

```bash
go build
```

# Reference

### Operators

|Operator           |Syntax |
|-------------------|-------|
|Positive bits count|`#`    |
|Bitwise NOT        |`~`    |
|Bitclear (AND NOT) |`&^`   |
|Bitwise XOR        |`^`    |
|Bitwise AND        |`&`    |
|Bitwise OR         |`\|`   |
|Right shift        |`>>`   |
|Left shift         |`<<`   |
|Modulo             |`%`    |
|Division           |`/`    |
|Exponentiation     |`**`   |
|Multiplication     |`*`    |
|Addition           |`+`    |
|Subtraction        |`-`    |
|Logical NOT        |`!`    |
|Logical AND        |`&&`   |
|Logical OR         |`\|\|` |
|Logical not equal  |`!=`   |
|Logical equal      |`==`   |
|Argument iterator  |`,`    |
|Divide and assign  |`/=`   |
|Mutiply and assign |`*=`   |
|Add and assign     |`+=`   |
|Subtract and assign|`-=`   |
|Assign             |`=`    |

### Built in constants

|Constant           |Value              |
|-------------------|-------------------|
|`pi`               |`3.141592653589793`|
|`true`             |`1`                |
|`false`            |`0`                |

### Built in functions

|Function           |Meaning                            |
|-------------------|-----------------------------------|
|`sin(x)`           |Sine                               |
|`cos(x)`           |Cosine                             |
|`pow(a,b)`         |Exponentiation                     |
|`vars()`           |List defined user variables        |
|`save(id)`         |Save working environment with `id` |
|`load(id)`         |Load working environment with `id` |
|`clear()`          |Clear terminal                     |
|`exit(error_code)` |Exit with error code               |
