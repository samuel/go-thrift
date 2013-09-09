// Copyright 2012 Samuel Stauffer. All rights reserved.
// Use of this source code is governed by a 3-clause BSD
// license that can be found in the LICENSE file.

package thrift

import (
	"bytes"
	"testing"
)

func testProtocol(t *testing.T, pr Protocol) {
	b := &bytes.Buffer{}

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
		t.Fatalf("read i16 returned %d expected 1234", v)
	}

	if err := pr.WriteI32(b, -1234); err != nil {
		t.Fatalf("write i32 failed: %+v", err)
	}
	if v, err := pr.ReadI32(b); err != nil {
		t.Fatalf("read i32 failed: %+v", err)
	} else if v != -1234 {
		t.Fatalf("read i32 returned %d expected -1234", v)
	}

	if err := pr.WriteI64(b, -1234); err != nil {
		t.Fatalf("write i64 failed: %+v", err)
	}
	if v, err := pr.ReadI64(b); err != nil {
		t.Fatalf("read i64 failed: %+v", err)
	} else if v != -1234 {
		t.Fatalf("read i64 returned %d expected -1234", v)
	}

	if err := pr.WriteDouble(b, -0.1234); err != nil {
		t.Fatalf("write double failed: %+v", err)
	}
	if v, err := pr.ReadDouble(b); err != nil {
		t.Fatalf("read double failed: %+v", err)
	} else if v != -0.1234 {
		t.Fatalf("read double returned %.4f expected -0.1234", v)
	}

	testString := "012345"
	for i := 0; i < 2; i += 1 {
		if err := pr.WriteString(b, testString); err != nil {
			t.Fatalf("write string failed: %+v", err)
		}
		if v, err := pr.ReadString(b); err != nil {
			t.Fatalf("read string failed: %+v", err)
		} else if v != testString {
			t.Fatalf("read string returned %s expected '%s'", v, testString)
		}
		testString += "012345"
	}

	// Write a message

	if err := pr.WriteMessageBegin(b, "msgName", 2, 123); err != nil {
		t.Fatalf("WriteMessageBegin failed: %+v", err)
	}
	if err := pr.WriteStructBegin(b, "struct"); err != nil {
		t.Fatalf("WriteStructBegin failed: %+v", err)
	}

	if err := pr.WriteFieldBegin(b, "boolTrue", TypeBool, 1); err != nil {
		t.Fatalf("WriteFieldBegin failed: %+v", err)
	}
	if err := pr.WriteBool(b, true); err != nil {
		t.Fatalf("WriteBool(true) failed: %+v", err)
	}
	if err := pr.WriteFieldEnd(b); err != nil {
		t.Fatalf("WriteFieldEnd failed: %+v", err)
	}

	if err := pr.WriteFieldBegin(b, "boolFalse", TypeBool, 3); err != nil {
		t.Fatalf("WriteFieldBegin failed: %+v", err)
	}
	if err := pr.WriteBool(b, false); err != nil {
		t.Fatalf("WriteBool(false) failed: %+v", err)
	}
	if err := pr.WriteFieldEnd(b); err != nil {
		t.Fatalf("WriteFieldEnd failed: %+v", err)
	}

	if err := pr.WriteFieldBegin(b, "str", TypeString, 2); err != nil {
		t.Fatalf("WriteFieldBegin failed: %+v", err)
	}
	if err := pr.WriteString(b, "foo"); err != nil {
		t.Fatalf("WriteString failed: %+v", err)
	}
	if err := pr.WriteFieldEnd(b); err != nil {
		t.Fatalf("WriteFieldEnd failed: %+v", err)
	}

	if err := pr.WriteFieldStop(b); err != nil {
		t.Fatalf("WriteStructEnd failed: %+v", err)
	}
	if err := pr.WriteStructEnd(b); err != nil {
		t.Fatalf("WriteStructEnd failed: %+v", err)
	}
	if err := pr.WriteMessageEnd(b); err != nil {
		t.Fatalf("WriteMessageEnd failed: %+v", err)
	}

	// Read the message

	if name, mtype, seqId, err := pr.ReadMessageBegin(b); err != nil {
		t.Fatalf("ReadMessageBegin failed: %+v", err)
	} else if name != "msgName" {
		t.Fatalf("ReadMessageBegin name mismatch: %s != %s", name, "msgName")
	} else if mtype != 2 {
		t.Fatalf("ReadMessageBegin type mismatch: %d != %d", mtype, 2)
	} else if seqId != 123 {
		t.Fatalf("ReadMessageBegin seqId mismatch: %d != %d", seqId, 123)
	}
	if err := pr.ReadStructBegin(b); err != nil {
		t.Fatalf("ReadStructBegin failed: %+v", err)
	}

	if fieldType, id, err := pr.ReadFieldBegin(b); err != nil {
		t.Fatalf("ReadFieldBegin failed: %+v", err)
	} else if fieldType != TypeBool {
		t.Fatalf("ReadFieldBegin type mismatch: %d != %d", fieldType, TypeBool)
	} else if id != 1 {
		t.Fatalf("ReadFieldBegin id mismatch: %d != %d", id, 1)
	}
	if v, err := pr.ReadBool(b); err != nil {
		t.Fatalf("ReaBool failed: %+v", err)
	} else if !v {
		t.Fatalf("ReadBool value mistmatch %+v != %+v", v, true)
	}
	if err := pr.ReadFieldEnd(b); err != nil {
		t.Fatalf("ReadFieldEnd failed: %+v", err)
	}

	if fieldType, id, err := pr.ReadFieldBegin(b); err != nil {
		t.Fatalf("ReadFieldBegin failed: %+v", err)
	} else if fieldType != TypeBool {
		t.Fatalf("ReadFieldBegin type mismatch: %d != %d", fieldType, TypeBool)
	} else if id != 3 {
		t.Fatalf("ReadFieldBegin id mismatch: %d != %d", id, 3)
	}
	if v, err := pr.ReadBool(b); err != nil {
		t.Fatalf("ReaBool failed: %+v", err)
	} else if v {
		t.Fatalf("ReadBool value mistmatch %+v != %+v", v, false)
	}
	if err := pr.ReadFieldEnd(b); err != nil {
		t.Fatalf("ReadFieldEnd failed: %+v", err)
	}

	if fieldType, id, err := pr.ReadFieldBegin(b); err != nil {
		t.Fatalf("ReadFieldBegin failed: %+v", err)
	} else if fieldType != TypeString {
		t.Fatalf("ReadFieldBegin type mismatch: %d != %d", fieldType, TypeString)
	} else if id != 2 {
		t.Fatalf("ReadFieldBegin id mismatch: %d != %d", id, 2)
	}
	if v, err := pr.ReadString(b); err != nil {
		t.Fatalf("ReadString failed: %+v", err)
	} else if v != "foo" {
		t.Fatalf("ReadString value mistmatch %s != %s", v, "foo")
	}
	if err := pr.ReadFieldEnd(b); err != nil {
		t.Fatalf("ReadFieldEnd failed: %+v", err)
	}

	if err := pr.ReadStructEnd(b); err != nil {
		t.Fatalf("ReadStructEnd failed: %+v", err)
	}
	if err := pr.ReadMessageEnd(b); err != nil {
		t.Fatalf("ReadMessageEnd failed: %+v", err)
	}
}

