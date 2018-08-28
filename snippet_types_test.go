package codegen

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSnippetTypeOf(t *testing.T) {
	tt := require.New(t)

	tt.Equal("*bytes.Buffer", Stringify(TypeOf(reflect.TypeOf(&bytes.Buffer{}))))
	tt.Equal(`chan struct {
Name string `+"`"+`json:"name"`+"`"+`
}`, Stringify(TypeOf(reflect.TypeOf(make(chan struct {
		Name string `json:"name"`
	})))))

	tt.Equal(`chan struct {
bytes.Buffer
Name string `+"`"+`json:"name"`+"`"+`
}`, Stringify(TypeOf(reflect.TypeOf(make(chan struct {
		bytes.Buffer
		Name string `json:"name"`
	})))))

	tt.Equal("*bytes.Buffer", Stringify(TypeOf(reflect.TypeOf(&bytes.Buffer{}))))
	tt.Equal("[]string", Stringify(TypeOf(reflect.TypeOf([]string{}))))
	tt.Equal("[0]string", Stringify(TypeOf(reflect.TypeOf([0]string{}))))
	tt.Equal("map[string]string", Stringify(TypeOf(reflect.TypeOf(map[string]string{}))))
}

func TestSnippetType(t *testing.T) {
	tt := require.New(t)

	tt.Equal("bool", Stringify(Bool))
	tt.Equal("Name", Stringify(Type("Name")))
}

func TestSnippetType_ModifiedOrComposedTypes(t *testing.T) {
	tt := require.New(t)

	tt.Equal("map[bool]bool", Stringify(Map(Bool, Bool)))
	tt.Equal("[]bool", Stringify(Slice(Bool)))
	tt.Equal("[0]bool", Stringify(Array(Bool, 0)))
	tt.Equal("chan *bool", Stringify(Chan(Star(Bool))))
	tt.Equal("chan bool", Stringify(Chan(Bool)))
	tt.Equal("chan Name", Stringify(Chan(Type("Name"))))
}

func TestSnippetType_InterfaceType(t *testing.T) {
	tt := require.New(t)

	tt.Equal("interface {}", Stringify(Interface()))
	tt.Equal(`interface {
Type()
String() (string)
}`, Stringify(Interface(
		Func().Named("Type"),
		Func().Return(Var(String)).Named("String"),
	)))
	tt.Equal(`interface {
Type
String() (string)
}`, Stringify(Interface(
		Type("Type"),
		Func().Return(Var(String)).Named("String"),
	)))
}

func TestSnippetType_FuncType(t *testing.T) {
	tt := require.New(t)

	tt.Equal("func ()", Stringify(Func()))

	tt.Equal("func Fn()", Stringify(Func().Named("Fn")))

	tt.Equal("func Fn(a, b string)", Stringify(Func(
		Var(String, "a", "b"),
	).Named("Fn")))

	tt.Equal("func Fn(a string, list ...string)", Stringify(Func(
		Var(String, "a"),
		Var(Ellipsis(String), "list"),
	).Named("Fn")))

	tt.Equal("func Fn(a string, b string) (string, error)", Stringify(
		Func(
			Var(String, "a"),
			Var(String, "b"),
		).
			Return(
				Var(String),
				Var(Error),
			).
			Named("Fn"),
	))

	tt.Equal("func (r R) Fn(a, b interface {}) (error)", Stringify(
		Func(
			Var(Interface(), "a", "b"),
		).
			Return(Var(Error)).
			MethodOf(Var(Type("R"), "r")).
			Named("Fn"),
	))
}

func TestSnippetType_StructType(t *testing.T) {
	tt := require.New(t)

	tt.Equal(`struct {
Embed
Key bool `+"`"+`json:"key,omitempty" validate:"@string"`+"`"+`
KeyA, KeyA1 bool
}`, Stringify(Struct(
		Var(Type("Embed")),
		Var(Bool, "Key").WithTags(map[string][]string{
			"json":     {"key", "omitempty"},
			"validate": {"@string"},
		}),
		Var(Bool, "KeyA", "KeyA1"),
	)))
}
