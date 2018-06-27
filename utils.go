package codegen

import (
	"fmt"
	"path/filepath"
	"reflect"
	"strings"
)

func IsGoFile(filename string) bool {
	return filepath.Ext(filename) == ".go"
}

func IsGoTestFile(filename string) bool {
	return strings.HasSuffix(filepath.Base(filename), "_test.go")
}

func GeneratedFileSuffix(filename string) string {
	dir := filepath.Dir(filename)
	base := filepath.Base(filename)
	ext := filepath.Ext(filename)

	if IsGoFile(filename) && IsGoTestFile(filename) {
		base = strings.Replace(base, "_test.go", "__generated_test.go", -1)
	} else {
		base = strings.Replace(base, ext, fmt.Sprintf("__generated%s", ext), -1)

	}
	return fmt.Sprintf("%s/%s", dir, base)
}

func IsEmptyValue(rv reflect.Value) bool {
	if rv.IsValid() && rv.CanInterface() {
		if zeroChecker, ok := rv.Interface().(interface{ IsZero() bool }); ok {
			return zeroChecker.IsZero()
		}
	}

	switch rv.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return rv.Len() == 0
	case reflect.Bool:
		return !rv.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return rv.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return rv.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return rv.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return rv.IsNil()
	}
	return false
}

func TryCatch(f func()) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%v", e)
		}
	}()

	f()

	return nil
}
