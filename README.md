# Jamlang

## What is it?
Jamlang is a simple, easy-to-use, interpreted programming language written in Go that is easy to extend via Go code, to add native functions/variables

It takes some parts of popular languages, such as:
* Rust
* Go
* JavaScript

## What does it look like?
```js
import "std/linkedlist.jam" /* use jamlang -i linkedlist or jamlang install linkedlist to get the linkedlist file */
fn Person(name, age): object {
    const this: object = {}
    this.getName = fn(): str {
        return name
    }

    this.getAge = fn(): i32 {
        return age
    }

    return this
}

const p: object = Person("Jamlang", 0.3)
println(p)
println(p.getName())
println((p.getName()).length)

const linkedlist: object = LinkedList()
linkedlist.add(1)
linkedlist.add(2)
linkedlist.add(3)
linkedlist.remove(2)
linkedlist.print()

const arr: list = [1,2,3,4]
foreach i in arr {
    println(i)
}
```

## Standard library
It has a very small standard library, which contains:
* LinkedList
* Math
* Random
* some algorithms

## Install
```sh
$ go install github.com/Jamlie/Jamlang@latest # or github.com/Jamlie/Jamlang@v1.5.0
```

## How to add native functions to it?
Adding functions via Go is rather simple, here's how to do it:

### First
Add the repository
```bash
$ go mod init your_package
$ go get github.com/Jamlie/Jamlang
```

### Second
```go
package main

import (
    "fmt"
    "os"

    "github.com/Jamlie/Jamlang/jamlang"
    "github.com/Jamlie/Jamlang/ast"
    . "github.com/Jamlie/Jamlang/runtimelang"
)

func main() {
    newEnv := CreateGlobalEnvironment()

    newEnv.DeclareVariable(/* name */ "foo", /* value */ MakeInt32Value(69), /* is const */ true, /* type */ ast.Int32Type)
    
    sumFn := MakeNativeFunction(func(args []RuntimeValue, env Environment) RuntimeValue {
        if len(args) != 2 {
            fmt.Fprintln(os.Stderr, "Error: sum takes 2 arguments")
            os.Exit(0)
        }

        if _, ok := args[0].(Int32Type); !ok {
            fmt.Fprintln(os.Stderr, "Error: arguments must be of type number")
            os.Exit(0)
        }
        num1 := args[0].(Int32Type).Value

        if _, ok := args[1].(Int32Type); !ok {
            fmt.Fprintln(os.Stderr, "Error: arguments must be of type number")
            os.Exit(0)
        }
        num2 := args[1].(Int32Type).Value
        
        return MakeInt32Value(num1 + num2)
    }, "sum")

    env.DeclareVariable("sum", sumFn, true, ast.Int32Type)
    
    jamlang.CallMain(newEnv)
}
```

### Third
go build

```bash
$ ./your_package -r fileName.jam
```

## Created by
**Omar Estietie (Jam)**
