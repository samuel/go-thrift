package thrift

import (
	"bytes"
	"testing"
	"testing/quick"
)

var builders = []ProtocolBuilder{BinaryProtocol}

func TestMessageBegin(t *testing.T) {
	buf := new(bytes.Buffer)

	for _, pb := range builders {
		if err := quick.Check(func(name string, messageType byte, seqid int32) bool {
			buf.Reset()

			pw := pb.NewProtocolWriter(buf)
			pr := pb.NewProtocolReader(buf)

			if err := pw.WriteMessageBegin(name, messageType, seqid); err != nil {
				t.Error(err)
				return false
			}

			name2, messageType2, seqid2, err := pr.ReadMessageBegin()
			if err != nil {
				t.Error(err)
				return false
			}

			if name2 != name || messageType2 != messageType || seqid2 != seqid {
				return false
			}

			return true
		}, nil); err != nil {
			t.Error(err)
		}
	}
}

// MessageEnd()
func TestMessageEnd(t *testing.T) {
	buf := new(bytes.Buffer)

	for _, pb := range builders {
		if err := quick.Check(func() bool {
			buf.Reset()

			pw := pb.NewProtocolWriter(buf)
			pr := pb.NewProtocolReader(buf)

			if err := pw.WriteMessageEnd(); err != nil {
				t.Error(err)
				return false
			}

			err := pr.ReadMessageEnd()
			if err != nil {
				t.Error(err)
				return false
			}
			return true
		}, nil); err != nil {
			t.Error(err)
		}
	}
}

// StructBegin(name string)
func TestStructBegin(t *testing.T) {
	buf := new(bytes.Buffer)

	for _, pb := range builders {
		if err := quick.Check(func(name string) bool {
			buf.Reset()

			pw := pb.NewProtocolWriter(buf)
			pr := pb.NewProtocolReader(buf)

			if err := pw.WriteStructBegin(name); err != nil {
				t.Error(err)
				return false
			}

			err := pr.ReadStructBegin()
			if err != nil {
				t.Error(err)
				return false
			}

			return true
		}, nil); err != nil {
			t.Error(err)
		}
	}
}

// StructEnd()
func TestStructEnd(t *testing.T) {
	buf := new(bytes.Buffer)

	for _, pb := range builders {
		if err := quick.Check(func() bool {
			buf.Reset()

			pw := pb.NewProtocolWriter(buf)
			pr := pb.NewProtocolReader(buf)

			if err := pw.WriteStructEnd(); err != nil {
				t.Error(err)
				return false
			}

			err := pr.ReadStructEnd()
			if err != nil {
				t.Error(err)
				return false
			}
			return true
		}, nil); err != nil {
			t.Error(err)
		}
	}
}

// FieldBegin(name string, fieldType byte, id int16)
func TestFieldBegin(t *testing.T) {
	buf := new(bytes.Buffer)

	for _, pb := range builders {
		if err := quick.Check(func(name string, fieldType byte, id int16) bool {
			buf.Reset()

			pw := pb.NewProtocolWriter(buf)
			pr := pb.NewProtocolReader(buf)

			if err := pw.WriteFieldBegin(name, fieldType, id); err != nil {
				t.Error(err)
				return false
			}

			fieldType2, id2, err := pr.ReadFieldBegin()
			if err != nil {
				t.Error(err)
				return false
			}

			if fieldType2 != fieldType || (fieldType != TypeStop && id2 != id) {
				t.Logf("%d != %d || %d != %d", fieldType2, fieldType, id2, id)
				return false
			}

			return true
		}, nil); err != nil {
			t.Error(err)
		}
	}
}

// FieldEnd()
func TestFieldEnd(t *testing.T) {
	buf := new(bytes.Buffer)

	for _, pb := range builders {
		if err := quick.Check(func() bool {
			buf.Reset()

			pw := pb.NewProtocolWriter(buf)
			pr := pb.NewProtocolReader(buf)

			if err := pw.WriteFieldEnd(); err != nil {
				t.Error(err)
				return false
			}

			err := pr.ReadFieldEnd()
			if err != nil {
				t.Error(err)
				return false
			}
			return true
		}, nil); err != nil {
			t.Error(err)
		}
	}
}

