package codegen

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSnippet_Star(t *testing.T) {
	tt := require.New(t)

	tt.Equal("*bool", Stringify(Star(Bool)))
	tt.Equal("**bool", Stringify(Star(Star(Bool))))
}

func TestSnippet_Sel(t *testing.T) {
	tt := require.New(t)

	tt.Equal(`fmt.Printf("1", 1, '1', true)`, Stringify(
		Sel(
			Id("fmt"),
			Call("Printf", Val("1"), Val(1), Val('1'), Val(true)),
		),
	))

	tt.Equal(`fmt.Printf("1", []interface {}{
"1",
"2",
}...)`, Stringify(
		Sel(
			Id("fmt"),
			Call("Printf", Val("1"), Val([]interface{}{"1", "2"})).WithEllipsis(),
		),
	))

	tt.Equal(`r.Request("GET").Do(req, &(resp))`, Stringify(
		Sel(
			Id("r"),
			Call("Request", Val("GET")),
			Call("Do", Id("req"), Unary(Paren(Id("resp")))),
		),
	))
}

func TestSnippet_Call(t *testing.T) {
	tt := require.New(t)

	tt.Equal(`fn("1", 1, '1', true, &a)`, Stringify(
		Call("fn", Val("1"), Val(1), Val('1'), Val(true), Unary(Id("a"))),
	))

	tt.Equal(`make([]string, 1)`, Stringify(
		Call("make", Slice(String), Val(1)),
	))

	tt.Equal(`defer func () {
}()`, Stringify(
		CallWith(Func().Do()).AsDefer(),
	))

	tt.Equal(`defer fn("1", 1, '1', true, &a)`, Stringify(
		Call("fn", Val("1"), Val(1), Val('1'), Val(true), Unary(Id("a"))).AsDefer(),
	))

	tt.Equal(`go fn("1", 1, '1', true, &a)`, Stringify(
		Call("fn", Val("1"), Val(1), Val('1'), Val(true), Unary(Id("a"))).AsGo(),
	))

	tt.Equal(`string("1")`, Stringify(
		Convert(String, Val("1")),
	))
}

func TestSnippet_TypeAssert(t *testing.T) {
	tt := require.New(t)

	tt.Equal(`a.(string)`, Stringify(
		TypeAssert(
			String,
			Id("a"),
		),
	))

	tt.Equal(`a.(interface {
IsZero() (bool)
})`, Stringify(
		TypeAssert(
			Interface(
				Func().Return(Var(Bool)).Named("IsZero"),
			),
			Id("a"),
		),
	))
}
