package thrift

import (
	"bytes"
	"testing"
)

type TestStruct struct {
	String string `thrift:"1"`
	Int    int    `thrift:"2"`
}

func TestDomainFilter(t *testing.T) {
	s := &TestStruct{"test", 123}
	enc := &Encoder{Writer: &bytes.Buffer{}}
	enc.WriteStruct(s)
}
