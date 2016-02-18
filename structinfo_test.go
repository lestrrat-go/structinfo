package structinfo_test

import (
	"reflect"
	"testing"

	"github.com/lestrrat/go-structinfo"
	"github.com/stretchr/testify/assert"
)

type X struct {
	private int
	Foo     string `json:"foo"`
	Bar     string `json:"bar,omitempty"`
	Baz     string `json:"baz"`
}

func TestStructFields(t *testing.T) {
	fields := make(map[string]struct{})
	for _, name := range structinfo.JSONFieldsFromStruct(reflect.ValueOf(X{})) {
		fields[name] = struct{}{}
	}

	expected := map[string]struct{}{
		"foo": {},
		"bar": {},
		"baz": {},
	}

	if !assert.Equal(t, expected, fields, "expected fields match") {
		return
	}
}

func TestLookupSructFieldFromJSONName(t *testing.T) {
	rv := reflect.ValueOf(X{})
	zero := reflect.Value{}

	data := map[string]string{
		"foo": "Foo",
		"bar": "Bar",
		"baz": "Baz",
	}

	for jsname, fname := range data {
		i := structinfo.StructFieldFromJSONName(rv, jsname)
		if !assert.NotEqual(t, i, -1, "should find '%s'", jsname) {
			return
		}

		sf := rv.Type().Field(i)
		if !assert.NotEqual(t, zero, sf, "should not be a zero value") {
			return
		}

		if !assert.Equal(t, sf.Name, fname, "'%s' should map to '%s'", jsname, fname) {
			return
		}
	}
}
