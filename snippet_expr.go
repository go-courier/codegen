package codegen

import (
	"bytes"
	"go/token"
	"regexp"
)

type SnippetExpr string

func (tpe SnippetExpr) Bytes() []byte {
	return []byte(string(tpe))
}

var Expr = createExpr(LowerSnakeCase)

func createExpr(aliaser ImportPathAliaser) func(f string, args ...interface{}) SnippetExpr {
	val := createVal(aliaser)

	return func(f string, args ...interface{}) SnippetExpr {
		idx := 0
		return SnippetExpr(reExprHolder.ReplaceAllStringFunc(f, func(i string) string {
			arg := args[idx]
			idx++
			if s, ok := arg.(Snippet); ok {
				return Stringify(s)
			}
			return Stringify(val(arg))
		}))
	}
}

var reExprHolder = regexp.MustCompile(`(\$\d+)|\?`)

func KeyValue(key Snippet, value Snippet) *SnippetKeyValueExpr {
	return &SnippetKeyValueExpr{
		Key:   key,
		Value: value,
	}
}

type SnippetKeyValueExpr struct {
	Snippet
	Key   Snippet
	Value Snippet
}

func (tpe *SnippetKeyValueExpr) Bytes() []byte {
	buf := &bytes.Buffer{}

	buf.Write(tpe.Key.Bytes())
	buf.WriteString(token.COLON.String())
	buf.WriteRune(' ')
	buf.Write(tpe.Value.Bytes())

	return buf.Bytes()
}

func Sel(x Snippet, selectors ...Snippet) *SnippetSelectorExpr {
	return &SnippetSelectorExpr{
		X:         x,
		Selectors: selectors,
	}
}

type SnippetSelectorExpr struct {
	X         Snippet
	Selectors []Snippet
}

func (expr *SnippetSelectorExpr) Bytes() []byte {
	buf := &bytes.Buffer{}

	buf.Write(expr.X.Bytes())

	for _, selectorExpr := range expr.Selectors {
		buf.WriteRune('.')
		buf.Write(selectorExpr.Bytes())
	}

	return buf.Bytes()
}

type SnippetCanAddr interface {
	Snippet
	snippetCanAddr()
}

func Star(tpe SnippetType) *SnippetStarExpr {
	return &SnippetStarExpr{
		X: tpe,
	}
}

type SnippetStarExpr struct {
	SnippetCanAddr
	SnippetType
	X SnippetType
}

func (expr *SnippetStarExpr) Bytes() []byte {
	buf := &bytes.Buffer{}

	buf.WriteString(token.MUL.String())
	buf.Write(expr.X.Bytes())

	return buf.Bytes()
}

func Unary(addr SnippetCanAddr) *SnippetUnaryExpr {
	return &SnippetUnaryExpr{
		Elem: addr,
	}
}

type SnippetUnaryExpr struct {
	Elem SnippetCanAddr
}

func (tpe *SnippetUnaryExpr) Bytes() []byte {
	buf := &bytes.Buffer{}

	buf.WriteString(token.AND.String())
	buf.Write(tpe.Elem.Bytes())

	return buf.Bytes()
}

func Paren(s Snippet) *SnippetParenExpr {
	return &SnippetParenExpr{
		Elem: s,
	}
}

type SnippetParenExpr struct {
	SnippetCanAddr
	Elem Snippet
}

func (tpe *SnippetParenExpr) Bytes() []byte {
	buf := &bytes.Buffer{}

	buf.WriteRune('(')
	buf.Write(tpe.Elem.Bytes())
	buf.WriteRune(')')

	return buf.Bytes()
}

func Convert(tpe SnippetType, target Snippet) *SnippetCallExpr {
	return CallWith(tpe, target)
}

func CallWith(s Snippet, params ...Snippet) *SnippetCallExpr {
	return &SnippetCallExpr{
		X:      s,
		Params: params,
	}
}

func Call(name string, params ...Snippet) *SnippetCallExpr {
	if isBuiltInFunc(name) {
		return &SnippetCallExpr{
			X:      SnippetIdent(name),
			Params: params,
		}
	}
	return &SnippetCallExpr{
		X:      Id(name),
		Params: params,
	}
}

type SnippetCallExpr struct {
	X        Snippet
	Params   []Snippet
	Ellipsis bool
	Modifier token.Token
}

func (expr SnippetCallExpr) AsDefer() *SnippetCallExpr {
	expr.Modifier = token.DEFER
	return &expr
}

func (expr SnippetCallExpr) AsGo() *SnippetCallExpr {
	expr.Modifier = token.GO
	return &expr
}

func (expr SnippetCallExpr) WithEllipsis() *SnippetCallExpr {
	expr.Ellipsis = true
	return &expr
}

func (expr *SnippetCallExpr) Bytes() []byte {
	buf := &bytes.Buffer{}

	if expr.Modifier > 0 {
		buf.WriteString(expr.Modifier.String())
		buf.WriteRune(' ')
	}

	buf.Write(expr.X.Bytes())

	buf.WriteRune('(')

	for i, p := range expr.Params {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.Write(p.Bytes())
	}

	if expr.Ellipsis {
		buf.WriteString(token.ELLIPSIS.String())
	}

	buf.WriteRune(')')

	return buf.Bytes()
}

func TypeAssert(tpe SnippetType, x Snippet) *SnippetTypeAssertExpr {
	return &SnippetTypeAssertExpr{
		X:    x,
		Type: tpe,
	}
}

type SnippetTypeAssertExpr struct {
	X    Snippet
	Type SnippetType
}

func (expr *SnippetTypeAssertExpr) Bytes() []byte {
	buf := &bytes.Buffer{}

	buf.Write(expr.X.Bytes())

	buf.WriteRune('.')
	buf.WriteRune('(')
	buf.Write(expr.Type.Bytes())
	buf.WriteRune(')')

	return buf.Bytes()
}
