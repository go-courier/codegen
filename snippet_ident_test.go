package codegen

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSnippetIdent(t *testing.T) {
	tt := require.New(t)

	tt.Equal("Test", Stringify(Id("Test")))

	invalidIdFuncs := []func(){
		func() { Id("$adsf") },
		func() { Id("1-21asd") },
		func() { Id("") },
		func() { Id("type") },
		func() { Id("time.%Time") },
	}

	for i := range invalidIdFuncs {
		f := invalidIdFuncs[i]
		err := TryCatch(f)
		tt.Error(err)
	}
}

func TestSnippetIdent_Converts(t *testing.T) {
	tt := require.New(t)

	id := Id("i_am_an_id")

	tt.Equal(Id("i_am_an_id"), id.LowerSnakeCase())
	tt.Equal(Id("I_AM_AN_ID"), id.UpperSnakeCase())
	tt.Equal(Id("iAmAnID"), id.LowerCamelCase())
	tt.Equal(Id("IAmAnID"), id.UpperCamelCase())
}

func Test_splitToWords(t *testing.T) {
	tt := require.New(t)

	tt.Equal([]string{}, splitToWords(""))
	tt.Equal([]string{"lowercase"}, splitToWords("lowercase"))
	tt.Equal([]string{"Class"}, splitToWords("Class"))
	tt.Equal([]string{"My", "Class"}, splitToWords("MyClass"))
	tt.Equal([]string{"My", "C"}, splitToWords("MyC"))
	tt.Equal([]string{"HTML"}, splitToWords("HTML"))
	tt.Equal([]string{"PDF", "Loader"}, splitToWords("PDFLoader"))
	tt.Equal([]string{"A", "String"}, splitToWords("AString"))
	tt.Equal([]string{"Simple", "XML", "Parser"}, splitToWords("SimpleXMLParser"))
	tt.Equal([]string{"vim", "RPC", "Plugin"}, splitToWords("vimRPCPlugin"))
	tt.Equal([]string{"GL11", "Version"}, splitToWords("GL11Version"))
	tt.Equal([]string{"99", "Bottles"}, splitToWords("99Bottles"))
	tt.Equal([]string{"May5"}, splitToWords("May5"))
	tt.Equal([]string{"BFG9000"}, splitToWords("BFG9000"))
	tt.Equal([]string{"Böse", "Überraschung"}, splitToWords("BöseÜberraschung"))
	tt.Equal([]string{"Two", "spaces"}, splitToWords("Two  spaces"))
	tt.Equal([]string{"BadUTF8\xe2\xe2\xa1"}, splitToWords("BadUTF8\xe2\xe2\xa1"))
	tt.Equal([]string{"snake", "case"}, splitToWords("snake_case"))
	tt.Equal([]string{"snake", "case"}, splitToWords("snake_ case"))
}
