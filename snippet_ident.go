package codegen

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

func Id(n string) *SnippetIdent {
	values := strings.Split(n, ".")

	if !IsValidIdent(values[0]) {
		panic(fmt.Errorf("`%s` is not a valid identifier", values[0]))
	}

	if len(values) == 2 {
		if !IsValidIdent(values[1]) {
			panic(fmt.Errorf("`%s` is not a valid identifier", values[1]))
		}
	}

	ident := SnippetIdent(n)
	return &ident
}

func IdsFromNames(names ...string) []*SnippetIdent {
	ids := make([]*SnippetIdent, 0)

	for _, name := range names {
		ids = append(ids, Id(name))
	}

	return ids
}

type SnippetIdent string

func (SnippetIdent) snippetCanAddr() {}

func (id SnippetIdent) Bytes() []byte {
	return []byte(id)
}

func IsValidIdent(s string) bool {
	if len(s) == 0 {
		return false
	}
	if isReservedWord(s) {
		return false
	}
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' {
			return false
		}
	}

	return true
}

func isBuiltInFunc(alias string) bool {
	for _, name := range builtInFuncs {
		if alias == name {
			return true
		}
	}
	return false
}

func isReservedWord(alias string) bool {
	for _, name := range reserved {
		if alias == name {
			return true
		}
	}
	return false
}

var builtInFuncs = []string{
	"append",
	"complex",
	"cap",
	"close",
	"copy",
	"delete",
	"imag",
	"len",
	"make",
	"new",
	"panic",
	"print",
	"println",
	"real",
	"recover",
}

var reserved = append([]string{
	Stringify(Break),
	Stringify(Fallthrough),
	Stringify(Continue),

	"default",
	"func",
	"interface",
	"select",
	"case",
	"defer",
	"go",
	"map",
	"struct",
	"chan",
	"else",
	"goto",
	"package",
	"switch",
	"const",
	"if",
	"range",
	"type",
	"for",
	"import",
	"return",
	"var",

	Stringify(Bool),
	Stringify(Int),
	Stringify(Int8),
	Stringify(Int16),
	Stringify(Int32),
	Stringify(Int64),

	Stringify(Uint),
	Stringify(Uint8),
	Stringify(Uint16),
	Stringify(Uint32),
	Stringify(Uint64),
	Stringify(Uintptr),

	Stringify(Float32),
	Stringify(Float64),
	Stringify(Complex64),
	Stringify(Complex128),

	Stringify(String),
	Stringify(Byte),
	Stringify(Rune),

	string(Error),
}, builtInFuncs...)

func (id SnippetIdent) UpperCamelCase() *SnippetIdent {
	return Id(UpperCamelCase(string(id)))
}

func (id SnippetIdent) LowerCamelCase() *SnippetIdent {
	return Id(LowerCamelCase(string(id)))
}

func (id SnippetIdent) UpperSnakeCase() *SnippetIdent {
	return Id(UpperSnakeCase(string(id)))
}

func (id SnippetIdent) LowerSnakeCase() *SnippetIdent {
	return Id(LowerSnakeCase(string(id)))
}

func UpperSnakeCase(s string) string {
	return rewords(s, func(result string, word string, idx int) string {
		newWord := strings.ToUpper(word)
		if idx == 0 || (len(newWord) == 1 && unicode.IsDigit(rune(newWord[0]))) {
			return result + newWord
		}
		return result + "_" + newWord
	})
}

func LowerSnakeCase(s string) string {
	return rewords(s, func(result string, word string, idx int) string {
		newWord := strings.ToLower(word)
		if idx == 0 || (len(newWord) == 1 && unicode.IsDigit(rune(newWord[0]))) {
			return result + newWord
		}
		return result + "_" + newWord
	})
}

func UpperCamelCase(s string) string {
	return rewords(s, func(result string, word string, idx int) string {
		return result + camelCase(word)
	})
}

func LowerCamelCase(s string) string {
	return rewords(s, func(result string, word string, idx int) string {
		if idx == 0 {
			return result + strings.ToLower(word)
		}
		return result + camelCase(word)
	})
}

func upperFirst(s string) string {
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

func camelCase(s string) string {
	upperString := strings.ToUpper(s)
	if commonInitialisms[upperString] {
		return upperString
	}
	return upperFirst(strings.ToLower(s))
}

func rewords(s string, reducer func(result string, word string, index int) string) string {
	words := splitToWords(string(s))

	var result = ""

	for idx, word := range words {
		result = reducer(result, word, idx)
	}

	return result
}

func splitToWords(s string) (entries []string) {
	if !utf8.ValidString(s) {
		return []string{s}
	}
	entries = []string{}
	var runes [][]rune
	lastClass := 0
	class := 0

	// split into fields based on class of unicode character
	for _, r := range s {
		switch true {
		case unicode.IsSpace(r):
			class = 1
		case unicode.IsLower(r):
			class = 2
		case unicode.IsUpper(r):
			class = 3
		case unicode.IsDigit(r):
			class = 4
		default:
			class = 5
		}
		if class == lastClass {
			runes[len(runes)-1] = append(runes[len(runes)-1], r)
		} else {
			runes = append(runes, []rune{r})
		}
		lastClass = class
	}

	// handle upper case -> lower case sequences, e.g.
	// "PDFL", "oader" -> "PDF", "Loader"
	for i := 0; i < len(runes)-1; i++ {
		if unicode.IsUpper(runes[i][0]) && unicode.IsLower(runes[i+1][0]) {
			runes[i+1] = append([]rune{runes[i][len(runes[i])-1]}, runes[i+1]...)
			runes[i] = runes[i][:len(runes[i])-1]
		}
	}

	// construct []string from results
	for i, s := range runes {
		if len(s) > 0 {
			if unicode.IsDigit(s[0]) {
				if i > 0 {
					entries[len(entries)-1] += string(s)
				} else {
					entries = append(entries, string(s))
				}
			}
			if unicode.IsLetter(s[0]) {
				entries = append(entries, string(s))
			}
		}
	}

	return
}

// https://github.com/golang/lint/blob/master/lint.go#L749-L788
var commonInitialisms = map[string]bool{
	"ACL":   true,
	"API":   true,
	"ASCII": true,
	"CPU":   true,
	"CSS":   true,
	"DNS":   true,
	"EOF":   true,
	"GUID":  true,
	"HTML":  true,
	"HTTP":  true,
	"HTTPS": true,
	"ID":    true,
	"IP":    true,
	"JSON":  true,
	"LHS":   true,
	"QPS":   true,
	"RAM":   true,
	"RHS":   true,
	"RPC":   true,
	"SLA":   true,
	"SMTP":  true,
	"SQL":   true,
	"SSH":   true,
	"TCP":   true,
	"TLS":   true,
	"TTL":   true,
	"UDP":   true,
	"UI":    true,
	"UID":   true,
	"UUID":  true,
	"URI":   true,
	"URL":   true,
	"UTF8":  true,
	"VM":    true,
	"XML":   true,
	"XMPP":  true,
	"XSRF":  true,
	"XSS":   true,
}
