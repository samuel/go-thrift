// Copyright 2012 Samuel Stauffer. All rights reserved.
// Use of this source code is governed by a 3-clause BSD
// license that can be found in the LICENSE file.

package thrift

import (
	"fmt"
	"io"
)

type ProtocolError struct {
	Protocol string
	Message  string
}

func (e ProtocolError) Error() string {
	return fmt.Sprintf("thrift: [%s] %s", e.Protocol, e.Message)
}

type Protocol interface {
	WriteMessageBegin(w io.Writer, name string, messageType byte, seqid int32) error
	WriteMessageEnd(w io.Writer) error
	WriteStructBegin(w io.Writer, name string) error
	WriteStructEnd(w io.Writer) error
	WriteFieldBegin(w io.Writer, name string, fieldType byte, id int16) error
	WriteFieldEnd(w io.Writer) error
	WriteFieldStop(w io.Writer) error
	WriteMapBegin(w io.Writer, keyType byte, valueType byte, size int) error
	WriteMapEnd(w io.Writer) error
	WriteListBegin(w io.Writer, elementType byte, size int) error
	WriteListEnd(w io.Writer) error
	WriteSetBegin(w io.Writer, elementType byte, size int) error
	WriteSetEnd(w io.Writer) error
	WriteBool(w io.Writer, value bool) error
	WriteByte(w io.Writer, value byte) error
	WriteI16(w io.Writer, value int16) error
	WriteI32(w io.Writer, value int32) error
	WriteI64(w io.Writer, value int64) error
	WriteDouble(w io.Writer, value float64) error
	WriteString(w io.Writer, value string) error
	WriteBytes(w io.Writer, value []byte) error

	ReadMessageBegin(r io.Reader) (name string, messageType byte, seqid int32, err error)
	ReadMessageEnd(r io.Reader) error
	ReadStructBegin(r io.Reader) error
	ReadStructEnd(r io.Reader) error
	ReadFieldBegin(r io.Reader) (fieldType byte, id int16, err error)
	ReadFieldEnd(r io.Reader) error
	ReadMapBegin(r io.Reader) (keyType byte, valueType byte, size int, err error)
	ReadMapEnd(r io.Reader) error
	ReadListBegin(r io.Reader) (elementType byte, size int, err error)
	ReadListEnd(r io.Reader) error
	ReadSetBegin(r io.Reader) (elementType byte, size int, err error)
	ReadSetEnd(r io.Reader) error
	ReadBool(r io.Reader) (bool, error)
	ReadByte(r io.Reader) (byte, error)
	ReadI16(r io.Reader) (int16, error)
	ReadI32(r io.Reader) (int32, error)
	ReadI64(r io.Reader) (int64, error)
	ReadDouble(r io.Reader) (float64, error)
	ReadString(r io.Reader) (string, error)
	ReadBytes(r io.Reader) ([]byte, error)
}
