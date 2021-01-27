package structinfo_test

import (
	"fmt"
	"log"
	"reflect"
	"testing"

	"github.com/lestrrat-go/structinfo"
	"github.com/stretchr/testify/assert"
)

type Quux struct {
	Baz string `json:"baz"`
}

type X struct {
	private int
	Quux
	Foo string `json:"foo"`
	Bar string `json:"bar,omitempty"`
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

	data := map[string]string{
		"foo": "Foo",
		"bar": "Bar",
		"baz": "Baz",
	}

	for jsname, fname := range data {
		fn := structinfo.StructFieldFromJSONName(rv, jsname)
		if !assert.NotEqual(t, fn, "", "should find '%s'", jsname) {
			return
		}

		sf, ok := rv.Type().FieldByName(fn)
		if !assert.True(t, ok, "should be able resolve '%s' (%s)", jsname, fn) {
			return
		}

		if !assert.Equal(t, sf.Name, fname, "'%s' should map to '%s'", jsname, fname) {
			return
		}
	}
}

func TestStore_LookupFromJSONName(t *testing.T) {
	x := &X{
		Foo: "abc",
		Bar: "def",
		Quux: Quux{
			Baz: "ghi",
		},
	}

	rv := reflect.ValueOf(x)
	data := map[string]string{
		"foo": "abc",
		"bar": "def",
		"baz": "ghi",
	}

	store := structinfo.NewStore()

	for jsname, value := range data {
		fv, err := store.FieldValue(rv, jsname)
		if !assert.NoError(t, err, `FieldValue should succeed`) {
			return
		}

		if !assert.Equal(t, fv.Interface(), value, `values should match`) {
			return
		}

		fv.Set(reflect.ValueOf("hacked"))
	}

	if !assert.Equal(t, x.Foo, "hacked", "x.Foo should be hacked") {
		return
	}
	if !assert.Equal(t, x.Bar, "hacked", "x.Bar should be hacked") {
		return
	}
	if !assert.Equal(t, x.Baz, "hacked", "x.Baz should be hacked") {
		return
	}
}

func ExampleStore_FieldValue() {
	var x struct {
		Name string `json:"name"`
	}

	nameVal, err := structinfo.DefaultStore.FieldValue(reflect.ValueOf(&x), "name")
	if err != nil {
		log.Fatal(err)
	}

	nameVal.SetString("foo")

	fmt.Println(x.Name)

	// Output:
	// foo
}
