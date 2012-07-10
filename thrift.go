package thrift

import (
	"reflect"
	"sync"
)

const (
	typeStop   = 0
	typeVoid   = 1
	typeBool   = 2
	typeByte   = 3
	typeI08    = 3
	typeDouble = 4
	typeI16    = 6
	typeI32    = 8
	typeI64    = 10
	typeString = 11
	typeUtf7   = 11
	typeStruct = 12
	typeMap    = 13
	typeSet    = 14
	typeList   = 15
	typeUtf8   = 16
	typeUtf16  = 17
)

var (
	versionMask uint32 = 0xffff0000
	version1    uint32 = 0x80010000
	typeMask    uint32 = 0x000000ff
)

type MissingRequiredField struct {
	StructName string
	FieldName  string
}

func (e *MissingRequiredField) Error() string {
	return "thrift: missing required field: " + e.StructName + "." + e.FieldName
}

type UnsupportedTypeError struct {
	Type reflect.Type
}

func (e *UnsupportedTypeError) Error() string {
	return "thrift: unsupported type: " + e.Type.String()
}

type UnsupportedValueError struct {
	Value reflect.Value
	Str   string
}

func (e *UnsupportedValueError) Error() string {
	return "thrift: unsupported value: " + e.Str
}

func fieldType(t reflect.Type) byte {
	switch t.Kind() {
	case reflect.Bool:
		return typeBool
	case reflect.Int8:
		return typeBool
	case reflect.Int16:
		return typeI16
	case reflect.Int32, reflect.Int:
		return typeI32
	case reflect.Int64:
		return typeI64
	case reflect.Map:
		return typeMap
	case reflect.Slice:
		return typeList
	case reflect.Struct:
		return typeStruct
	case reflect.String:
		return typeString
	case reflect.Interface, reflect.Ptr:
		return fieldType(t.Elem())
	}
	panic(&UnsupportedTypeError{t})
}

func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}

// encodeField contains information about how to encode a field of a
// struct.
type encodeField struct {
	i         int // field index in struct
	id        int
	required  bool
	keepEmpty bool
}

var (
	typeCacheLock     sync.RWMutex
	encodeFieldsCache = make(map[reflect.Type]map[int]encodeField)
)

// encodeFields returns a slice of encodeField for a given
// struct type.
func encodeFields(t reflect.Type) map[int]encodeField {
	typeCacheLock.RLock()
	fs, ok := encodeFieldsCache[t]
	typeCacheLock.RUnlock()
	if ok {
		return fs
	}

	typeCacheLock.Lock()
	defer typeCacheLock.Unlock()
	fs, ok = encodeFieldsCache[t]
	if ok {
		return fs
	}

	fs = make(map[int]encodeField)
	v := reflect.Zero(t)
	n := v.NumField()
	for i := 0; i < n; i++ {
		f := t.Field(i)
		if f.PkgPath != "" {
			continue
		}
		if f.Anonymous {
			// We want to do a better job with these later,
			// so for now pretend they don't exist.
			continue
		}
		var ef encodeField
		ef.i = i
		ef.id = 0

		tv := f.Tag.Get("thrift")
		if tv != "" {
			if tv == "-" {
				continue
			}
			id, opts := parseTag(tv)
			ef.id = id
			ef.required = opts.Contains("required")
			ef.keepEmpty = opts.Contains("keepempty")
		}
		fs[ef.id] = ef
	}
	encodeFieldsCache[t] = fs
	return fs
}
