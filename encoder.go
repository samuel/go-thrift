package thrift

import (
	"io"
	"reflect"
	"runtime"
)

type Encoder interface {
	EncodeThrift(io.Writer, Protocol) error
}

type encoder struct {
	w io.Writer
	p Protocol
}

func EncodeStruct(w io.Writer, protocol Protocol, v interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			}
			err = r.(error)
		}
	}()
	e := &encoder{w, protocol}
	vo := reflect.ValueOf(v)
	e.writeStruct(vo)
	return nil
}

func (e *encoder) error(err interface{}) {
	panic(err)
}

func (e *encoder) writeStruct(v reflect.Value) {
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		e.error(&UnsupportedValueError{Value: v, Str: "expected a struct"})
	}
	if err := e.p.WriteStructBegin(e.w, v.Type().Name()); err != nil {
		e.error(err)
	}
	for _, ef := range encodeFields(v.Type()).fields {
		structField := v.Type().Field(ef.i)
		fieldValue := v.Field(ef.i)

		if !ef.required && !ef.keepEmpty && isEmptyValue(fieldValue) {
			continue
		}

		if fieldValue.Kind() == reflect.Ptr {
			if ef.required && fieldValue.IsNil() {
				e.error(&MissingRequiredField{v.Type().Name(), structField.Name})
			}
		}

		ftype := ef.fieldType

		if err := e.p.WriteFieldBegin(e.w, structField.Name, ftype, int16(ef.id)); err != nil {
			e.error(err)
		}
		e.writeValue(fieldValue, ftype)
		if err := e.p.WriteFieldEnd(e.w); err != nil {
			e.error(err)
		}
	}
	e.p.WriteFieldStop(e.w)
	if err := e.p.WriteStructEnd(e.w); err != nil {
		e.error(err)
	}
}

func (e *encoder) writeValue(v reflect.Value, thriftType byte) {
	kind := v.Kind()
	if kind == reflect.Ptr || kind == reflect.Interface {
		v = v.Elem()
		kind = v.Kind()
	}

	var err error = nil
	switch thriftType {
	case typeBool:
		err = e.p.WriteBool(e.w, v.Bool())
	case typeByte:
		err = e.p.WriteByte(e.w, byte(v.Int()))
	case typeI16:
		err = e.p.WriteI16(e.w, int16(v.Int()))
	case typeI32:
		err = e.p.WriteI32(e.w, int32(v.Int()))
	case typeI64:
		err = e.p.WriteI64(e.w, int64(v.Int()))
	case typeDouble:
		err = e.p.WriteDouble(e.w, v.Float())
	case typeString:
		if kind == reflect.Slice {
			elemType := v.Type().Elem()
			if elemType.Kind() == reflect.Uint8 && elemType.Name() == "byte" {
				err = e.p.WriteBytes(e.w, v.Bytes())
			} else {
				err = &UnsupportedValueError{Value: v, Str: "expected a byte array"}
			}
		} else {
			err = e.p.WriteString(e.w, v.String())
		}
	case typeStruct:
		e.writeStruct(v)
	case typeMap:
		keyType := v.Type().Key()
		valueType := v.Type().Elem()
		keyThriftType := fieldType(keyType)
		valueThriftType := fieldType(valueType)
		if er := e.p.WriteMapBegin(e.w, keyThriftType, valueThriftType, v.Len()); er != nil {
			e.error(er)
		}
		for _, k := range v.MapKeys() {
			e.writeValue(k, fieldType(keyType))
			e.writeValue(v.MapIndex(k), fieldType(valueType))
		}
		err = e.p.WriteMapEnd(e.w)
	case typeList:
		elemType := v.Type().Elem()
		if elemType.Kind() == reflect.Uint8 && elemType.Name() == "byte" {
			err = e.p.WriteBytes(e.w, v.Bytes())
		} else {
			elemThriftType := fieldType(elemType)
			if er := e.p.WriteListBegin(e.w, elemThriftType, v.Len()); er != nil {
				e.error(er)
			}
			n := v.Len()
			for i := 0; i < n; i++ {
				e.writeValue(v.Index(i), elemThriftType)
			}
			err = e.p.WriteListEnd(e.w)
		}
	case typeSet:
		elemType := v.Type().Elem()
		elemThriftType := fieldType(elemType)
		if er := e.p.WriteSetBegin(e.w, elemThriftType, v.Len()); er != nil {
			e.error(er)
		}
		n := v.Len()
		for i := 0; i < n; i++ {
			e.writeValue(v.Index(i), elemThriftType)
		}
		err = e.p.WriteSetEnd(e.w)
	default:
		e.error(&UnsupportedTypeError{v.Type()})
	}

	if err != nil {
		e.error(err)
	}
}
