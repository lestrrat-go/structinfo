// Package structinfo contains tools to inspect structs.

package structinfo

import (
	"reflect"
	"sync"
)

type jsonFieldMap struct {
	lock   sync.Mutex
	fields map[string]int
}

var type2jfm = map[reflect.Type]jsonFieldMap{}
var type2jfmMutex = sync.Mutex{}

// StructFieldFromJSONName returns the struct field index on the
// given struct value. Value of -1 means the field is either not
// public, or does not exist.
//
// This can be used to map JSON field names to actual struct fields.
func StructFieldFromJSONName(v reflect.Value, name string) int {
	m := getType2jfm(v.Type())
	m.lock.Lock()
	defer m.lock.Unlock()

	i, ok := m.fields[name]
	if !ok {
		return -1
	}
	return i
}

func getType2jfm(t reflect.Type) jsonFieldMap {
	type2jfmMutex.Lock()
	defer type2jfmMutex.Unlock()

	fm, ok := type2jfm[t]
	if ok {
		return fm
	}

	fm = constructJfm(t)
	type2jfm[t] = fm
	return fm
}

func constructJfm(t reflect.Type) jsonFieldMap {
	fm := jsonFieldMap{
		fields: make(map[string]int),
	}
	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)
		if sf.PkgPath != "" { // unexported
			continue
		}

		tag := sf.Tag.Get("json")
		if tag == "" || tag == "-" || tag[0] == ',' {
			fm.fields[sf.Name] = i
			continue
		}

		flen := 0
		for j := 0; j < len(tag); j++ {
			if tag[j] == ',' {
				break
			}
			flen = j
		}
		fm.fields[tag[:flen+1]] = i
	}

	return fm
}