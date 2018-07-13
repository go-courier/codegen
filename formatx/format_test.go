package formatx

import (
	"os"
	"os/exec"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProcess(t *testing.T) {
	exec.Command("go", "get", "github.com/davecgh/go-spew/spew").CombinedOutput()

	cwd, _ := os.Getwd()
	result, _ := Format(path.Join(cwd, "format2_test.go"), []byte( /* language=go */ `package formatx

import (
	"github.com/go-courier/codegen"

	// spew
	s "github.com/davecgh/go-spew/spew" 
	testing "testing" // testing
)

func Test(t *testing.T) {
	s.Dump(codegen.String)
}
`), SortImportsProcess)

	require.Equal(t /* language=go */, `package formatx

import (
	"testing" // testing

	// spew
	s "github.com/davecgh/go-spew/spew"

	"github.com/go-courier/codegen"
)

func Test(t *testing.T) {
	s.Dump(codegen.String)
}
`, string(result))
}
