package codegen

import (
	"bytes"
	"fmt"
	"go/token"
	"reflect"
)

type SnippetType interface {
	Snippet
	snippetType()
}

type ImportPathAliaser func(importPath string) string

var TypeOf = createTypeOf(LowerSnakeCase)

func createTypeOf(aliaser ImportPathAliaser) func(tpe reflect.Type) SnippetType {
	return func(tpe reflect.Type) SnippetType {
		if tpe.PkgPath() != "" {
			return Type(aliaser(tpe.PkgPath()) + "." + tpe.Name())
		}

		typeof := createTypeOf(aliaser)

		switch tpe.Kind() {
		case reflect.Ptr:
			return Star(typeof(tpe.Elem()))
		case reflect.Chan:
			return Chan(typeof(tpe.Elem()))
		case reflect.Struct:
			fields := make([]*SnippetField, 0)

			for i := 0; i < tpe.NumField(); i++ {
				f := tpe.Field(i)
				if f.Anonymous {
					fields = append(fields, Var(typeof(f.Type)).WithTag(string(f.Tag)))
				} else {
					fields = append(fields, Var(typeof(f.Type), f.Name).WithTag(string(f.Tag)))
				}
			}

			return Struct(fields...)
		case reflect.Array:
			return Array(typeof(tpe.Elem()), tpe.Len())
		case reflect.Slice:
			return Slice(typeof(tpe.Elem()))
		case reflect.Map:
			return Map(typeof(tpe.Key()), typeof(tpe.Elem()))
		default:
			return BuiltInType(tpe.String())
		}
	}
}

func Ellipsis(tpe SnippetType) *EllipsisType {
	return &EllipsisType{
		Elem: tpe,
	}
}

type EllipsisType struct {
	SnippetType
	Elem SnippetType
}

func (tpe *EllipsisType) Bytes() []byte {
	buf := &bytes.Buffer{}

	buf.WriteString(token.ELLIPSIS.String())
	buf.Write(tpe.Elem.Bytes())

	return buf.Bytes()
}

func Chan(tpe SnippetType) *ChanType {
	return &ChanType{
		Elem: tpe,
	}
}

type ChanType struct {
	SnippetType
	Elem SnippetType
}

func (tpe *ChanType) Bytes() []byte {
	buf := &bytes.Buffer{}

	buf.WriteString(token.CHAN.String() + " ")
	buf.Write(tpe.Elem.Bytes())

	return buf.Bytes()
}

func Type(name string) *NamedType {
	return &NamedType{
		Name: Id(name),
	}
}

type NamedType struct {
	SnippetType
	SnippetCanBeInterfaceMethod
	SnippetCanAddr
	Name *SnippetIdent
}

func (tpe *NamedType) Bytes() []byte {
	return tpe.Name.Bytes()
}

func Func(params ...*SnippetField) *FuncType {
	return &FuncType{
		Params: params,
	}
}

type FuncType struct {
	SnippetType
	SnippetCanBeInterfaceMethod
	Name    *SnippetIdent
	Recv    *SnippetField
	Params  []*SnippetField
	Results []*SnippetField
	Body    []Snippet

	noFuncToken bool
}

func (f FuncType) withoutFuncToken() *FuncType {
	f.noFuncToken = true
	return &f
}

func (f FuncType) Do(bodies ...Snippet) *FuncType {
	f.Body = append([]Snippet{}, bodies...)
	return &f
}

func (f FuncType) Named(name string) *FuncType {
	f.Name = Id(name)
	return &f
}

func (f FuncType) MethodOf(recv *SnippetField) *FuncType {
	f.Recv = recv
	return &f
}

func (f FuncType) Return(results ...*SnippetField) *FuncType {
	f.Results = results
	return &f
}

func (f *FuncType) Bytes() []byte {
	buf := &bytes.Buffer{}

	if !f.noFuncToken {
		buf.WriteString(token.FUNC.String())
		buf.WriteRune(' ')
	}

	if f.Recv != nil {
		buf.WriteByte('(')
		buf.Write(f.Recv.Bytes())
		buf.WriteString(") ")
	}

	if f.Name != nil {
		buf.Write(f.Name.Bytes())
	}

	buf.WriteByte('(')

	for i := range f.Params {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.Write(f.Params[i].WithoutTag().Bytes())
	}

	buf.WriteByte(')')

	hasResults := len(f.Results) > 0

	if hasResults {
		buf.WriteString(" (")
	}

	for i := range f.Results {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.Write(f.Results[i].WithoutTag().Bytes())
	}

	if hasResults {
		buf.WriteByte(')')
	}

	if f.Body != nil {
		buf.WriteRune(' ')
		buf.Write(Body(f.Body).Bytes())
	}

	return buf.Bytes()
}

