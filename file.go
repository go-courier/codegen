package codegen

import (
	"bytes"
	"fmt"
	"go/build"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
)

func NewFile(pkgName string, filename string) *File {
	return &File{
		PkgName:  LowerSnakeCase(pkgName),
		filename: filename,
	}
}

type File struct {
	PkgName  string
	filename string
	imports  map[string]string
	bytes.Buffer
}

func (file *File) WriteBlock(ss ...Snippet) {
	for _, s := range ss {
		file.Write(s.Bytes())
		file.WriteString("\n\n")
	}
}

func (file *File) Bytes() []byte {
	buf := &bytes.Buffer{}

	buf.WriteString(`package ` + LowerSnakeCase(file.PkgName) + `
`)

	if file.imports != nil {
		buf.WriteString(`import (
`)
		for importPath, alias := range file.imports {
			buf.WriteString(alias)
			buf.WriteString(" ")
			buf.WriteString(strconv.Quote(importPath))
			buf.WriteString("\n")
		}

		buf.WriteString(`)
`)
	}

	io.Copy(buf, &file.Buffer)

	return Format(file.filename, buf.Bytes())
}

func Format(filename string, src []byte) []byte {
	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, filename, src, parser.ParseComments)
	if err != nil {
		for i, line := range bytes.Split(src, []byte("\n")) {
			fmt.Printf("\t%d\t%s\n", i+1, line)
		}
		panic(fmt.Errorf("go codes parse failed: %s in %s", err.Error(), filename))
	}
	buf := &bytes.Buffer{}
	if err := format.Node(buf, fileSet, file); err != nil {
		for i, line := range bytes.Split(src, []byte("\n")) {
			fmt.Printf("\t%d\t%s\n", i+1, line)
		}
		panic(fmt.Errorf("go codes format failed: %s in %s", err.Error(), filename))
	}
	return buf.Bytes()
}

func (file *File) Expr(f string, args ...interface{}) SnippetExpr {
	return createExpr(file.importAliaser)(f, args...)
}

func (file *File) TypeOf(tpe reflect.Type) SnippetType {
	return createTypeOf(file.importAliaser)(tpe)
}

func (file *File) Val(v interface{}) Snippet {
	return createVal(file.importAliaser)(v)
}

func (file *File) importAliaser(importPath string) string {
	if file.imports == nil {
		file.imports = map[string]string{}
	}
	if file.imports[importPath] == "" {
		pkg, err := build.Import(importPath, "", build.ImportComment)
		if err != nil {
			panic(err)
		}
		importPath = deVendor(pkg.ImportPath)
		file.imports[importPath] = LowerSnakeCase(importPath)
	}
	return file.imports[importPath]
}

func (file *File) Use(importPath string, exposedName string) string {
	return file.importAliaser(importPath) + "." + exposedName
}

func deVendor(importPath string) string {
	parts := strings.Split(importPath, "/vendor/")
	return parts[len(parts)-1]
}

func (file *File) WriteFile() (int, error) {
	dir := filepath.Dir(file.filename)

	if dir != "" {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return -1, err
		}
	}

	f, err := os.Create(file.filename)
	defer f.Close()
	if err != nil {
		return -1, err
	}

	n3, err := f.Write(file.Bytes())
	if err != nil {
		return -1, err
	}

	if err := f.Sync(); err != nil {
		return -1, err
	}

	return n3, nil
}
