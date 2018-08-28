package formatx

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProcess(t *testing.T) {
	cwd, _ := os.Getwd()
	result, _ := Format(path.Join(cwd, "format2_test.go"), []byte( /* language=go */ `package formatx

import (
	"github.com/go-courier/codegen"
	
	"unicode"
	
	"unicode/utf8"

	// spew
	s "github.com/davecgh/go-spew/spew" 
	"testing" // testing
)

func Test(t *testing.T) {
	s.Dump(codegen.String)
	s.Dump(unicode.Armenian)
	s.Dump(utf8.DecodeLastRune)
}
`), SortImportsProcess)

	fmt.Println(string(result))

	require.Equal(t /* language=go */, `package formatx

import (
	"testing" // testing
	"unicode"
	"unicode/utf8"

	// spew
	s "github.com/davecgh/go-spew/spew"

	"github.com/go-courier/codegen"
)

func Test(t *testing.T) {
	s.Dump(codegen.String)
	s.Dump(unicode.Armenian)
	s.Dump(utf8.DecodeLastRune)
}
`, string(result))
}
