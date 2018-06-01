package codegen

import (
	"bytes"
	"go/token"
)

func Select(clauses ...*SnippetClause) *SnippetSelectStmt {
	return &SnippetSelectStmt{
		Clauses: clauses,
	}
}

type SnippetSelectStmt struct {
	Clauses []*SnippetClause
}

func (stmt *SnippetSelectStmt) Bytes() []byte {
	buf := &bytes.Buffer{}

	buf.WriteString(token.SELECT.String())

	buf.WriteString(" {\n")

	for _, clause := range stmt.Clauses {
		buf.Write(clause.Bytes())
	}

	buf.WriteString("}")

	return buf.Bytes()
}

func Switch(cond Snippet) *SnippetSwitchStmt {
	return &SnippetSwitchStmt{
		Cond: cond,
	}
}

type SnippetSwitchStmt struct {
	Init    Snippet
	Cond    Snippet
	Clauses []*SnippetClause
}

func (stmt SnippetSwitchStmt) InitWith(init Snippet) *SnippetSwitchStmt {
	stmt.Init = init
	return &stmt
}

func (stmt SnippetSwitchStmt) When(clauses ...*SnippetClause) *SnippetSwitchStmt {
	stmt.Clauses = append([]*SnippetClause{}, clauses...)
	return &stmt
}

func (stmt *SnippetSwitchStmt) Bytes() []byte {
	buf := &bytes.Buffer{}

	buf.WriteString(token.SWITCH.String())

	if stmt.Cond != nil {
		if stmt.Init != nil {
			buf.WriteRune(' ')
			buf.Write(stmt.Init.Bytes())
			buf.WriteString(";")
		}

		buf.WriteRune(' ')
		buf.Write(stmt.Cond.Bytes())
	}

	buf.WriteString(" {\n")

	for _, clause := range stmt.Clauses {
		buf.Write(clause.Bytes())
	}

	buf.WriteString("}")

	return buf.Bytes()
}

func Clause(ss ...Snippet) *SnippetClause {
	return &SnippetClause{
		List: ss,
	}
}

type SnippetClause struct {
	List []Snippet
	Body []Snippet
}

func (stmt SnippetClause) Do(bodies ...Snippet) *SnippetClause {
	stmt.Body = bodies
	return &stmt
}

func (stmt *SnippetClause) Bytes() []byte {
	buf := &bytes.Buffer{}

	if len(stmt.List) == 0 {
		buf.WriteString(token.DEFAULT.String())
	} else {
		buf.WriteString(token.CASE.String() + " ")
		for i, s := range stmt.List {
			if i > 0 {
				buf.WriteString(", ")
			}
			buf.Write(s.Bytes())
		}
	}

	buf.WriteRune(':')

	for _, s := range stmt.Body {
		buf.WriteRune('\n')
		buf.Write(s.Bytes())
	}

	buf.WriteRune('\n')

	return buf.Bytes()
}

func ForRange(x Snippet, keyAndValue ...string) *SnippetRangeStmt {
	stmt := &SnippetRangeStmt{
		X: x,
	}

	l := len(keyAndValue)

	if l > 0 && keyAndValue[0] != "" {
		stmt.Key = Id(keyAndValue[0])
	}

	if l > 1 && keyAndValue[1] != "" {
		stmt.Value = Id(keyAndValue[1])
	}

	return stmt
}

type SnippetRangeStmt struct {
	Key   *SnippetIdent
	Value *SnippetIdent
	X     Snippet
	Body  []Snippet
}

func (stmt SnippetRangeStmt) Do(bodies ...Snippet) *SnippetRangeStmt {
	stmt.Body = bodies
	return &stmt
}

func (stmt *SnippetRangeStmt) Bytes() []byte {
	buf := &bytes.Buffer{}

	buf.WriteString(token.FOR.String())
	buf.WriteRune(' ')

	if stmt.Key != nil {
		buf.Write(stmt.Key.Bytes())
	}

	if stmt.Value != nil {
		buf.WriteRune(',')
		buf.WriteRune(' ')
		buf.Write(stmt.Value.Bytes())
	}

	if stmt.Key != nil || stmt.Value != nil {
		buf.WriteString(" " + token.DEFINE.String() + " ")
	}

	buf.WriteString(token.RANGE.String() + " ")
	buf.Write(stmt.X.Bytes())

	buf.WriteRune(' ')
	buf.Write(Body(stmt.Body).Bytes())

	return buf.Bytes()
}

func For(init Snippet, cond Snippet, post Snippet) *SnippetForStmt {
	return &SnippetForStmt{
		Init: init,
		Cond: cond,
		Post: post,
	}
}

type SnippetForStmt struct {
	Init Snippet
	Cond Snippet
	Post Snippet
	Body []Snippet
}

