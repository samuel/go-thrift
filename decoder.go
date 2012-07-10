package thrift

import (
	"reflect"
	// "runtime"
)

type Decoder struct {
	Protocol Protocol
}

func (d *Decoder) error(err interface{}) {
	panic(err)
}

func (d *Decoder) ReadStruct(v interface{}) (err error) {
	// defer func() {
	// 	if r := recover(); r != nil {
	// 		if _, ok := r.(runtime.Error); ok {
	// 			panic(r)
	// 		}
	// 		err = r.(error)
	// 	}
	// }()
	vo := reflect.ValueOf(v)
	for vo.Kind() != reflect.Ptr {
		d.error(&UnsupportedValueError{Value: vo, Str: "pointer to struct expected"})
	}
	if vo.Elem().Kind() != reflect.Struct {
		d.error(&UnsupportedValueError{Value: vo, Str: "expected a struct"})
	}
	d.readValue(typeStruct, vo.Elem())
	return nil
}

func (d *Decoder) readValue(ftype byte, rf reflect.Value) {
	v := rf
	if rf.Kind() == reflect.Ptr {
		if rf.IsNil() {
			rf.Set(reflect.New(rf.Type().Elem()))
		}
		v = rf.Elem()
	}

	if ftype != fieldType(v.Type()) {
		d.error(&UnsupportedValueError{Value: v, Str: "type mistmatch"})
	}

	var err error = nil
	switch v.Kind() {
	case reflect.Bool:
		if val, err := d.Protocol.ReadBool(); err != nil {
			d.error(err)
		} else {
			v.SetBool(val)
		}
	case reflect.Int8:
		if val, err := d.Protocol.ReadByte(); err != nil {
			d.error(err)
		} else {
			v.SetInt(int64(val))
		}
	case reflect.Int16:
		if val, err := d.Protocol.ReadI16(); err != nil {
			d.error(err)
		} else {
			v.SetInt(int64(val))
		}
	case reflect.Int32, reflect.Int:
		if val, err := d.Protocol.ReadI32(); err != nil {
			d.error(err)
		} else {
			v.SetInt(int64(val))
		}
	case reflect.Int64:
		if val, err := d.Protocol.ReadI64(); err != nil {
			d.error(err)
		} else {
			v.SetInt(val)
		}
	case reflect.Float64:
		if val, err := d.Protocol.ReadDouble(); err != nil {
			d.error(err)
		} else {
			v.SetFloat(val)
		}
	case reflect.String:
		if val, err := d.Protocol.ReadString(); err != nil {
			d.error(err)
		} else {
			v.SetString(val)
		}
	case reflect.Struct:
		if err := d.Protocol.ReadStructBegin(); err != nil {
			d.error(err)
		}

		fields := encodeFields(v.Type())
		for {
			ftype, id, err := d.Protocol.ReadFieldBegin()
			if err != nil {
				d.error(err)
			}
			if ftype == typeStop {
				break
			}

			ef, ok := fields[int(id)]
			if !ok {
				// Ignore unknown fields
				// TODO
				d.error(&UnsupportedValueError{Str: "TODO"})
			} else {
				fieldValue := v.Field(ef.i)
				d.readValue(ftype, fieldValue)
			}

			if err = d.Protocol.ReadFieldEnd(); err != nil {
				d.error(err)
			}
		}

		if err := d.Protocol.ReadStructEnd(); err != nil {
			d.error(err)
		}
	case reflect.Map:
		keyType := v.Type().Key()
		valueType := v.Type().Elem()
		ktype, vtype, n, err := d.Protocol.ReadMapBegin()
		if err != nil {
			d.error(err)
		}
		v.Set(reflect.MakeMap(v.Type()))
		for i := 0; i < n; i++ {
			key := reflect.New(keyType).Elem()
			val := reflect.New(valueType).Elem()
			d.readValue(ktype, key)
			d.readValue(vtype, val)
			v.SetMapIndex(key, val)
		}
		if err := d.Protocol.ReadMapEnd(); err != nil {
			d.error(err)
		}
	case reflect.Slice, reflect.Array:
		elemType := v.Type().Elem()
		et, n, err := d.Protocol.ReadListBegin()
		if err != nil {
			d.error(err)
		}
		for i := 0; i < n; i++ {
			val := reflect.New(elemType)
			d.readValue(et, val.Elem())
			v.Set(reflect.Append(v, val.Elem()))
		}
		if err := d.Protocol.ReadListEnd(); err != nil {
			d.error(err)
		}
	default:
		d.error(&UnsupportedTypeError{v.Type()})
	}

	if err != nil {
		d.error(err)
	}

	return
}

func (d *Decoder) readSimpleValue(fieldType int) (val interface{}) {
	var err error = nil
	switch fieldType {
	case typeBool:
		val, err = d.Protocol.ReadBool()
	case typeByte:
		val, err = d.Protocol.ReadByte()
	case typeI16:
		val, err = d.Protocol.ReadI16()
	case typeI32:
		val, err = d.Protocol.ReadI32()
	case typeI64:
		val, err = d.Protocol.ReadI64()
	case typeDouble:
		val, err = d.Protocol.ReadDouble()
	case typeString:
		val, err = d.Protocol.ReadString()
	default:
		d.error(&UnsupportedTypeError{})
	}

	if err != nil {
		d.error(err)
	}

	return
}
