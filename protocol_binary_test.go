package thrift

import (
	"bytes"
	"testing"
)

func TestBinaryProtocol(t *testing.T) {
	b := &bytes.Buffer{}
	var pr Protocol = &BinaryProtocol{Writer: b, Reader: b, StrictWrite: true, StrictRead: false}

	if err := pr.WriteBool(true); err != nil {
		t.Fatalf("write bool true failed: %+v", err)
	}
	if b, err := pr.ReadBool(); err != nil {
		t.Fatalf("read bool true failed: %+v", err)
	} else if !b {
		t.Fatal("read bool true returned false")
	}

	if err := pr.WriteBool(false); err != nil {
		t.Fatalf("write bool false failed: %+v", err)
	}
	if b, err := pr.ReadBool(); err != nil {
		t.Fatalf("read bool false failed: %+v", err)
	} else if b {
		t.Fatal("read bool false returned true")
	}

	if err := pr.WriteI16(1234); err != nil {
		t.Fatalf("write i16 failed: %+v", err)
	}
	if v, err := pr.ReadI16(); err != nil {
		t.Fatalf("read i16 failed: %+v", err)
	} else if v != 1234 {
		t.Fatal("read i16 returned %d expected 1234", v)
	}

	if err := pr.WriteI32(-1234); err != nil {
		t.Fatalf("write i32 failed: %+v", err)
	}
	if v, err := pr.ReadI32(); err != nil {
		t.Fatalf("read i32 failed: %+v", err)
	} else if v != -1234 {
		t.Fatal("read i32 returned %d expected -1234", v)
	}

	if err := pr.WriteI64(-1234); err != nil {
		t.Fatalf("write i64 failed: %+v", err)
	}
	if v, err := pr.ReadI64(); err != nil {
		t.Fatalf("read i64 failed: %+v", err)
	} else if v != -1234 {
		t.Fatal("read i64 returned %d expected -1234", v)
	}

	if err := pr.WriteDouble(-0.1234); err != nil {
		t.Fatalf("write double failed: %+v", err)
	}
	if v, err := pr.ReadDouble(); err != nil {
		t.Fatalf("read double failed: %+v", err)
	} else if v != -0.1234 {
		t.Fatal("read double returned %.4f expected -0.1234", v)
	}

	testString := "012345"
	for i := 0; i < 2; i += 1 {
		if err := pr.WriteString(testString); err != nil {
			t.Fatalf("write string failed: %+v", err)
		}
		if v, err := pr.ReadString(); err != nil {
			t.Fatalf("read string failed: %+v", err)
		} else if v != testString {
			t.Fatal("read string returned %s expected '%s'", v, testString)
		}
		testString += "012345"
	}
}

// func BenchmarkReadByte(b *testing.B) {
// 	// buf := bytes.NewBuffer(make([]byte, b.N*8))
// 	// rd := BinaryProtocolReader{Reader: buf, Strict: false}
// 	// for i := 0; i < b.N; i++ {
// 	// 	pr.ReadI16()
// 	// }

// 	buf := &bytes.Buffer{}
// 	pr := BinaryProtocol{Writer: buf, StrictWrite: true}
// 	for i := 0; i < b.N; i++ {
// 		pr.WriteI32(1)
// 	}
// }
