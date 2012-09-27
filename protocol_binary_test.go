package thrift

import (
	"bytes"
	"testing"
)

type nullWriter int

func (n nullWriter) Write(b []byte) (int, error) {
	return len(b), nil
}

func TestBinaryProtocol(t *testing.T) {
	b := &bytes.Buffer{}
	pr := BinaryProtocol

	if err := pr.WriteBool(b, true); err != nil {
		t.Fatalf("write bool true failed: %+v", err)
	}
	if b, err := pr.ReadBool(b); err != nil {
		t.Fatalf("read bool true failed: %+v", err)
	} else if !b {
		t.Fatal("read bool true returned false")
	}

	if err := pr.WriteBool(b, false); err != nil {
		t.Fatalf("write bool false failed: %+v", err)
	}
	if b, err := pr.ReadBool(b); err != nil {
		t.Fatalf("read bool false failed: %+v", err)
	} else if b {
		t.Fatal("read bool false returned true")
	}

	if err := pr.WriteI16(b, 1234); err != nil {
		t.Fatalf("write i16 failed: %+v", err)
	}
	if v, err := pr.ReadI16(b); err != nil {
		t.Fatalf("read i16 failed: %+v", err)
	} else if v != 1234 {
		t.Fatal("read i16 returned %d expected 1234", v)
	}

	if err := pr.WriteI32(b, -1234); err != nil {
		t.Fatalf("write i32 failed: %+v", err)
	}
	if v, err := pr.ReadI32(b); err != nil {
		t.Fatalf("read i32 failed: %+v", err)
	} else if v != -1234 {
		t.Fatal("read i32 returned %d expected -1234", v)
	}

	if err := pr.WriteI64(b, -1234); err != nil {
		t.Fatalf("write i64 failed: %+v", err)
	}
	if v, err := pr.ReadI64(b); err != nil {
		t.Fatalf("read i64 failed: %+v", err)
	} else if v != -1234 {
		t.Fatal("read i64 returned %d expected -1234", v)
	}

	if err := pr.WriteDouble(b, -0.1234); err != nil {
		t.Fatalf("write double failed: %+v", err)
	}
	if v, err := pr.ReadDouble(b); err != nil {
		t.Fatalf("read double failed: %+v", err)
	} else if v != -0.1234 {
		t.Fatal("read double returned %.4f expected -0.1234", v)
	}

	testString := "012345"
	for i := 0; i < 2; i += 1 {
		if err := pr.WriteString(b, testString); err != nil {
			t.Fatalf("write string failed: %+v", err)
		}
		if v, err := pr.ReadString(b); err != nil {
			t.Fatalf("read string failed: %+v", err)
		} else if v != testString {
			t.Fatal("read string returned %s expected '%s'", v, testString)
		}
		testString += "012345"
	}
}

func BenchmarkBinaryProtocolReadByte(b *testing.B) {
	buf := bytes.NewBuffer(make([]byte, b.N))
	for i := 0; i < b.N; i++ {
		BinaryProtocol.ReadByte(buf)
	}
}

func BenchmarkBinaryProtocolReadI32(b *testing.B) {
	buf := bytes.NewBuffer(make([]byte, b.N*4))
	for i := 0; i < b.N; i++ {
		BinaryProtocol.ReadI32(buf)
	}
}

func BenchmarkBinaryProtocolWriteByte(b *testing.B) {
	buf := nullWriter(0)
	for i := 0; i < b.N; i++ {
		BinaryProtocol.WriteByte(buf, 1)
	}
}

func BenchmarkBinaryProtocolWriteI32(b *testing.B) {
	buf := nullWriter(0)
	for i := 0; i < b.N; i++ {
		BinaryProtocol.WriteI32(buf, 1)
	}
}

func BenchmarkBinaryProtocolWriteString4(b *testing.B) {
	buf := nullWriter(0)
	for i := 0; i < b.N; i++ {
		BinaryProtocol.WriteString(buf, "test")
	}
}

func BenchmarkBinaryProtocolBufferedReadByte(b *testing.B) {
	buf := bytes.NewBuffer(make([]byte, b.N))
	p := NewBinaryProtocol(true, false, 256)
	for i := 0; i < b.N; i++ {
		p.ReadByte(buf)
	}
}

func BenchmarkBinaryProtocolBufferedReadI32(b *testing.B) {
	buf := bytes.NewBuffer(make([]byte, b.N*4))
	p := NewBinaryProtocol(true, false, 256)
	for i := 0; i < b.N; i++ {
		p.ReadI32(buf)
	}
}

func BenchmarkBinaryProtocolBufferedWriteByte(b *testing.B) {
	buf := nullWriter(0)
	p := NewBinaryProtocol(true, false, 256)
	for i := 0; i < b.N; i++ {
		p.WriteByte(buf, 1)
	}
}

func BenchmarkBinaryProtocolBufferedWriteI32(b *testing.B) {
	buf := nullWriter(0)
	p := NewBinaryProtocol(true, false, 256)
	for i := 0; i < b.N; i++ {
		p.WriteI32(buf, 1)
	}
}

func BenchmarkBinaryProtocolBufferedWriteString4(b *testing.B) {
	buf := nullWriter(0)
	p := NewBinaryProtocol(true, false, 256)
	for i := 0; i < b.N; i++ {
		p.WriteString(buf, "test")
	}
}
