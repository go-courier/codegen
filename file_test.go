package codegen

import (
	"fmt"
)

func ExampleNewFile_hello() {
	file := NewFile("main", "examples/hello/hello_test.go")

	file.WriteBlock(
		Func().Named("main").Do(
			Call(file.Use("fmt", "Println"), file.Val("Hello, 世界")),
		),
	)

	fmt.Println(string(file.Bytes()))
	// Output:
	//package main
	//
	//import (
	//	fmt "fmt"
	//)
	//
	//func main() {
	//	fmt.Println("Hello, 世界")
	//}
}

func ExampleNewFile_main() {
	file := NewFile("main", "examples/range-and-close/range-and-close.go")

	file.WriteBlock(
		Func(Var(Int, "n"), Var(Chan(Int), "c")).Named("fibonacci").Do(
			Expr("x, y := 0, 1"),
			For(Expr("i := 0"), Expr("i < n"), Expr("i++")).Do(
				file.Expr("c <- x"),
				file.Expr("x, y = y, x+y"),
			),
			Call("close", Id("c")),
		),
		Func().Named("main").Do(
			Define(Id("c")).By(Call("make", Chan(Int), Val(10))),
			Call("fibonacci", Call("cap", Id("c")), Id("c")).AsGo(),
			ForRange(Id("c"), "i").Do(
				Call(file.Use("fmt", "Println"), Id("i")),
			),
		),
	)

	fmt.Println(string(file.Bytes()))
	// Output:
	//package main
	//
	//import (
	//	fmt "fmt"
	//)
	//
	//func fibonacci(n int, c chan int) {
	//	x, y := 0, 1
	//	for i := 0; i < n; i++ {
	//		c <- x
	//		x, y = y, x+y
	//	}
	//	close(c)
	//}
	//
	//func main() {
	//	c := make(chan int, 10)
	//	go fibonacci(cap(c), c)
	//	for i := range c {
	//		fmt.Println(i)
	//	}
	//}
}
