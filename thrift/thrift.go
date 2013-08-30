package thrift

import (
	"errors"
	"fmt"
	"io"
	"reflect"
	"sync"
)

const (
	TypeStop   = 0
	TypeVoid   = 1
	TypeBool   = 2
	TypeByte   = 3
	TypeI08    = 3
	TypeDouble = 4
	TypeI16    = 6
	TypeI32    = 8
	TypeI64    = 10
	TypeString = 11
	TypeUtf7   = 11
	TypeStruct = 12
	TypeMap    = 13
	TypeSet    = 14
	TypeList   = 15
	TypeUtf8   = 16
	TypeUtf16  = 17
)

const (
	MessageTypeCall      = 1
	MessageTypeReply     = 2
	MessageTypeException = 3
	MessageTypeOneway    = 4
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
		return TypeBool
	case reflect.Int8, reflect.Uint8:
		return TypeByte
	case reflect.Int16:
		return TypeI16
	case reflect.Int32, reflect.Uint32, reflect.Int:
		return TypeI32
	case reflect.Int64, reflect.Uint64:
		return TypeI64
	case reflect.Float64:
		return TypeDouble
	case reflect.Map:
		valueType := t.Elem()
		if valueType.Kind() == reflect.Struct && valueType.Name() == "" && valueType.NumField() == 0 {
			return TypeSet
		}
		return TypeMap
	case reflect.Slice:
		elemType := t.Elem()
		if elemType.Kind() == reflect.Uint8 {
			return TypeString
		} else {
			return TypeList
		}
	case reflect.Struct:
		return TypeStruct
	case reflect.String:
		return TypeString
	case reflect.Interface, reflect.Ptr:
		return fieldType(t.Elem())
	}
	panic(&UnsupportedTypeError{t})
}

func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice:
		return v.IsNil()
	case reflect.String:
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
				ef.fieldType = TypeSet
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
	case TypeBool:
		_, err = p.ReadBool(r)
	case TypeByte:
		_, err = p.ReadByte(r)
	case TypeI16:
		_, err = p.ReadI16(r)
	case TypeI32:
		_, err = p.ReadI32(r)
	case TypeI64:
		_, err = p.ReadI64(r)
	case TypeDouble:
		_, err = p.ReadDouble(r)
	case TypeString:
		_, err = p.ReadBytes(r)
	case TypeStruct:
		if err := p.ReadStructBegin(r); err != nil {
			return err
		}
		for {
			ftype, _, err := p.ReadFieldBegin(r)
			if err != nil {
				return err
			}
			if ftype == TypeStop {
				break
			}
			if err = SkipValue(r, p, ftype); err != nil {
				return err
			}
			if err = p.ReadFieldEnd(r); err != nil {
				return err
			}
		}
		return p.ReadStructEnd(r)
	case TypeMap:
		keyType, valueType, n, err := p.ReadMapBegin(r)
		if err != nil {
			return err
		}

		for i := 0; i < n; i++ {
			if err = SkipValue(r, p, keyType); err != nil {
				return err
			}
			if err = SkipValue(r, p, valueType); err != nil {
				return err
			}
		}

		return p.ReadMapEnd(r)
	case TypeList:
		valueType, n, err := p.ReadListBegin(r)
		if err != nil {
			return err
		}
		for i := 0; i < n; i++ {
			if err = SkipValue(r, p, valueType); err != nil {
				return err
			}
		}
		return p.ReadListEnd(r)
	case TypeSet:
		valueType, n, err := p.ReadSetBegin(r)
		if err != nil {
			return err
		}
		for i := 0; i < n; i++ {
			if err = SkipValue(r, p, valueType); err != nil {
				return err
			}
		}
		return p.ReadSetEnd(r)
	}
	return err
}

func ReadValue(r io.Reader, p Protocol, thriftType byte) (interface{}, error) {
	switch thriftType {
	case TypeBool:
		return p.ReadBool(r)
	case TypeByte:
		return p.ReadByte(r)
	case TypeI16:
		return p.ReadI16(r)
	case TypeI32:
		return p.ReadI32(r)
	case TypeI64:
		return p.ReadI64(r)
	case TypeDouble:
		return p.ReadDouble(r)
	case TypeString:
		return p.ReadString(r)
	case TypeStruct:
		if err := p.ReadStructBegin(r); err != nil {
			return nil, err
		}
		st := make(map[int]interface{})
		for {
			ftype, id, err := p.ReadFieldBegin(r)
			if err != nil {
				return st, err
			}
			if ftype == TypeStop {
				break
			}
			v, err := ReadValue(r, p, ftype)
			if err != nil {
				return st, err
			}
			st[int(id)] = v
			if err = p.ReadFieldEnd(r); err != nil {
				return st, err
			}
		}
		return st, p.ReadStructEnd(r)
	case TypeMap:
		keyType, valueType, n, err := p.ReadMapBegin(r)
		if err != nil {
			return nil, err
		}

		mp := make(map[interface{}]interface{})
		for i := 0; i < n; i++ {
			k, err := ReadValue(r, p, keyType)
			if err != nil {
				return mp, err
			}
			v, err := ReadValue(r, p, valueType)
			if err != nil {
				return mp, err
			}
			mp[k] = v
		}

		return mp, p.ReadMapEnd(r)
	case TypeList:
		valueType, n, err := p.ReadListBegin(r)
		if err != nil {
			return nil, err
		}
		lst := make([]interface{}, 0)
		for i := 0; i < n; i++ {
			v, err := ReadValue(r, p, valueType)
			if err != nil {
				return lst, err
			}
			lst = append(lst, v)
		}
		return lst, p.ReadListEnd(r)
	case TypeSet:
		valueType, n, err := p.ReadSetBegin(r)
		if err != nil {
			return nil, err
		}
		set := make([]interface{}, 0)
		for i := 0; i < n; i++ {
			v, err := ReadValue(r, p, valueType)
			if err != nil {
				return set, err
			}
			set = append(set, v)
		}
		return set, p.ReadSetEnd(r)
	}
	return nil, errors.New("thrift: unknown type")
}
