package thrift

import (
	"reflect"
	"runtime"
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

type Encoder struct {
	Protocol Protocol
}

// func Marshal(v interface{}) ([]byte, error) {
// 	e := &encodeState{}
// 	err := e.marshal(v)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return e.Bytes(), nil
// }

func (e *Encoder) WriteStruct(v interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			}
			err = r.(error)
		}
	}()
	e.reflectValue(reflect.ValueOf(v))
	return nil
}

func (e *Encoder) reflectValue(v reflect.Value) {
	if !v.IsValid() {
		// e.WriteString("null")
		return
	}

	switch v.Kind() {
	case reflect.Bool:
		if err := e.Protocol.WriteBool(v.Bool()); err != nil {
			panic(err)
		}
	case reflect.Int8:
		if err := e.Protocol.WriteByte(byte(v.Int())); err != nil {
			panic(err)
		}
	case reflect.Int16:
		if err := e.Protocol.WriteI16(int16(v.Int())); err != nil {
			panic(err)
		}
	case reflect.Int32, reflect.Int:
		if err := e.Protocol.WriteI32(int32(v.Int())); err != nil {
			panic(err)
		}
	case reflect.Int64:
		if err := e.Protocol.WriteI64(int64(v.Int())); err != nil {
			panic(err)
		}
	// case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
	// case reflect.Float32, reflect.Float64:
	case reflect.Float64:
		if err := e.Protocol.WriteDouble(v.Float()); err != nil {
			panic(err)
		}
	case reflect.String:
		if err := e.Protocol.WriteString(v.String()); err != nil {
			panic(err)
		}
	case reflect.Struct:
		if err := e.Protocol.WriteStructBegin("TODO"); err != nil {
			panic(err)
		}
		for _, ef := range encodeFields(v.Type()) {
			fieldValue := v.Field(ef.i)

			// TODO
			// if ef.omitEmpty && isEmptyValue(fieldValue) {
			// 	continue
			// }

			e.Protocol.WriteFieldBegin(fieldValue.Name, 1, int16(ef.id))
			e.reflectValue(fieldValue)
			e.Protocol.WriteFieldEnd()
		}
		if err := e.Protocol.WriteStructEnd(); err != nil {
			panic(err)
		}
	// case reflect.Map:
	// 	if v.Type().Key().Kind() != reflect.String {
	// 		e.error(&UnsupportedTypeError{v.Type()})
	// 	}
	// 	if v.IsNil() {
	// 		e.WriteString("null")
	// 		break
	// 	}
	// 	e.WriteByte('{')
	// 	var sv stringValues = v.MapKeys()
	// 	sort.Sort(sv)
	// 	for i, k := range sv {
	// 		if i > 0 {
	// 			e.WriteByte(',')
	// 		}
	// 		e.string(k.String())
	// 		e.WriteByte(':')
	// 		e.reflectValue(v.MapIndex(k))
	// 	}
	// 	e.WriteByte('}')

	// case reflect.Slice:
	// 	if v.IsNil() {
	// 		e.WriteString("null")
	// 		break
	// 	}
	// 	if v.Type().Elem().Kind() == reflect.Uint8 {
	// 		// Byte slices get special treatment; arrays don't.
	// 		s := v.Bytes()
	// 		e.WriteByte('"')
	// 		if len(s) < 1024 {
	// 			// for small buffers, using Encode directly is much faster.
	// 			dst := make([]byte, base64.StdEncoding.EncodedLen(len(s)))
	// 			base64.StdEncoding.Encode(dst, s)
	// 			e.Write(dst)
	// 		} else {
	// 			// for large buffers, avoid unnecessary extra temporary
	// 			// buffer space.
	// 			enc := base64.NewEncoder(base64.StdEncoding, e)
	// 			enc.Write(s)
	// 			enc.Close()
	// 		}
	// 		e.WriteByte('"')
	// 		break
	// 	}
	// 	// Slices can be marshalled as nil, but otherwise are handled
	// 	// as arrays.
	// 	fallthrough
	// case reflect.Array:
	// 	e.WriteByte('[')
	// 	n := v.Len()
	// 	for i := 0; i < n; i++ {
	// 		if i > 0 {
	// 			e.WriteByte(',')
	// 		}
	// 		e.reflectValue(v.Index(i))
	// 	}
	// 	e.WriteByte(']')

	// case reflect.Interface, reflect.Ptr:
	// 	if v.IsNil() {
	// 		e.WriteString("null")
	// 		return
	// 	}
	// 	e.reflectValue(v.Elem())

	default:
		e.error(&UnsupportedTypeError{v.Type()})
	}
	return
}

// encodeField contains information about how to encode a field of a
// struct.
type encodeField struct {
	i        int // field index in struct
	id       int
	required bool
}

var (
	typeCacheLock     sync.RWMutex
	encodeFieldsCache = make(map[reflect.Type][]encodeField)
)

// encodeFields returns a slice of encodeField for a given
// struct type.
func encodeFields(t reflect.Type) []encodeField {
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
		}
		fs = append(fs, ef)
	}
	encodeFieldsCache[t] = fs
	return fs
}
