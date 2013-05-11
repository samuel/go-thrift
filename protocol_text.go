// Copyright 2012 Samuel Stauffer. All rights reserved.
// Use of this source code is governed by a 3-clause BSD
// license that can be found in the LICENSE file.

package thrift

import (
	"errors"
	"fmt"
	"io"
)

var (
	ErrUnimplemented = errors.New("thrift: unimplemented")
)

type textProtocol struct {
	indentation string
}

func NewTextProtocol() Protocol {
	return &textProtocol{}
}

func (p *textProtocol) indent() {
	p.indentation += "\t"
}

func (p *textProtocol) unindent() {
	p.indentation = p.indentation[:len(p.indentation)-1]
}

func (p *textProtocol) WriteMessageBegin(w io.Writer, name string, messageType byte, seqid int32) error {
	fmt.Fprintf(w, "%sMessageBegin(%s, %d, %.8x)\n", p.indentation, name, messageType, seqid)
	p.indent()
	return nil
}

func (p *textProtocol) WriteMessageEnd(w io.Writer) error {
	p.unindent()
	fmt.Fprintf(w, "%sMessageEnd()\n", p.indentation)
	return nil
}

func (p *textProtocol) WriteStructBegin(w io.Writer, name string) error {
	fmt.Fprintf(w, "%sStructBegin(%s)\n", p.indentation, name)
	p.indent()
	return nil
}

func (p *textProtocol) WriteStructEnd(w io.Writer) error {
	p.unindent()
	fmt.Fprintf(w, "%sStructEnd()\n", p.indentation)
	return nil
}

func (p *textProtocol) WriteFieldBegin(w io.Writer, name string, fieldType byte, id int16) error {
	fmt.Fprintf(w, "%sFieldBegin(%s, %d, %d)\n", p.indentation, name, fieldType, id)
	p.indent()
	return nil
}

func (p *textProtocol) WriteFieldEnd(w io.Writer) error {
	p.unindent()
	fmt.Fprintf(w, "%sFieldEnd()\n", p.indentation)
	return nil
}

func (p *textProtocol) WriteFieldStop(w io.Writer) error {
	fmt.Fprintf(w, "%sFieldStop()\n", p.indentation)
	return nil
}

func (p *textProtocol) WriteMapBegin(w io.Writer, keyType byte, valueType byte, size int) error {
	fmt.Fprintf(w, "%sMapBegin(%d, %d, %d)\n", p.indentation, keyType, valueType, size)
	p.indent()
	return nil
}

func (p *textProtocol) WriteMapEnd(w io.Writer) error {
	p.unindent()
	fmt.Fprintf(w, "%sMapEnd()\n", p.indentation)
	return nil
}

func (p *textProtocol) WriteListBegin(w io.Writer, elementType byte, size int) error {
	fmt.Fprintf(w, "%sListBegin(%d, %d)\n", p.indentation, elementType, size)
	p.indent()
	return nil
}

func (p *textProtocol) WriteListEnd(w io.Writer) error {
	p.unindent()
	fmt.Fprintf(w, "%sListEnd()\n", p.indentation)
	return nil
}

func (p *textProtocol) WriteSetBegin(w io.Writer, elementType byte, size int) error {
	fmt.Fprintf(w, "%sSetBegin(%d, %d)\n", p.indentation, elementType, size)
	p.indent()
	return nil
}

func (p *textProtocol) WriteSetEnd(w io.Writer) error {
	p.unindent()
	fmt.Fprintf(w, "%sSetEnd()\n", p.indentation)
	return nil
}

func (p *textProtocol) WriteBool(w io.Writer, value bool) error {
	fmt.Fprintf(w, "%sBool(%+v)\n", p.indentation, value)
	return nil
}

func (p *textProtocol) WriteByte(w io.Writer, value byte) error {
	fmt.Fprintf(w, "%sByte(%d)\n", p.indentation, value)
	return nil
}

func (p *textProtocol) WriteI16(w io.Writer, value int16) error {
	fmt.Fprintf(w, "%sI16(%d)\n", p.indentation, value)
	return nil
}

func (p *textProtocol) WriteI32(w io.Writer, value int32) error {
	fmt.Fprintf(w, "%sI32(%d)\n", p.indentation, value)
	return nil
}

func (p *textProtocol) WriteI64(w io.Writer, value int64) error {
	fmt.Fprintf(w, "%sI64(%d)\n", p.indentation, value)
	return nil
}

func (p *textProtocol) WriteDouble(w io.Writer, value float64) error {
	fmt.Fprintf(w, "%sDouble(%f)\n", p.indentation, value)
	return nil
}

func (p *textProtocol) WriteString(w io.Writer, value string) error {
	fmt.Fprintf(w, "%sString(%s)\n", p.indentation, value)
	return nil
}

func (p *textProtocol) WriteBytes(w io.Writer, value []byte) error {
	fmt.Fprintf(w, "%sBytes(%+v)\n", p.indentation, value)
	return nil
}

func (p *textProtocol) ReadMessageBegin(r io.Reader) (name string, messageType byte, seqid int32, err error) {
	return "", 0, 0, ErrUnimplemented
}

func (p *textProtocol) ReadMessageEnd(r io.Reader) error {
	return ErrUnimplemented
}

func (p *textProtocol) ReadStructBegin(r io.Reader) error {
	return ErrUnimplemented
}

func (p *textProtocol) ReadStructEnd(r io.Reader) error {
	return ErrUnimplemented
}

func (p *textProtocol) ReadFieldBegin(r io.Reader) (fieldType byte, id int16, err error) {
	return 0, 0, ErrUnimplemented
}

func (p *textProtocol) ReadFieldEnd(r io.Reader) error {
	return ErrUnimplemented
}

func (p *textProtocol) ReadMapBegin(r io.Reader) (keyType byte, valueType byte, size int, err error) {
	return 0, 0, 0, ErrUnimplemented
}

func (p *textProtocol) ReadMapEnd(r io.Reader) error {
	return ErrUnimplemented
}

func (p *textProtocol) ReadListBegin(r io.Reader) (elementType byte, size int, err error) {
	return 0, 0, ErrUnimplemented
}

func (p *textProtocol) ReadListEnd(r io.Reader) error {
	return ErrUnimplemented
}

func (p *textProtocol) ReadSetBegin(r io.Reader) (elementType byte, size int, err error) {
	return 0, 0, ErrUnimplemented
}

func (p *textProtocol) ReadSetEnd(r io.Reader) error {
	return ErrUnimplemented
}

func (p *textProtocol) ReadBool(r io.Reader) (bool, error) {
	return false, ErrUnimplemented
}

func (p *textProtocol) ReadByte(r io.Reader) (byte, error) {
	return 0, ErrUnimplemented
}

func (p *textProtocol) ReadI16(r io.Reader) (int16, error) {
	return 0, ErrUnimplemented
}

func (p *textProtocol) ReadI32(r io.Reader) (int32, error) {
	return 0, ErrUnimplemented
}

func (p *textProtocol) ReadI64(r io.Reader) (int64, error) {
	return 0, ErrUnimplemented
}

func (p *textProtocol) ReadDouble(r io.Reader) (float64, error) {
	return 0.0, ErrUnimplemented
}

func (p *textProtocol) ReadString(r io.Reader) (string, error) {
	return "", ErrUnimplemented
}

func (p *textProtocol) ReadBytes(r io.Reader) ([]byte, error) {
	return nil, ErrUnimplemented
}
