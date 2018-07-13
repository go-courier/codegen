package formatx

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
)

type Process func(fset *token.FileSet, f *ast.File, filename string) error

func MustFormat(filename string, src []byte, processes ...Process) []byte {
	codes, err := Format(filename, src, processes...)
	if err != nil {
		panic(err)
	}
	return codes
}

func Format(filename string, src []byte, processes ...Process) ([]byte, error) {
	fset, f, err := ParseFile(filename, src)
	if err != nil {
		return nil, err
	}

	if processes != nil {
		for i := range processes {
			if err := processes[i](fset, f, filename); err != nil {
				return nil, err
			}
		}
	}

	buf := bytes.NewBuffer(nil)
	if err := format.Node(buf, fset, f); err != nil {
		printCodesWithLineNumber(src)
		return nil, fmt.Errorf("go codes format failed: %s in %s", err.Error(), filename)
	}
	return buf.Bytes(), nil
}

func ParseFile(filename string, src []byte) (*token.FileSet, *ast.File, error) {
	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, filename, src, parser.ParseComments)
	if err != nil {
		printCodesWithLineNumber(src)
		return nil, nil, fmt.Errorf("go codes parse failed: %s in %s", err.Error(), filename)
	}
	return fileSet, file, nil
}

func printCodesWithLineNumber(src []byte) {
	for i, line := range bytes.Split(src, []byte("\n")) {
		fmt.Printf("\t%d\t%s\n", i+1, line)
	}
}
