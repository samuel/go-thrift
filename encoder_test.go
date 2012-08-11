package thrift

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
)

type TestStruct2 struct {
	Str    string `thrift:"1"`
	Binary []byte `thrift:"2"`
}

func (t *TestStruct2) String() string {
	return fmt.Sprintf("{Str:%s Binary:%+v}", t.Str, t.Binary)
}

type TestStruct struct {
	String  string            `thrift:"1"`
	Int     *int              `thrift:"2"`
	List    []string          `thrift:"3"`
	Map     map[string]string `thrift:"4"`
	Struct  *TestStruct2      `thrift:"5"`
	List2   []*string         `thrift:"6"`
	Struct2 TestStruct2       `thrift:"7"`
	Binary  []byte            `thrift:"8"`
	Set     []string          `thrift:"9,set"`
}

func TestKeepEmpty(t *testing.T) {
	buf := &bytes.Buffer{}

	s := struct {
		Str1 string `thrift:"1"`
	}{}
	err := EncodeStruct(buf, DefaultBinaryProtocol, s)
	if err != nil {
		t.Fatal(err)
	}
	if buf.Len() != 1 || buf.Bytes()[0] != 0 {
		t.Fatal("missing keepempty should mean empty fields are not serialized")
	}

	buf.Reset()
	s2 := struct {
		Str1 string `thrift:"1,keepempty"`
	}{}
	err = EncodeStruct(buf, DefaultBinaryProtocol, s2)
	if err != nil {
		t.Fatal(err)
	}
	if buf.Len() != 8 {
		t.Fatal("keepempty should cause empty fields to be serialized")
	}
}

func TestEncodeRequired(t *testing.T) {
	buf := &bytes.Buffer{}

	s := struct {
		Str1 string `thrift:"1,required"`
	}{}
	err := EncodeStruct(buf, DefaultBinaryProtocol, s)
	if err != nil {
		t.Fatal(err)
	}
	if buf.Len() != 8 {
		t.Fatal("Non-pointer required fields that aren't 'keepempty' should be serialized empty")
	}

	buf.Reset()
	s2 := struct {
		Str1 *string `thrift:"1,required"`
	}{}
	err = EncodeStruct(buf, DefaultBinaryProtocol, s2)
	_, ok := err.(*MissingRequiredField)
	if !ok {
		t.Fatalf("Missing required field should throw MissingRequiredField instead of %+v", err)
	}
}

func TestBasics(t *testing.T) {
	i := 123
	str := "bar"
	ts2 := TestStruct2{"qwerty", []byte{1, 2, 3}}
	s := &TestStruct{
		"test",
		&i,
		[]string{"a", "b"},
		map[string]string{"a": "b", "1": "2"},
		&ts2,
		[]*string{&str},
		ts2,
		[]byte{1, 2, 3},
		[]string{"a", "b"},
	}
	buf := &bytes.Buffer{}

	err := EncodeStruct(buf, DefaultBinaryProtocol, s)
	if err != nil {
		t.Fatal(err)
	}

	s2 := &TestStruct{}
	err = DecodeStruct(buf, DefaultBinaryProtocol, s2)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(s, s2) {
		t.Fatalf("encdec doesn't match: %+v != %+v", s, s2)
	}
}