func (stmt SnippetForStmt) Do(bodies ...Snippet) *SnippetForStmt {
	stmt.Body = bodies
	return &stmt
}

func (stmt *SnippetForStmt) Bytes() []byte {
	buf := &bytes.Buffer{}

	buf.WriteString(token.FOR.String())

	if stmt.Init != nil {
		buf.WriteRune(' ')
		buf.Write(stmt.Init.Bytes())
		buf.WriteRune(';')
	}

	if stmt.Cond != nil {
		buf.WriteRune(' ')
		buf.Write(stmt.Cond.Bytes())
	}

	if stmt.Post != nil {
		buf.WriteRune(';')
		buf.WriteRune(' ')
		buf.Write(stmt.Post.Bytes())
	}

	buf.WriteRune(' ')
	buf.Write(Body(stmt.Body).Bytes())

	return buf.Bytes()
}

func If(cond Snippet) *SnippetIfStmt {
	return &SnippetIfStmt{
		Cond: cond,
	}
}

type SnippetIfStmt struct {
	Init     Snippet
	Cond     Snippet
	Body     []Snippet
	ElseList []*SnippetIfStmt
	AsElse   bool
}

func (stmt SnippetIfStmt) InitWith(init Snippet) *SnippetIfStmt {
	stmt.Init = init
	return &stmt
}

func (stmt SnippetIfStmt) WithoutInit() *SnippetIfStmt {
	stmt.Init = nil
	return &stmt
}

func (stmt SnippetIfStmt) Else(ifStmt *SnippetIfStmt) *SnippetIfStmt {
	stmt.ElseList = append(stmt.ElseList, ifStmt.WithoutInit())
	return &stmt
}

func (stmt SnippetIfStmt) Do(bodies ...Snippet) *SnippetIfStmt {
	stmt.Body = bodies
	return &stmt
}

func (stmt *SnippetIfStmt) Bytes() []byte {
	buf := &bytes.Buffer{}

	if stmt.Cond != nil {
		buf.WriteString(token.IF.String())
	}

	if stmt.Init != nil {
		buf.WriteRune(' ')
		buf.Write(stmt.Init.Bytes())
		buf.WriteRune(';')
	}

	if stmt.Cond != nil {
		buf.WriteRune(' ')
		buf.Write(stmt.Cond.Bytes())
	}

	buf.WriteRune(' ')
	buf.Write(Body(stmt.Body).Bytes())

	for _, then := range stmt.ElseList {
		buf.WriteString(" " + token.ELSE.String())
		if then.Cond != nil {
			buf.WriteRune(' ')
		}
		buf.Write(then.Bytes())
	}

	return buf.Bytes()
}

func Define(lhs ...SnippetCanAddr) *SnippetAssignStmt {
	return &SnippetAssignStmt{
		Token: token.DEFINE,
		Lhs:   lhs,
	}
}

func Assign(lhs ...SnippetCanAddr) *SnippetAssignStmt {
	return &SnippetAssignStmt{
		Token: token.ASSIGN,
		Lhs:   lhs,
	}
}

func AssignWith(tok token.Token, lhs ...SnippetCanAddr) *SnippetAssignStmt {
	return &SnippetAssignStmt{
		Token: tok,
		Lhs:   lhs,
	}
}

type SnippetAssignStmt struct {
	SnippetSpec
	Token token.Token
	Lhs   []SnippetCanAddr
	Rhs   []Snippet
}

func (stmt SnippetAssignStmt) By(rhs ...Snippet) *SnippetAssignStmt {
	stmt.Rhs = rhs
	return &stmt
}

func (stmt *SnippetAssignStmt) Bytes() []byte {
	buf := &bytes.Buffer{}

	for i, n := range stmt.Lhs {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.Write(n.Bytes())
	}

	if len(stmt.Rhs) > 0 {
		buf.WriteRune(' ')
		buf.WriteString(stmt.Token.String())
		buf.WriteRune(' ')

		for i, n := range stmt.Rhs {
			if i > 0 {
				buf.WriteString(", ")
			}
			buf.Write(n.Bytes())
		}

		return buf.Bytes()
	}

	return buf.Bytes()
}

func Return(snippets ...Snippet) *SnippetReturnStmt {
	return &SnippetReturnStmt{
		Results: snippets,
	}
}

type SnippetReturnStmt struct {
	Results []Snippet
}

func (stmt *SnippetReturnStmt) Bytes() []byte {
	buf := &bytes.Buffer{}

	buf.WriteString("return")

	for i, n := range stmt.Results {
		if i > 0 {
			buf.WriteRune(',')
		}
		buf.WriteRune(' ')
		buf.Write(n.Bytes())
	}

	return buf.Bytes()
}
