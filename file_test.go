package codegen

import (
	"reflect"
	"testing"
)

func TestHello(t *testing.T) {
	file := NewFile("main", "examples/hello/hello_test.go")

	file.WriteBlock(
		Func(Var(file.TypeOf(reflect.TypeOf(t)))).Named("Test").Do(
			Call(file.Use("fmt", "Println"), file.Val("Hello, 世界")),
		),
	)

	file.WriteFile()
}

func TestRangeAndClose(t *testing.T) {
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

	file.WriteFile()
}
