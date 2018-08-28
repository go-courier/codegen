package codegen

import (
	"bytes"
	"go/token"
	"sort"
	"strconv"
	"strings"
)

type SnippetSpec interface {
	Snippet
	snippetSpec()
}

func DeclConst(specs ...SnippetSpec) *SnippetTypeDecl {
	return &SnippetTypeDecl{
		Token: token.CONST,
		Specs: specs,
	}
}

func DeclVar(specs ...SnippetSpec) *SnippetTypeDecl {
	return &SnippetTypeDecl{
		Token: token.VAR,
		Specs: specs,
	}
}

func DeclType(specs ...SnippetSpec) *SnippetTypeDecl {
	return &SnippetTypeDecl{
		Token: token.TYPE,
		Specs: specs,
	}
}

type SnippetTypeDecl struct {
	Token token.Token
	Specs []SnippetSpec
}

func (decl *SnippetTypeDecl) Bytes() []byte {
	buf := &bytes.Buffer{}

	buf.WriteString(decl.Token.String())
	buf.WriteRune(' ')

	multi := len(decl.Specs) > 1

	if multi {
		buf.WriteString("(")
		buf.WriteRune('\n')
	}

	for i, spec := range decl.Specs {
		if i > 0 {
			buf.WriteRune('\n')
		}
		buf.Write(spec.Bytes())
	}

	if multi {
		buf.WriteRune('\n')
		buf.WriteString(")")
	}

	return buf.Bytes()
}

func Var(tpe SnippetType, names ...string) *SnippetField {
	return &SnippetField{
		Type:  tpe,
		Names: IdsFromNames(names...),
	}
}

type SnippetField struct {
	SnippetSpec
	SnippetCanAddr
	Type  SnippetType
	Names []*SnippetIdent
	Tag   string
	Alias bool
	SnippetComments
}

func (f SnippetField) AsAlias() *SnippetField {
	f.Alias = true
	return &f
}

func (f SnippetField) WithComments(comments ...string) *SnippetField {
	f.SnippetComments = Comments(comments...)
	return &f
}

func (f SnippetField) WithTag(tag string) *SnippetField {
	f.Tag = tag
	return &f
}

func (f SnippetField) WithTags(tags map[string][]string) *SnippetField {
	buf := &bytes.Buffer{}

	tagNames := make([]string, 0)
	for tag := range tags {
		tagNames = append(tagNames, tag)
	}
	sort.Strings(tagNames)

	for i, tag := range tagNames {
		if i > 0 {
			buf.WriteRune(' ')
		}

		values := make([]string, 0)
		for j := range tags[tag] {
			v := tags[tag][j]
			if v != "" {
				values = append(values, v)
			}
		}

		buf.WriteString(tag)
		buf.WriteRune(':')
		buf.WriteString(strconv.Quote(strings.Join(values, ",")))
	}
	f.Tag = buf.String()
	return &f
}

func (f SnippetField) WithoutTag() *SnippetField {
	f.Tag = ""
	return &f
}

func (f *SnippetField) Bytes() []byte {
	buf := &bytes.Buffer{}

	if f.SnippetComments != nil {
		buf.Write(f.SnippetComments.Bytes())
	}

	for i := range f.Names {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.Write(f.Names[i].Bytes())
	}

	if len(f.Names) > 0 {
		if f.Alias {
			buf.WriteString(" = ")
		} else {
			buf.WriteRune(' ')
		}
	}

	buf.Write(f.Type.Bytes())

	if f.Tag != "" {
		buf.WriteRune(' ')
		buf.WriteRune('`')
		buf.WriteString(f.Tag)
		buf.WriteRune('`')
	}

	return buf.Bytes()
}