func Struct(fields ...*SnippetField) *StructType {
	return &StructType{
		Fields: fields,
	}
}

type StructType struct {
	SnippetType
	Fields []*SnippetField
}

func (tpe *StructType) Bytes() []byte {
	buf := &bytes.Buffer{}

	buf.WriteString(token.STRUCT.String() + " {")

	for i := range tpe.Fields {
		buf.WriteRune('\n')
		buf.Write(tpe.Fields[i].Bytes())
	}

	buf.WriteRune('\n')
	buf.WriteRune('}')

	return buf.Bytes()
}

func Interface(methods ...SnippetCanBeInterfaceMethod) *InterfaceType {
	return &InterfaceType{
		Methods: methods,
	}
}

type SnippetCanBeInterfaceMethod interface {
	canBeInterfaceMethod()
}

type InterfaceType struct {
	SnippetType
	Methods []SnippetCanBeInterfaceMethod
}

func (tpe *InterfaceType) Bytes() []byte {
	buf := &bytes.Buffer{}

	buf.WriteString(token.INTERFACE.String() + " {")

	for i := range tpe.Methods {
		if i == 0 {
			buf.WriteRune('\n')
		}
		methodType := tpe.Methods[i]
		switch methodType.(type) {
		case *FuncType:
			buf.Write(methodType.(*FuncType).withoutFuncToken().Bytes())
		case *NamedType:
			buf.Write(methodType.(*NamedType).Bytes())
		}
		buf.WriteRune('\n')
	}

	buf.WriteRune('}')

	return buf.Bytes()
}

func Map(key SnippetType, value SnippetType) *MapType {
	return &MapType{
		Key:   key,
		Value: value,
	}
}

type MapType struct {
	SnippetType
	Key   SnippetType
	Value SnippetType
}

func (tpe *MapType) Bytes() []byte {
	buf := &bytes.Buffer{}

	buf.WriteString(token.MAP.String() + "[")
	buf.Write(tpe.Key.Bytes())
	buf.WriteRune(']')
	buf.Write(tpe.Value.Bytes())

	return buf.Bytes()
}

func Slice(tpe SnippetType) *SliceType {
	return &SliceType{
		Elem: tpe,
	}
}

type SliceType struct {
	SnippetType
	Elem SnippetType
}

func (tpe *SliceType) Bytes() []byte {
	buf := &bytes.Buffer{}

	buf.WriteString("[]")
	buf.Write(tpe.Elem.Bytes())

	return buf.Bytes()
}

func Array(tpe SnippetType, len int) *ArrayType {
	return &ArrayType{
		Elem: tpe,
		Len:  len,
	}
}

type ArrayType struct {
	SnippetType
	Elem SnippetType
	Len  int
}

func (tpe *ArrayType) Bytes() []byte {
	buf := &bytes.Buffer{}

	buf.WriteString(fmt.Sprintf("[%d]", tpe.Len))
	buf.Write(tpe.Elem.Bytes())

	return buf.Bytes()
}

type BuiltInType string

func (BuiltInType) snippetType() {}

func (tpe BuiltInType) Bytes() []byte {
	return []byte(string(tpe))
}

const (
	Bool BuiltInType = "bool"

	Int   BuiltInType = "int"
	Int8  BuiltInType = "int8"
	Int16 BuiltInType = "int16"
	Int32 BuiltInType = "int32"
	Int64 BuiltInType = "int64"

	Uint    BuiltInType = "uint"
	Uint8   BuiltInType = "uint8"
	Uint16  BuiltInType = "uint16"
	Uint32  BuiltInType = "uint32"
	Uint64  BuiltInType = "uint64"
	Uintptr BuiltInType = "uintptr"

	Float32    BuiltInType = "float32"
	Float64    BuiltInType = "float64"
	Complex64  BuiltInType = "complex64"
	Complex128 BuiltInType = "complex128"

	String BuiltInType = "string"
	Byte   BuiltInType = "byte"
	Rune   BuiltInType = "rune"

	Error BuiltInType = "error"
)