// MapBegin(keyType byte, valueType byte, size int)
func TestMapBegin(t *testing.T) {
	buf := new(bytes.Buffer)

	for _, pb := range builders {
		if err := quick.Check(func(keyType byte, valueType byte, size int32) bool {
			buf.Reset()

			pw := pb.NewProtocolWriter(buf)
			pr := pb.NewProtocolReader(buf)

			if err := pw.WriteMapBegin(keyType, valueType, int(size)); err != nil {
				t.Error(err)
				return false
			}

			keyType2, valueType2, size2, err := pr.ReadMapBegin()
			if err != nil {
				t.Error(err)
				return false
			}

			if keyType2 != keyType || valueType2 != valueType || size2 != int(size) {
				t.Logf("%d != %d || %d != %d || %d != %d", keyType2, keyType, valueType2, valueType, size2, size)
				return false
			}

			return true
		}, nil); err != nil {
			t.Error(err)
		}
	}
}

// MapEnd()
func TestMapEnd(t *testing.T) {
	buf := new(bytes.Buffer)

	for _, pb := range builders {
		if err := quick.Check(func() bool {
			buf.Reset()

			pw := pb.NewProtocolWriter(buf)
			pr := pb.NewProtocolReader(buf)

			if err := pw.WriteMapEnd(); err != nil {
				t.Error(err)
				return false
			}

			err := pr.ReadMapEnd()
			if err != nil {
				t.Error(err)
				return false
			}

			return true
		}, nil); err != nil {
			t.Error(err)
		}
	}
}

// ListBegin(elementType byte, size int)
func TestListBegin(t *testing.T) {
	buf := new(bytes.Buffer)

	for _, pb := range builders {
		if err := quick.Check(func(elementType byte, size int32) bool {
			buf.Reset()

			pw := pb.NewProtocolWriter(buf)
			pr := pb.NewProtocolReader(buf)

			if err := pw.WriteListBegin(elementType, int(size)); err != nil {
				t.Error(err)
				return false
			}

			elementType2, size2, err := pr.ReadListBegin()
			if err != nil {
				t.Error(err)
				return false
			}

			if elementType2 != elementType || size2 != int(size) {
				t.Logf("%d != %d || %d != %d", elementType2, elementType, size2, size)
				return false
			}

			return true
		}, nil); err != nil {
			t.Error(err)
		}
	}
}

// ListEnd()
func TestListEnd(t *testing.T) {
	buf := new(bytes.Buffer)

	for _, pb := range builders {
		if err := quick.Check(func() bool {
			buf.Reset()

			pw := pb.NewProtocolWriter(buf)
			pr := pb.NewProtocolReader(buf)

			if err := pw.WriteListEnd(); err != nil {
				t.Error(err)
				return false
			}

			err := pr.ReadListEnd()
			if err != nil {
				t.Error(err)
				return false
			}

			return true
		}, nil); err != nil {
			t.Error(err)
		}
	}
}

// SetBegin(elementType byte, size int)
func TestSetBegin(t *testing.T) {
	buf := new(bytes.Buffer)

	for _, pb := range builders {
		if err := quick.Check(func(elementType byte, size int32) bool {
			buf.Reset()

			pw := pb.NewProtocolWriter(buf)
			pr := pb.NewProtocolReader(buf)

			if err := pw.WriteSetBegin(elementType, int(size)); err != nil {
				t.Error(err)
				return false
			}

			elementType2, size2, err := pr.ReadSetBegin()
			if err != nil {
				t.Error(err)
				return false
			}

			if elementType2 != elementType || size2 != int(size) {
				t.Logf("%d != %d || %d != %d", elementType2, elementType, size2, size)
				return false
			}

			return true
		}, nil); err != nil {
			t.Error(err)
		}
	}
}

// SetEnd()
func TestSetEnd(t *testing.T) {
	buf := new(bytes.Buffer)

	for _, pb := range builders {
		if err := quick.Check(func() bool {
			buf.Reset()

			pw := pb.NewProtocolWriter(buf)
			pr := pb.NewProtocolReader(buf)

			if err := pw.WriteSetEnd(); err != nil {
				t.Error(err)
				return false
			}

			err := pr.ReadSetEnd()
			if err != nil {
				t.Error(err)
				return false
			}
			return true
		}, nil); err != nil {
			t.Error(err)
		}
	}
}

// Bool(value bool)
func TestBool(t *testing.T) {
	buf := new(bytes.Buffer)

	for _, pb := range builders {
		if err := quick.Check(func(value bool) bool {
			buf.Reset()

			pw := pb.NewProtocolWriter(buf)
			pr := pb.NewProtocolReader(buf)

			if err := pw.WriteBool(value); err != nil {
				t.Error(err)
				return false
			}

			value2, err := pr.ReadBool()
			if err != nil {
				t.Error(err)
				return false
			}

			return value2 == value
		}, nil); err != nil {
			t.Error(err)
		}
	}
}

