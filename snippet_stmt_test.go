package codegen

import (
	"go/token"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSnippet_SwitchStmt(t *testing.T) {
	tt := require.New(t)

	tt.Equal(`switch os := runtime.GOOS; os {
case "darwin", "darwin32":
case "linux":
default:
}`, Stringify(
		Switch(
			Id("os"),
		).InitWith(Expr("os := runtime.GOOS")).When(
			Clause(Val("darwin"), Val("darwin32")),
			Clause(Val("linux")),
			Clause(),
		),
	))

	tt.Equal(`switch i {
case 0:
case fn():
}`, Stringify(
		Switch(
			Id("i"),
		).When(
			Clause(Val(0)),
			Clause(Call("fn")),
		),
	))

	tt.Equal(`switch {
case t.Hour() < 12:
case t.Hour() < 17:
default:
i++
}`, Stringify(
		Switch(nil).When(
			Clause(Expr("t.Hour() < 12")),
			Clause(Expr("t.Hour() < 17")),
			Clause().Do(Expr("i++")),
		),
	))
}

func TestSnippet_SelectStmt(t *testing.T) {
	tt := require.New(t)

	tt.Equal(`select {
case c <- x:
x, y = y, x+y
case <- quit:
fmt.Println("quit")
return
}`, Stringify(
		Select(
			Clause(Expr("c <- x")).Do(Expr("x, y = y, x+y")),
			Clause(Expr("<- quit")).Do(
				Expr(`fmt.Println("quit")`),
				Return(),
			),
		),
	))
}

func TestSnippet_IfStmt(t *testing.T) {
	tt := require.New(t)

	tt.Equal(`if i := 1; i == 1 {
}`, Stringify(
		If(Expr("i == 1")).InitWith(Expr("i := 1")),
	))

	tt.Equal(`if i == 1 {
} else {
}`, Stringify(
		If(Expr("i == 1")).Else(
			If(nil),
		),
	))

	tt.Equal(`if i := 1; i == 1 {
i++
} else if i == 2 {
}`, Stringify(
		If(Expr("i == 1")).InitWith(Expr("i := 1")).Do(Expr("i++")).
			Else(If(Expr("i == 2"))),
	))
}

func TestSnippet_For(t *testing.T) {
	tt := require.New(t)

	tt.Equal(`for i := 1; i < 10; i++ {
}`, Stringify(
		For(Expr("i := 1"), Expr("i < 10"), Expr("i++")).Do(),
	))

	tt.Equal(`for a < 1 {
}`, Stringify(
		For(nil, Expr("a < 1"), nil),
	))

	tt.Equal(`for {
}`, Stringify(
		For(nil, nil, nil),
	))

	tt.Equal(`for _, v := range []string{
"1",
"2",
} {
}`, Stringify(
		ForRange(Val([]string{"1", "2"}), "_", "v"),
	))

	tt.Equal(`for range []string{
"1",
"2",
} {
i++
}`, Stringify(
		ForRange(Val([]string{"1", "2"})).Do(Expr("i++")),
	))
}

func TestSnippet_AssignStmt(t *testing.T) {
	tt := require.New(t)

	tt.Equal(`a = Fn("1")`, Stringify(
		Assign(Id("a")).By(Call("Fn", Val("1"))),
	))

	tt.Equal(`a += Fn("1")`, Stringify(
		AssignWith(token.ADD_ASSIGN, Id("a")).By(Call("Fn", Val("1"))),
	))

	tt.Equal(`a := Fn("1")`, Stringify(
		Define(Id("a")).By(Call("Fn", Val("1"))),
	))

	tt.Equal(`a, b := "1", 1`, Stringify(
		Define(Id("a"), Id("b")).By(Val("1"), Val(1)),
	))
}

func TestSnippet_ReturnStmt(t *testing.T) {
	tt := require.New(t)

	tt.Equal(`return`, Stringify(
		Return(),
	))

	tt.Equal(`return a`, Stringify(
		Return(Id("a")),
	))

	tt.Equal(`return a, fmt.Sprintf("%s", 1)`, Stringify(
		Return(Id("a"), Call("fmt.Sprintf", Val("%s"), Val(1))),
	))
}