func TestBinaryProtocolBadStringLength(t *testing.T) {
	b := &bytes.Buffer{}
	pr := NewBinaryProtocol(true, false)

	// zero string length
	if err := pr.WriteI32(b, 0); err != nil {
		t.Fatal(err)
	}
	if st, err := pr.ReadString(b); err != nil {
		t.Fatal(err)
	} else if st != "" {
		t.Fatal("BinaryProtocol.ReadString didn't return an empty string given a length of 0")
	}

	// negative string length
	if err := pr.WriteI32(b, -1); err != nil {
		t.Fatal(err)
	}
	if _, err := pr.ReadString(b); err == nil {
		t.Fatal("BinaryProtocol.ReadString didn't return an error given a negative length")
	}
}

func TestBinaryProtocol(t *testing.T) {
	testProtocol(t, NewBinaryProtocol(true, false))
}

func BenchmarkBinaryProtocolReadByte(b *testing.B) {
	buf := &loopingReader{}
	p := NewBinaryProtocol(true, false)
	p.WriteByte(buf, 123)
	for i := 0; i < b.N; i++ {
		p.ReadByte(buf)
	}
}

func BenchmarkBinaryProtocolReadI32(b *testing.B) {
	buf := &loopingReader{}
	p := NewBinaryProtocol(true, false)
	p.WriteI32(buf, 1234567890)
	for i := 0; i < b.N; i++ {
		p.ReadI32(buf)
	}
}

func BenchmarkBinaryProtocolWriteByte(b *testing.B) {
	buf := nullWriter(0)
	p := NewBinaryProtocol(true, false)
	for i := 0; i < b.N; i++ {
		p.WriteByte(buf, 1)
	}
}

func BenchmarkBinaryProtocolWriteI32(b *testing.B) {
	buf := nullWriter(0)
	p := NewBinaryProtocol(true, false)
	for i := 0; i < b.N; i++ {
		p.WriteI32(buf, 1)
	}
}

func BenchmarkBinaryProtocolWriteString4(b *testing.B) {
	buf := nullWriter(0)
	p := NewBinaryProtocol(true, false)
	for i := 0; i < b.N; i++ {
		p.WriteString(buf, "test")
	}
}

func BenchmarkBinaryProtocolWriteFullMessage(b *testing.B) {
	buf := nullWriter(0)
	p := NewBinaryProtocol(true, false)
	for i := 0; i < b.N; i++ {
		p.WriteMessageBegin(buf, "", 2, 123)
		p.WriteStructBegin(buf, "")
		p.WriteFieldBegin(buf, "", TypeBool, 1)
		p.WriteBool(buf, true)
		p.WriteFieldEnd(buf)
		p.WriteFieldBegin(buf, "", TypeBool, 3)
		p.WriteBool(buf, false)
		p.WriteFieldEnd(buf)
		p.WriteFieldBegin(buf, "", TypeString, 2)
		p.WriteString(buf, "foo")
		p.WriteFieldEnd(buf)
		p.WriteFieldStop(buf)
		p.WriteStructEnd(buf)
		p.WriteMessageEnd(buf)
	}
}
