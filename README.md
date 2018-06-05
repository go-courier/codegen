## Codegen

[![GoDoc Widget](https://godoc.org/github.com/go-courier/codegen?status.svg)](https://godoc.org/github.com/go-courier/codegen)
[![Build Status](https://travis-ci.org/go-courier/codegen.svg?branch=master)](https://travis-ci.org/go-courier/codegen)
[![codecov](https://codecov.io/gh/go-courier/codegen/branch/master/graph/badge.svg)](https://codecov.io/gh/go-courier/codegen)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-courier/codegen)](https://goreportcard.com/report/github.com/go-courier/codegen)

Helpers fo generating go codes

## Example 

```go
package main

import (
	"github.com/go-courier/codegen"
)

func main()  {
	file := codegen.NewFile("main", "examples/range-and-close/range-and-close.go")
 
    file.WriteBlock(
        codegen.Func(codegen.Var(codegen.Int, "n"), codegen.Var(codegen.Chan(codegen.Int), "c"), ).Named("fibonacci").Do(
            file.Expr("x, y := 0, 1"),
            codegen.For(file.Expr("i := 0"), file.Expr("i < n"), file.Expr("i++")).Do(
                file.Expr("c <- x"),
                file.Expr("x, y = y, x+y"),
            ),
            codegen.Call("close", codegen.Id("c")),
        ),
        codegen.Func().Named("main").Do(
            codegen.Define(codegen.Id("c")).By(codegen.Call("make", codegen.Chan(codegen.Int), file.Val(10))),
            codegen.Call("fibonacci", codegen.Call("cap", codegen.Id("c")), codegen.Id("c")).AsGo(),
            codegen.ForRange(codegen.Id("c"), "i").Do(
                codegen.Call(file.Use("fmt", "Println"), codegen.Id("i")),
            ),
        ),
    )
    
    file.WriteFile()
}
```

will generate file `examples/range-and-close/range-and-close.go`

```go
package main

import (
	fmt "fmt"
)

func fibonacci(n int, c chan int) {
	x, y := 0, 1
	for i := 0; i < n; i++ {
		c <- x
		x, y = y, x+y
	}
	close(c)
}
func main() {
	c := make(chan int, 10)
	go fibonacci(cap(c), c)
	for i := range c {
		fmt.Println(i)
	}
}
```