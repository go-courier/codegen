package codegen

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDeclConst(t *testing.T) {
	tt := require.New(t)

	tt.Equal(`const a = "1"`, Stringify(
		DeclConst(
			Assign(Id("a")).By(Val("1")),
		),
	))

	tt.Equal(`const (
a int = iota
b
c
)`, Stringify(
		DeclConst(
			Assign(Var(Int, "a")).By(Iota),
			Assign(Id("b")),
			Assign(Id("c")),
		),
	))
}

func TestDeclVar(t *testing.T) {
	tt := require.New(t)

	tt.Equal(`var a = Fn("1")`, Stringify(
		DeclVar(
			Assign(Id("a")).By(Call("Fn", Val("1"))),
		),
	))

	tt.Equal(`var a, b string = Fn("1")`, Stringify(
		DeclVar(
			Assign(Var(String, "a", "b")).By(Call("Fn", Val("1"))),
		),
	))
}

func TestDeclType(t *testing.T) {
	tt := require.New(t)

	tt.Equal("type N string", Stringify(
		DeclType(
			Var(String, "N"),
		),
	))

	tt.Equal(`type (
// test
M = time.Time
N = time.Time
)`, Stringify(
		DeclType(
			Var(Type("time.Time"), "M").AsAlias().WithComments("test"),
			Var(Type("time.Time"), "N").AsAlias(),
		),
	))
}
