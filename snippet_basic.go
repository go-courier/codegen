package codegen

import (
	"bytes"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

func Stringify(snippet Snippet) string {
	return string(snippet.Bytes())
}

type Snippet interface {
	Bytes() []byte
}

func Block(bodies ...Snippet) Body {
	return Body(bodies)
}

type Body []Snippet

func (body Body) Bytes() []byte {
	buf := &bytes.Buffer{}

	buf.WriteRune('{')

	for _, stmt := range body {
		if stmt == nil {
			continue
		}
		buf.WriteRune('\n')
		buf.Write(stmt.Bytes())
	}

	buf.WriteRune('\n')
	buf.WriteRune('}')

	return buf.Bytes()
}

type SnippetBuiltIn string

func (tpe SnippetBuiltIn) Bytes() []byte {
	return []byte(string(tpe))
}

const (
	Iota        SnippetBuiltIn = "iota"
	True        SnippetBuiltIn = "true"
	False       SnippetBuiltIn = "false"
	Nil         SnippetBuiltIn = "nil"
	Break       SnippetBuiltIn = "break"
	Continue    SnippetBuiltIn = "continue"
	Fallthrough SnippetBuiltIn = "fallthrough"
)

func Comments(lines ...string) SnippetComments {
	finalLines := make([]string, 0)
	for _, line := range lines {
		if line != "" {
			finalLines = append(finalLines, strings.Split(line, "\n")...)
		}
	}
	return SnippetComments(finalLines)
}

type SnippetComments []string

func (comments SnippetComments) Bytes() []byte {
	buf := &bytes.Buffer{}

	for _, n := range comments {
		buf.WriteString("// ")
		buf.WriteString(n)
		buf.WriteRune('\n')
	}

	return buf.Bytes()
}

var Val = createVal(LowerSnakeCase)

func createVal(aliaser ImportPathAliaser) func(v interface{}) Snippet {
	return func(v interface{}) Snippet {
		rv := reflect.ValueOf(v)
		tpe := reflect.TypeOf(v)

		val := createVal(aliaser)
		typeof := createTypeOf(aliaser)

		switch rv.Kind() {
		case reflect.Ptr:
			return Unary(Paren(val(rv.Elem().Interface())))
		case reflect.Struct:
			values := make([]Snippet, 0)
			for i := 0; i < rv.NumField(); i++ {
				f := rv.Field(i)
				ft := tpe.Field(i)

				if !IsEmptyValue(f) {
					values = append(values, KeyValue(Id(ft.Name), val(f.Interface())))
				}
			}
			return Compose(typeof(tpe), values...)
		case reflect.Map:
			values := make([]Snippet, 0)
			for _, key := range rv.MapKeys() {
				values = append(values, KeyValue(val(key.Interface()), val(rv.MapIndex(key).Interface())))
			}
			sort.Slice(values, func(i, j int) bool {
				return string(values[i].(*SnippetKeyValueExpr).Key.Bytes()) < string(values[j].(*SnippetKeyValueExpr).Key.Bytes())
			})
			return Compose(typeof(tpe), values...)
		case reflect.Slice, reflect.Array:
			values := make([]Snippet, 0)
			for i := 0; i < rv.Len(); i++ {
				values = append(values, val(rv.Index(i).Interface()))
			}
			return Compose(typeof(tpe), values...)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int64,
			reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint8:
			return Lit(fmt.Sprintf("%d", v))
		case reflect.Int32:
			if b, ok := v.(rune); ok {
				r := strconv.QuoteRune(b)
				if len(r) == 3 {
					return Lit(r)
				}
			}
			return Lit(fmt.Sprintf("%d", v))
		case reflect.Bool:
			return Lit(strconv.FormatBool(v.(bool)))
		case reflect.Float32:
			return Lit(strconv.FormatFloat(float64(v.(float32)), 'f', -1, 32))
		case reflect.Float64:
			return Lit(strconv.FormatFloat(v.(float64), 'f', -1, 64))
		case reflect.String:
			return Lit(strconv.Quote(v.(string)))
		case reflect.Invalid:
			return Nil
		default:
			panic(fmt.Errorf("%v is an unsupported type", v))
		}
	}
}

func Compose(tpe SnippetType, elts ...Snippet) *SnippetCompositeLit {
	return &SnippetCompositeLit{
		Type: tpe,
		Elts: elts,
	}
}

type SnippetCompositeLit struct {
	Type SnippetType
	Elts []Snippet
}

func (lit *SnippetCompositeLit) Bytes() []byte {
	buf := &bytes.Buffer{}

	if lit.Type != nil {
		buf.Write(lit.Type.Bytes())
	}

	buf.WriteRune('{')

	for _, n := range lit.Elts {
		buf.WriteRune('\n')
		buf.Write(n.Bytes())
		buf.WriteRune(',')
	}

	buf.WriteRune('\n')
	buf.WriteRune('}')

	return buf.Bytes()
}

func Lit(s string) *SnippetLit {
	lit := SnippetLit(s)
	return &lit
}

type SnippetLit string

func (lit SnippetLit) Bytes() []byte {
	return []byte(lit)
}
