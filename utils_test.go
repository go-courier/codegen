package codegen

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

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
