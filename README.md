Thrift Package for Go
=====================

WARNING
-------

*This package is a work in progress.*

Overview
--------

So why another thrift package? While the existing one
([thrift4go](https://github.com/pomack/thrift4go/)) works well, my philosophy
is that interface should match the language rather than being standardized
across disparate styles.

As an example, Go already has the idea of a thrift transport in the
Reader/Writer interfaces.

Another design decision was to keep the generated code as terse as possible.
The generator only creates a struct and the encoding/decoding is done through
reflection. Annotations are used to set thrift ID for a field and options such
as 'required'.

Example struct:

    type User struct {
        Id        int64    `thrift:"1,required"`
        Name      string   `thrift:"2"`
        PostCount int32    `thrift:"3,keepempty"`
        Flags     []string `thrift:"4"`
    }