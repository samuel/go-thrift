Thrift Package for Go
=====================

[![Build Status](https://travis-ci.org/samuel/go-thrift.png)](https://travis-ci.org/samuel/go-thrift)

API Documentation: <http://godoc.org/github.com/samuel/go-thrift>

License
-------

3-clause BSD. See LICENSE file.

Overview
--------

So why another thrift package? While the existing one
([thrift4go](https://github.com/pomack/thrift4go/)) works well, my philosophy
is that interfaces should match the language. Most Thrift libraries try
to match the API of the original which makes them awkward to use in
other languages.

As an example, Go already has the idea of a thrift transport in the
ReadWriteCloser interfaces.

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
        SomeSet   []string `thrift:"5,set"`
    }

Types
-----

Most types map directly to the native Go types, but there are some
quirks and limitations.

* Go supports a more limited set of types for map keys than Thrift
* To use a set define the field as []type and provide a tag of "set":

        StringSet []string `thrift:"1,set"`

* []byte get encoded/decoded as a string because the Thrift binary type
  is the same as string on the wire.

RPC
---

The standard Go net/rpc package is used to provide RPC. Although, one
incompatibility is the net/rpc's use of ServiceName.Method for naming
RPC methods. To get around this the Thrift ServerCodec prefixes method
names with "Thrift".

### Transport

There are no specific transport "classes" as there are in most Thrift
libraries. Instead, the standard `io.ReadWriteCloser` is used as the
interface. If the value also implements the thrift.Flusher interface
then `Flush() error` is called after `protocol.WriteMessageEnd`.

_Framed transport_ is supported by wrapping a value implementing
`io.ReadWriteCloser` with `thrift.NewFramedReadWriteCloser(value)`

### One-way requests

#### Client

One-way request support needs to be enabled on the RPC codec explicitly.
The reason they're not allowed by default is because the Go RPC package
doesn't actually support one-way requests. To get around this requires
a rather janky hack of using channels to track pending requests in the
codec and faking responses.

#### Server

One-way requests aren't yet implemented on the server side.

Parser & Code Generator
-----------------------

The "parser" subdirectory contains a Thrift IDL parser, and "generator"
contains a Go code generator. It could be extended to include other
languages.

How to use the generator:

    $ go install github.com/samuel/go-thrift/generator

    $ generator --help
    Usage of parsimony: [options] inputfile outputpath
      -go.binarystring=false: Always use string for binary instead of []byte
      -go.json.enumnum=false: For JSON marshal enums by number instead of name
      -go.pointers=false: Make all fields pointers

    $ mkdir $GOPATH/src/cassandra
    $ generator cassandra.thrift $GOPATH/src/

TODO
----

* "extends" for services
* default values (is it worth it? maybe just annotations?)
