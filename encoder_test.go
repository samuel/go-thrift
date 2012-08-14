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

type TestStructRequiredOptional struct {
	RequiredPtr *string `thrift:"1,required"`
	Required    string  `thrift:"2,required"`
	OptionalPtr *string `thrift:"3"`
	Optional    string  `thrift:"4"`
}

type TestEmptyStruct struct{}

func TestKeepEmpty(t *testing.T) {
	buf := &bytes.Buffer{}

	s := struct {
		Str1 string `thrift:"1"`
	}{}
	err := EncodeStruct(buf, BinaryProtocol, s)
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
	err = EncodeStruct(buf, BinaryProtocol, s2)
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
	err := EncodeStruct(buf, BinaryProtocol, s)
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
	err = EncodeStruct(buf, BinaryProtocol, s2)
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

	err := EncodeStruct(buf, BinaryProtocol, s)
	if err != nil {
		t.Fatal(err)
	}

	s2 := &TestStruct{}
	err = DecodeStruct(buf, BinaryProtocol, s2)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(s, s2) {
		t.Fatalf("encdec doesn't match: %+v != %+v", s, s2)
	}
}

func TestEncodeRequiredFields(t *testing.T) {
	buf := &bytes.Buffer{}

	// encode nil pointer required field

	s := &TestStructRequiredOptional{nil, "", nil, ""}
	err := EncodeStruct(buf, BinaryProtocol, s)
	if err == nil {
		t.Fatal("Expected MissingRequiredField exception")
	}
	e, ok := err.(*MissingRequiredField)
	if !ok {
		t.Fatalf("Expected MissingRequiredField exception instead %+v", err)
	}
	if e.StructName != "TestStructRequiredOptional" || e.FieldName != "RequiredPtr" {
		t.Fatalf("Expected MissingRequiredField{'TestStructRequiredOptional', 'RequiredPtr'} instead %+v", e)
	}

	// encode empty non-pointer required field

	str := "foo"
	s = &TestStructRequiredOptional{&str, "", nil, ""}
	err = EncodeStruct(buf, BinaryProtocol, s)
	if err != nil {
		t.Fatal("Empty non-pointer required fields shouldn't return an error")
	}
}

func TestDecodeRequiredFields(t *testing.T) {
	buf := &bytes.Buffer{}

	s := &TestEmptyStruct{}
	err := EncodeStruct(buf, BinaryProtocol, s)
	if err != nil {
		t.Fatal("Failed to encode empty struct")
	}

	s2 := &TestStructRequiredOptional{}
	err = DecodeStruct(buf, BinaryProtocol, s2)
	if err == nil {
		t.Fatal("Expected MissingRequiredField exception")
	}
	e, ok := err.(*MissingRequiredField)
	if !ok {
		t.Fatalf("Expected MissingRequiredField exception instead %+v", err)
	}
	if e.StructName != "TestStructRequiredOptional" || e.FieldName != "RequiredPtr" {
		t.Fatalf("Expected MissingRequiredField{'TestStructRequiredOptional', 'RequiredPtr'} instead %+v", e)
	}
}

func TestDecodeUnknownFields(t *testing.T) {
	buf := &bytes.Buffer{}

	str := "foo"
	s := &TestStructRequiredOptional{&str, str, &str, str}
	err := EncodeStruct(buf, BinaryProtocol, s)
	if err != nil {
		t.Fatal("Failed to encode TestStructRequiredOptional struct")
	}

	s2 := &TestEmptyStruct{}
	err = DecodeStruct(buf, BinaryProtocol, s2)
	if err != nil {
		t.Fatalf("Unknown fields during decode weren't ignored: %+v", err)
	}
}

func BenchmarkEncodeEmptyStruct(b *testing.B) {
	buf := nullWriter(0)
	st := &struct{}{}
	for i := 0; i < b.N; i++ {
		EncodeStruct(buf, BinaryProtocol, st)
	}
}

func BenchmarkDecodeEmptyStruct(b *testing.B) {
	b.StopTimer()
	buf1 := &bytes.Buffer{}
	st := &struct{}{}
	EncodeStruct(buf1, BinaryProtocol, st)
	buf := bytes.NewBuffer(bytes.Repeat(buf1.Bytes(), b.N))
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		DecodeStruct(buf, BinaryProtocol, st)
	}
}

func BenchmarkEncodeSimpleStruct(b *testing.B) {
	buf := nullWriter(0)
	st := &struct {
		Str string `thrift:"1,required"`
		Int int32  `thrift:"2,required"`
	}{
		Str: "test",
		Int: 123,
	}
	for i := 0; i < b.N; i++ {
		EncodeStruct(buf, BinaryProtocol, st)
	}
}
