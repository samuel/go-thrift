package thrift

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
)

type TestStruct2 struct {
	Str string `thrift:"1"`
}

func (t *TestStruct2) String() string {
	return fmt.Sprintf("{Str:%s}", t.Str)
}

type TestStruct struct {
	String string            `thrift:"1"`
	Int    *int              `thrift:"2"`
	List   []string          `thrift:"3"`
	Map    map[string]string `thrift:"4"`
	Struct *TestStruct2      `thrift:"5"`
	List2  []*string         `thrift:"6"`
}

func TestDomainFilter(t *testing.T) {
	i := 123
	y := "bar"
	o := TestStruct2{"qwerty"}
	s := &TestStruct{
		"test",
		&i,
		[]string{"a", "b"},
		map[string]string{"a": "b", "1": "2"},
		&o,
		[]*string{&y},
	}
	buf := &bytes.Buffer{}
	p := &BinaryProtocol{Writer: buf, Reader: buf, StrictWrite: true, StrictRead: false}

	enc := &Encoder{Protocol: p}
	err := enc.WriteStruct(s)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v\n", buf.Bytes())

	s2 := &TestStruct{}
	dec := &Decoder{Protocol: p}
	err = dec.ReadStruct(s2)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(s, s2) {
		t.Fatalf("encdec doesn't match: %+v != %+v", s, s2)
	}
}