// Byte(value byte)
func TestByte(t *testing.T) {
	buf := new(bytes.Buffer)

	for _, pb := range builders {
		if err := quick.Check(func(value byte) bool {
			buf.Reset()

			pw := pb.NewProtocolWriter(buf)
			pr := pb.NewProtocolReader(buf)

			if err := pw.WriteByte(value); err != nil {
				t.Error(err)
				return false
			}

			value2, err := pr.ReadByte()
			if err != nil {
				t.Error(err)
				return false
			}
			return value2 == value
		}, nil); err != nil {
			t.Error(err)
		}
	}
}

// I16(value int16)
func TestI16(t *testing.T) {
	buf := new(bytes.Buffer)

	for _, pb := range builders {
		if err := quick.Check(func(value int16) bool {
			buf.Reset()

			pw := pb.NewProtocolWriter(buf)
			pr := pb.NewProtocolReader(buf)

			if err := pw.WriteI16(value); err != nil {
				t.Error(err)
				return false
			}

			value2, err := pr.ReadI16()
			if err != nil {
				t.Error(err)
				return false
			}
			return value2 == value
		}, nil); err != nil {
			t.Error(err)
		}
	}
}

// I32(value int32)
func TestI32(t *testing.T) {
	buf := new(bytes.Buffer)

	for _, pb := range builders {
		if err := quick.Check(func(value int32) bool {
			buf.Reset()

			pw := pb.NewProtocolWriter(buf)
			pr := pb.NewProtocolReader(buf)

			if err := pw.WriteI32(value); err != nil {
				t.Error(err)
				return false
			}

			value2, err := pr.ReadI32()
			if err != nil {
				t.Error(err)
				return false
			}
			return value2 == value
		}, nil); err != nil {
			t.Error(err)
		}
	}
}

// I64(value int64)
func TestI64(t *testing.T) {
	buf := new(bytes.Buffer)

	for _, pb := range builders {
		if err := quick.Check(func(value int64) bool {
			buf.Reset()

			pw := pb.NewProtocolWriter(buf)
			pr := pb.NewProtocolReader(buf)

			if err := pw.WriteI64(value); err != nil {
				t.Error(err)
				return false
			}

			value2, err := pr.ReadI64()
			if err != nil {
				t.Error(err)
				return false
			}
			return value2 == value
		}, nil); err != nil {
			t.Error(err)
		}
	}
}

// Double(value float64)
func TestDouble(t *testing.T) {
	buf := new(bytes.Buffer)

	for _, pb := range builders {
		if err := quick.Check(func(value float64) bool {
			buf.Reset()

			pw := pb.NewProtocolWriter(buf)
			pr := pb.NewProtocolReader(buf)

			if err := pw.WriteDouble(value); err != nil {
				t.Error(err)
				return false
			}

			value2, err := pr.ReadDouble()
			if err != nil {
				t.Error(err)
				return false
			}
			return value2 == value
		}, nil); err != nil {
			t.Error(err)
		}
	}
}

// String(value string)
func TestString(t *testing.T) {
	buf := new(bytes.Buffer)

	for _, pb := range builders {
		if err := quick.Check(func(value string) bool {
			buf.Reset()

			pw := pb.NewProtocolWriter(buf)
			pr := pb.NewProtocolReader(buf)

			if err := pw.WriteString(value); err != nil {
				t.Error(err)
				return false
			}

			value2, err := pr.ReadString()
			if err != nil {
				t.Error(err)
				return false
			}
			return value2 == value
		}, nil); err != nil {
			t.Error(err)
		}
	}
}

// Bytes(value []byte)
func TestBytes(t *testing.T) {
	buf := new(bytes.Buffer)

	for _, pb := range builders {
		if err := quick.Check(func(value []byte) bool {
			buf.Reset()

			pw := pb.NewProtocolWriter(buf)
			pr := pb.NewProtocolReader(buf)

			if err := pw.WriteBytes(value); err != nil {
				t.Error(err)
				return false
			}

			value2, err := pr.ReadBytes()
			if err != nil {
				t.Error(err)
				return false
			}
			return bytes.Equal(value2, value)
		}, nil); err != nil {
			t.Error(err)
		}
	}
}
