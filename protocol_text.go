package thrift

import (
	"errors"
	"fmt"
	"io"
)

var (
	ErrUnimplemented = errors.New("unimplemented")
)

type TextProtocol struct {
	indentation string
}

func (p *TextProtocol) indent() {
	p.indentation += "\t"
}

func (p *TextProtocol) unindent() {
	p.indentation = p.indentation[:len(p.indentation)-1]
}

func (p *TextProtocol) WriteMessageBegin(w io.Writer, name string, messageType byte, seqid int32) error {
	fmt.Fprintf(w, "%sMessageBegin(%s, %d, %.8x)\n", p.indentation, name, messageType, seqid)
	p.indent()
	return nil
}

func (p *TextProtocol) WriteMessageEnd(w io.Writer) error {
	p.unindent()
	fmt.Fprintf(w, "%sMessageEnd()\n", p.indentation)
	return nil
}

func (p *TextProtocol) WriteStructBegin(w io.Writer, name string) error {
	fmt.Fprintf(w, "%sStructBegin(%s)\n", p.indentation, name)
	p.indent()
	return nil
}

func (p *TextProtocol) WriteStructEnd(w io.Writer) error {
	p.unindent()
	fmt.Fprintf(w, "%sStructEnd()\n", p.indentation)
	return nil
}

func (p *TextProtocol) WriteFieldBegin(w io.Writer, name string, fieldType byte, id int16) error {
	fmt.Fprintf(w, "%sFieldBegin(%s, %d, %d)\n", p.indentation, name, fieldType, id)
	p.indent()
	return nil
}

func (p *TextProtocol) WriteFieldEnd(w io.Writer) error {
	p.unindent()
	fmt.Fprintf(w, "%sFieldEnd()\n", p.indentation)
	return nil
}

func (p *TextProtocol) WriteFieldStop(w io.Writer) error {
	fmt.Fprintf(w, "%sFieldStop()\n", p.indentation)
	return nil
}

func (p *TextProtocol) WriteMapBegin(w io.Writer, keyType byte, valueType byte, size int) error {
	fmt.Fprintf(w, "%sMapBegin(%d, %d, %d)\n", p.indentation, keyType, valueType, size)
	p.indent()
	return nil
}

func (p *TextProtocol) WriteMapEnd(w io.Writer) error {
	p.unindent()
	fmt.Fprintf(w, "%sMapEnd()\n", p.indentation)
	return nil
}

func (p *TextProtocol) WriteListBegin(w io.Writer, elementType byte, size int) error {
	fmt.Fprintf(w, "%sListBegin(%d, %d)\n", p.indentation, elementType, size)
	p.indent()
	return nil
}

func (p *TextProtocol) WriteListEnd(w io.Writer) error {
	p.unindent()
	fmt.Fprintf(w, "%sListEnd()\n", p.indentation)
	return nil
}

func (p *TextProtocol) WriteSetBegin(w io.Writer, elementType byte, size int) error {
	fmt.Fprintf(w, "%sSetBegin(%d, %d)\n", p.indentation, elementType, size)
	p.indent()
	return nil
}

func (p *TextProtocol) WriteSetEnd(w io.Writer) error {
	p.unindent()
	fmt.Fprintf(w, "%sSetEnd()\n", p.indentation)
	return nil
}

func (p *TextProtocol) WriteBool(w io.Writer, value bool) error {
	fmt.Fprintf(w, "%sBool(%+v)\n", p.indentation, value)
	return nil
}

func (p *TextProtocol) WriteByte(w io.Writer, value byte) error {
	fmt.Fprintf(w, "%sByte(%d)\n", p.indentation, value)
	return nil
}

func (p *TextProtocol) WriteI16(w io.Writer, value int16) error {
	fmt.Fprintf(w, "%sI16(%d)\n", p.indentation, value)
	return nil
}

func (p *TextProtocol) WriteI32(w io.Writer, value int32) error {
	fmt.Fprintf(w, "%sI32(%d)\n", p.indentation, value)
	return nil
}

func (p *TextProtocol) WriteI64(w io.Writer, value int64) error {
	fmt.Fprintf(w, "%sI64(%d)\n", p.indentation, value)
	return nil
}

func (p *TextProtocol) WriteDouble(w io.Writer, value float64) error {
	fmt.Fprintf(w, "%sDouble(%f)\n", p.indentation, value)
	return nil
}

func (p *TextProtocol) WriteString(w io.Writer, value string) error {
	fmt.Fprintf(w, "%sString(%s)\n", p.indentation, value)
	return nil
}

func (p *TextProtocol) WriteBytes(w io.Writer, value []byte) error {
	fmt.Fprintf(w, "%sBytes(%+v)\n", p.indentation, value)
	return nil
}

func (p *TextProtocol) ReadMessageBegin(r io.Reader) (name string, messageType byte, seqid int32, err error) {
	return "", 0, 0, ErrUnimplemented
}

func (p *TextProtocol) ReadMessageEnd(r io.Reader) error {
	return ErrUnimplemented
}

func (p *TextProtocol) ReadStructBegin(r io.Reader) error {
	return ErrUnimplemented
}

func (p *TextProtocol) ReadStructEnd(r io.Reader) error {
	return ErrUnimplemented
}

func (p *TextProtocol) ReadFieldBegin(r io.Reader) (fieldType byte, id int16, err error) {
	return 0, 0, ErrUnimplemented
}

func (p *TextProtocol) ReadFieldEnd(r io.Reader) error {
	return ErrUnimplemented
}

func (p *TextProtocol) ReadMapBegin(r io.Reader) (keyType byte, valueType byte, size int, err error) {
	return 0, 0, 0, ErrUnimplemented
}

func (p *TextProtocol) ReadMapEnd(r io.Reader) error {
	return ErrUnimplemented
}

func (p *TextProtocol) ReadListBegin(r io.Reader) (elementType byte, size int, err error) {
	return 0, 0, ErrUnimplemented
}

func (p *TextProtocol) ReadListEnd(r io.Reader) error {
	return ErrUnimplemented
}

func (p *TextProtocol) ReadSetBegin(r io.Reader) (elementType byte, size int, err error) {
	return 0, 0, ErrUnimplemented
}

func (p *TextProtocol) ReadSetEnd(r io.Reader) error {
	return ErrUnimplemented
}

func (p *TextProtocol) ReadBool(r io.Reader) (bool, error) {
	return false, ErrUnimplemented
}

func (p *TextProtocol) ReadByte(r io.Reader) (byte, error) {
	return 0, ErrUnimplemented
}

func (p *TextProtocol) ReadI16(r io.Reader) (int16, error) {
	return 0, ErrUnimplemented
}

func (p *TextProtocol) ReadI32(r io.Reader) (int32, error) {
	return 0, ErrUnimplemented
}

func (p *TextProtocol) ReadI64(r io.Reader) (int64, error) {
	return 0, ErrUnimplemented
}

func (p *TextProtocol) ReadDouble(r io.Reader) (float64, error) {
	return 0.0, ErrUnimplemented
}

func (p *TextProtocol) ReadString(r io.Reader) (string, error) {
	return "", ErrUnimplemented
}

func (p *TextProtocol) ReadBytes(r io.Reader) ([]byte, error) {
	return nil, ErrUnimplemented
}
