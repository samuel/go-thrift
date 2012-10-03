package thrift

import (
	"bytes"
	"testing"
)

func TestCompactProtocol(t *testing.T) {
	testProtocol(t, NewCompactProtocol())
}

func TestCompactI16(t *testing.T) {
	p := NewCompactProtocol()

	exp := map[int16][]byte{
		0:     []byte{0},
		-1:    []byte{1},
		1:     []byte{2},
		12345: []byte{242, 192, 1},
	}

	for expValue, expBytes := range exp {
		b := &bytes.Buffer{}
		err := p.WriteI16(b, expValue)
		if err != nil {
			t.Fatalf("WriteI16 returned an error: %+v", err)
		}
		out := b.Bytes()
		if bytes.Compare(out, expBytes) != 0 {
			t.Fatalf("CompactProtocol.WriteI16 wrote %+v which did not match expected %+v", out, expBytes)
		}

		b = bytes.NewBuffer(expBytes)
		v, err := p.ReadI16(b)
		if err != nil {
			t.Fatalf("ReadI16 returned an error: %+v", err)
		}
		if v != expValue {
			t.Fatalf("ReadI16 returned the wrong value %d instead of %d", v, expValue)
		}
	}
}

func TestCompactI32(t *testing.T) {
	p := NewCompactProtocol()

	exp := map[int32][]byte{
		0:          []byte{0},
		-1:         []byte{1},
		1:          []byte{2},
		1234567890: []byte{164, 139, 176, 153, 9},
	}

	for expValue, expBytes := range exp {
		b := &bytes.Buffer{}
		err := p.WriteI32(b, expValue)
		if err != nil {
			t.Fatalf("WriteI32 returned an error: %+v", err)
		}
		out := b.Bytes()
		if bytes.Compare(out, expBytes) != 0 {
			t.Fatalf("CompactProtocol.WriteI32 wrote %+v which did not match expected %+v", out, expBytes)
		}

		b = bytes.NewBuffer(expBytes)
		v, err := p.ReadI32(b)
		if err != nil {
			t.Fatalf("Read32 returned an error: %+v", err)
		}
		if v != expValue {
			t.Fatalf("Read32 returned the wrong value %d instead of %d", v, expValue)
		}
	}
}
