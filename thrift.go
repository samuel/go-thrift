package thrift

import (
	"fmt"
	"io"
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

const (
	messageTypeCall      = 1
	messageTypeReply     = 2
	messageTypeException = 3
	messageTypeOneway    = 4
)

const (
	ExceptionUnknown            = 0
	ExceptionUnknownMethod      = 1
	ExceptionInvalidMessageType = 2
	ExceptionWrongMethodName    = 3
	ExceptionBadSequenceId      = 4
	ExceptionMissingResult      = 5
	ExceptionInternalError      = 6
	ExceptionProtocolError      = 7
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
	return fmt.Sprintf("thrift: unsupported value (%+v): %s", e.Value, e.Str)
}

// Application level thrift exception
type ApplicationException struct {
	Message string `thrift:"1"`
	Type    int32  `thrift:"2"`
}

func (e *ApplicationException) String() string {
	typeStr := "Unknown Exception"
	switch e.Type {
	case ExceptionUnknownMethod:
		typeStr = "Unknown Method"
	case ExceptionInvalidMessageType:
		typeStr = "Invalid Message Type"
	case ExceptionWrongMethodName:
		typeStr = "Wrong Method Name"
	case ExceptionBadSequenceId:
		typeStr = "Bad Sequence Id"
	case ExceptionMissingResult:
		typeStr = "Missing Result"
	case ExceptionInternalError:
		typeStr = "Internal Error"
	case ExceptionProtocolError:
		typeStr = "Protocol Error"
	}
	return fmt.Sprintf("%s: %s", typeStr, e.Message)
}

func fieldType(t reflect.Type) byte {
	switch t.Kind() {
	case reflect.Bool:
		return typeBool
	case reflect.Int8:
		return typeByte
	case reflect.Int16:
		return typeI16
	case reflect.Int32, reflect.Int:
		return typeI32
	case reflect.Int64:
		return typeI64
	case reflect.Map:
		return typeMap
	case reflect.Slice:
		elemType := t.Elem()
		if elemType.Kind() == reflect.Uint8 && elemType.Name() == "byte" {
			return typeString
		} else {
			return typeList
		}
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
	fieldType byte
	name      string
}

type structMeta struct {
	required uint64 // bitmap of required fields
	fields   map[int]encodeField
}

var (
	typeCacheLock     sync.RWMutex
	encodeFieldsCache = make(map[reflect.Type]structMeta)
)

// encodeFields returns a slice of encodeField for a given
// struct type.
func encodeFields(t reflect.Type) structMeta {
	typeCacheLock.RLock()
	m, ok := encodeFieldsCache[t]
	typeCacheLock.RUnlock()
	if ok {
		return m
	}

	typeCacheLock.Lock()
	defer typeCacheLock.Unlock()
	m, ok = encodeFieldsCache[t]
	if ok {
		return m
	}

	fs := make(map[int]encodeField)
	m = structMeta{fields: fs}
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
		tv := f.Tag.Get("thrift")
		if tv != "" {
			var ef encodeField
			ef.i = i
			ef.id = 0

			if tv == "-" {
				continue
			}
			id, opts := parseTag(tv)
			if id >= 64 {
				// TODO: figure out a better way to deal with this
				panic("thrift: field id must be < 64")
			}
			ef.id = id
			ef.name = f.Name
			ef.required = opts.Contains("required")
			if ef.required {
				m.required |= 1 << byte(id)
			}
			ef.keepEmpty = opts.Contains("keepempty")
			if opts.Contains("set") {
				ef.fieldType = typeSet
			} else {
				ef.fieldType = fieldType(f.Type)
			}

			fs[ef.id] = ef
		}
	}
	encodeFieldsCache[t] = m
	return m
}

func SkipValue(r io.Reader, p Protocol, thriftType byte) error {
	var err error
	switch thriftType {
	case typeBool:
		_, err = p.ReadBool(r)
	case typeByte:
		_, err = p.ReadByte(r)
	case typeI16:
		_, err = p.ReadI16(r)
	case typeI32:
		_, err = p.ReadI32(r)
	case typeI64:
		_, err = p.ReadI64(r)
	case typeDouble:
		_, err = p.ReadDouble(r)
	case typeString:
		_, err = p.ReadBytes(r)
	case typeStruct:
		if err := p.ReadStructBegin(r); err != nil {
			return err
		}
		for {
			ftype, _, err := p.ReadFieldBegin(r)
			if err != nil {
				return err
			}
			if ftype == typeStop {
				break
			}
			SkipValue(r, p, ftype)
			if err = p.ReadFieldEnd(r); err != nil {
				return err
			}
		}
		return p.ReadStructEnd(r)
	case typeMap:
		keyType, valueType, n, err := p.ReadMapBegin(r)
		if err != nil {
			return err
		}

		for i := 0; i < n; i++ {
			SkipValue(r, p, keyType)
			SkipValue(r, p, valueType)
		}

		return p.ReadMapEnd(r)
	case typeList:
		valueType, n, err := p.ReadListBegin(r)
		if err != nil {
			return err
		}
		for i := 0; i < n; i++ {
			SkipValue(r, p, valueType)
		}
		return p.ReadListEnd(r)
	case typeSet:
		valueType, n, err := p.ReadSetBegin(r)
		if err != nil {
			return err
		}
		for i := 0; i < n; i++ {
			SkipValue(r, p, valueType)
		}
		return p.ReadSetEnd(r)
	}
	return err
}
