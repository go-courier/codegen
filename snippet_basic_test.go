package codegen

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExpr(t *testing.T) {
	tt := require.New(t)
	tt.Equal(`1 + 1`, Stringify(Expr("? + ?", 1, 1)))
	tt.Equal(`"1" + "1"`, Stringify(Expr("? + ?", "1", "1")))
	tt.Equal(`i++`, Stringify(Expr("?++", Id("i"))))
}

func TestComments(t *testing.T) {
	tt := require.New(t)
	tt.Equal(`// 123123
// 123123
`, Stringify(Comments("123123", "123123")))
}

func TestVal(t *testing.T) {

	tt := require.New(t)

	tt.Error(TryCatch(func() {
		v := make(chan string)
		Val(v)
	}))

	tt.Equal(`"string"`, Stringify(Val("string")))
	tt.Equal(`1`, Stringify(Val(1)))
	tt.Equal(`1.2`, Stringify(Val(1.2)))
	tt.Equal(`1.2`, Stringify(Val(float32(1.2))))

	tt.Equal(`'b'`, Stringify(Val('b')))
	tt.Equal(`'1'`, Stringify(Val('1')))
	tt.Equal(`1`, Stringify(Val(int32(1))))
	tt.Equal(`true`, Stringify(Val(true)))

	tt.Equal(`[]string{
"1",
"2",
}`, Stringify(Val([]string{"1", "2"})))

	tt.Equal(`struct {
Name string `+"`"+`json:"name"`+"`"+`
Empty string
}{
Name: "123",
}`, Stringify(Val(struct {
		Name  string `json:"name"`
		Empty string
	}{
		Name: "123",
	})))

	{
		type Nested struct {
			Name string
			*Nested
		}

		n := Nested{}
		n.Name = "string"
		n.Nested = &Nested{
			Name: "string",
		}

		tt.Equal(`github_com_go_courier_codegen.Nested{
Name: "string",
Nested: &(github_com_go_courier_codegen.Nested{
Name: "string",
}),
}`,
			Stringify(Val(n)))
	}

	tt.Equal(`[]interface {}{
"1",
nil,
}`, Stringify(Val([]interface{}{"1", nil})))

	tt.Equal(`[2]interface {}{
"1",
nil,
}`, Stringify(Val([2]interface{}{"1", nil})))

	tt.Equal(`map[string]int{
"1": 1,
"2": 2,
}`, Stringify(Val(map[string]int{"2": 2, "1": 1})))
}

func TestSnippet_Block(t *testing.T) {
	tt := require.New(t)

	tt.Equal(`{
a := Fn("1")
v := Fn("1")
}`, Stringify(
		Block(
			Define(Id("a")).By(Call("Fn", Val("1"))),
			Define(Id("v")).By(Call("Fn", Val("1"))),
		),
	))
}
