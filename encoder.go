package thrift

import (
	"reflect"
	"runtime"
)

type Encoder interface {
	EncodeThrift(protocol Protocol) error
}

type encoder struct {
	Protocol Protocol
}

func EncodeStruct(protocol Protocol, v interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			}
			err = r.(error)
		}
	}()
	e := &encoder{protocol}
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
	if err := e.Protocol.WriteStructBegin(v.Type().Name()); err != nil {
		e.error(err)
	}
	for _, ef := range encodeFields(v.Type()) {
		if ef.id == 0 {
			continue
		}

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

		if err := e.Protocol.WriteFieldBegin(structField.Name, fieldType(fieldValue.Type()), int16(ef.id)); err != nil {
			e.error(err)
		}
		e.writeValue(fieldValue)
		if err := e.Protocol.WriteFieldEnd(); err != nil {
			e.error(err)
		}
	}
	e.Protocol.WriteFieldStop()
	if err := e.Protocol.WriteStructEnd(); err != nil {
		e.error(err)
	}
}

func (e *encoder) writeValue(v reflect.Value) {
	var err error = nil
	switch v.Kind() {
	case reflect.Bool:
		err = e.Protocol.WriteBool(v.Bool())
	case reflect.Int8:
		err = e.Protocol.WriteByte(byte(v.Int()))
	case reflect.Int16:
		err = e.Protocol.WriteI16(int16(v.Int()))
	case reflect.Int32, reflect.Int:
		err = e.Protocol.WriteI32(int32(v.Int()))
	case reflect.Int64:
		err = e.Protocol.WriteI64(int64(v.Int()))
	case reflect.Float64:
		err = e.Protocol.WriteDouble(v.Float())
	case reflect.String:
		err = e.Protocol.WriteString(v.String())
	case reflect.Struct:
		e.writeStruct(v)
	case reflect.Map:
		keyType := v.Type().Key()
		valueType := v.Type().Elem()
		if er := e.Protocol.WriteMapBegin(fieldType(keyType), fieldType(valueType), v.Len()); er != nil {
			e.error(er)
		}
		for _, k := range v.MapKeys() {
			e.writeValue(k)
			e.writeValue(v.MapIndex(k))
		}
		err = e.Protocol.WriteMapEnd()
	case reflect.Slice, reflect.Array:
		elemType := v.Type().Elem()
		if er := e.Protocol.WriteListBegin(fieldType(elemType), v.Len()); er != nil {
			e.error(er)
		}
		n := v.Len()
		for i := 0; i < n; i++ {
			e.writeValue(v.Index(i))
		}
		err = e.Protocol.WriteListEnd()
	case reflect.Ptr, reflect.Interface:
		e.writeValue(v.Elem())
	default:
		e.error(&UnsupportedTypeError{v.Type()})
	}

	if err != nil {
		e.error(err)
	}

	return
}
