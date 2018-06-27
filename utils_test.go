package codegen

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGeneratedFileSuffix(t *testing.T) {
	cases := []struct {
		from string
		to   string
	}{
		{
			"./main.go",
			"./main__generated.go",
		},
		{
			"./main_test.go",
			"./main__generated_test.go",
		},
	}

	for _, c := range cases {
		assert.Equal(t, GeneratedFileSuffix(c.from), c.to)
	}
}

func TestIsEmptyValue(t *testing.T) {
	tt := assert.New(t)

	c := make(chan int)

	tt.False(IsEmptyValue(reflect.ValueOf(c)))
	tt.True(IsEmptyValue(reflect.ValueOf(0)))
	tt.True(IsEmptyValue(reflect.ValueOf(float32(0))))
	tt.True(IsEmptyValue(reflect.ValueOf("")))
	tt.True(IsEmptyValue(reflect.ValueOf(false)))
	tt.True(IsEmptyValue(reflect.ValueOf((*int)(nil))))
	tt.True(IsEmptyValue(reflect.ValueOf(uint(0))))
	tt.True(IsEmptyValue(reflect.ValueOf(time.Time{})))
}
