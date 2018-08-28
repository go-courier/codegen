package formatx

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/token"
	"io/ioutil"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

func SortImportsProcess(fset *token.FileSet, f *ast.File, filename string) error {
	ast.SortImports(fset, f)
	dir := path.Dir(filename)

	for _, decl := range f.Decls {
		d, ok := decl.(*ast.GenDecl)
		if !ok || d.Tok != token.IMPORT || len(d.Specs) == 0 {
			break
		}

		g := &groupSet{}

		for i := range d.Specs {
			g.register(d.Specs[i].(*ast.ImportSpec), dir)
		}

		fileSet, file, err := ParseFile(filename, bytes.Replace(formatNode(fset, f), formatNode(fset, d), g.Bytes(), 1))
		if err != nil {
			return err
		}
		*fset = *fileSet
		*f = *file
	}
	return nil
}

func formatNode(fileSet *token.FileSet, node ast.Node) []byte {
	buf := &bytes.Buffer{}
	if err := format.Node(buf, fileSet, node); err != nil {
		panic(err)
	}
	return buf.Bytes()
}

type groupSet [4][]*dep

type dep struct {
	name       string
	pkgPath    string
	importSpec *ast.ImportSpec
}

func (group *groupSet) Bytes() []byte {
	buf := bytes.NewBuffer(nil)

	buf.WriteString("import (")

	for _, deps := range group {
		for _, d := range deps {
			buf.WriteRune('\n')

			importSpec := d.importSpec
			if importSpec.Doc != nil {
				for _, c := range importSpec.Doc.List {
					buf.WriteString(c.Text)
					buf.WriteRune('\n')
				}
			}
			if importSpec.Name != nil {
				buf.WriteString(importSpec.Name.String())
				buf.WriteRune(' ')
			}
			buf.WriteString(importSpec.Path.Value)
			if importSpec.Comment != nil {
				for _, c := range importSpec.Comment.List {
					buf.WriteString(c.Text)
				}
			}
		}
		buf.WriteRune('\n')
	}

	buf.WriteRune(')')
	return buf.Bytes()
}

type stdLibSet map[string]bool

func (s stdLibSet) read(dir string, prefix string) {
	files, _ := ioutil.ReadDir(dir)

	for _, f := range files {
		if f.IsDir() {
			importPath := f.Name()
			if prefix != "" {
				importPath = filepath.Join(prefix, f.Name())
			}
			s.read(filepath.Join(dir, f.Name()), importPath)
			s[importPath] = true
			continue
		}
	}
}

var stdLibs = func() map[string]bool {
	set := stdLibSet{}
	set.read(runtime.GOROOT()+"/src", "")
	return set
}()

func (group *groupSet) register(importSpec *ast.ImportSpec, dir string) {
	importPath, _ := strconv.Unquote(importSpec.Path.Value)

	appendTo := func(i int) {
		group[i] = append(group[i], &dep{
			pkgPath:    importPath,
			importSpec: importSpec,
		})
	}

	// std
	if stdLibs[strings.ToLower(importPath)] {
		appendTo(0)
		return
	}

	// local
	if strings.Index(strings.ToLower(dir), strings.ToLower(importPath)) > -1 {
		appendTo(2)
		return
	}

	// vendor
	appendTo(1)
	return
}
