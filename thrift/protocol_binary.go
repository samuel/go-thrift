// Copyright 2012 Samuel Stauffer. All rights reserved.
// Use of this source code is governed by a 3-clause BSD
// license that can be found in the LICENSE file.

package thrift

import (
	"encoding/binary"
	"io"
	"math"
)

const (
	versionMask uint32 = 0xffff0000
	version1    uint32 = 0x80010000
	typeMask    uint32 = 0x000000ff
)

const (
	maxMessageNameSize = 128
)

type binaryProtocol struct {
	strictWrite bool
	strictRead  bool
	writeBuf    []byte
	readBuf     []byte
}

func NewBinaryProtocol(strictWrite bool, strictRead bool) Protocol {
	p := &binaryProtocol{
		strictWrite: strictWrite,
		strictRead:  strictRead,
		writeBuf:    make([]byte, 32),
		readBuf:     make([]byte, 32),
	}
	return p
}

func (p *binaryProtocol) WriteMessageBegin(w io.Writer, name string, messageType byte, seqid int32) error {
	if p.strictWrite {
		if err := p.WriteI32(w, int32(version1|uint32(messageType))); err != nil {
			return err
		}
		if err := p.WriteString(w, name); err != nil {
			return err
		}
	} else {
		if err := p.WriteString(w, name); err != nil {
			return err
		}
		if err := p.WriteByte(w, messageType); err != nil {
			return err
		}
	}
	return p.WriteI32(w, seqid)
}

func (p *binaryProtocol) WriteMessageEnd(w io.Writer) error {
	return nil
}

func (p *binaryProtocol) WriteStructBegin(w io.Writer, name string) error {
	return nil
}

func (p *binaryProtocol) WriteStructEnd(w io.Writer) error {
	return nil
}

func (p *binaryProtocol) WriteFieldBegin(w io.Writer, name string, fieldType byte, id int16) error {
	if err := p.WriteByte(w, fieldType); err != nil {
		return err
	}
	return p.WriteI16(w, id)
}

func (p *binaryProtocol) WriteFieldEnd(w io.Writer) error {
	return nil
}

func (p *binaryProtocol) WriteFieldStop(w io.Writer) error {
	return p.WriteByte(w, TypeStop)
}

func (p *binaryProtocol) WriteMapBegin(w io.Writer, keyType byte, valueType byte, size int) error {
	if err := p.WriteByte(w, keyType); err != nil {
		return err
	}
	if err := p.WriteByte(w, valueType); err != nil {
		return err
	}
	return p.WriteI32(w, int32(size))
}

func (p *binaryProtocol) WriteMapEnd(w io.Writer) error {
	return nil
}

func (p *binaryProtocol) WriteListBegin(w io.Writer, elementType byte, size int) error {
	if err := p.WriteByte(w, elementType); err != nil {
		return err
	}
	return p.WriteI32(w, int32(size))
}

func (p *binaryProtocol) WriteListEnd(w io.Writer) error {
	return nil
}

func (p *binaryProtocol) WriteSetBegin(w io.Writer, elementType byte, size int) error {
	if err := p.WriteByte(w, elementType); err != nil {
		return err
	}
	return p.WriteI32(w, int32(size))
}

func (p *binaryProtocol) WriteSetEnd(w io.Writer) error {
	return nil
}

func (p *binaryProtocol) WriteBool(w io.Writer, value bool) error {
	if value {
		return p.WriteByte(w, 1)
	}
	return p.WriteByte(w, 0)
}

func (p *binaryProtocol) WriteByte(w io.Writer, value byte) error {
	b := p.writeBuf
	if b == nil {
		b = []byte{value}
	} else {
		b[0] = value
	}
	_, err := w.Write(b[:1])
	return err
}

func (p *binaryProtocol) WriteI16(w io.Writer, value int16) (err error) {
	b := p.writeBuf
	if b == nil {
		b = []byte{0, 0}
	}
	binary.BigEndian.PutUint16(b, uint16(value))
	_, err = w.Write(b[:2])
	return
}

func (p *binaryProtocol) WriteI32(w io.Writer, value int32) (err error) {
	b := p.writeBuf
	if b == nil {
		b = []byte{0, 0, 0, 0}
	}
	binary.BigEndian.PutUint32(b, uint32(value))
	_, err = w.Write(b[:4])
	return
}

func (p *binaryProtocol) WriteI64(w io.Writer, value int64) (err error) {
	b := p.writeBuf
	if b == nil {
		b = []byte{0, 0, 0, 0, 0, 0, 0, 0}
	}
	binary.BigEndian.PutUint64(b, uint64(value))
	_, err = w.Write(b[:8])
	return
}

func (p *binaryProtocol) WriteDouble(w io.Writer, value float64) (err error) {
	b := p.writeBuf
	if b == nil {
		b = []byte{0, 0, 0, 0, 0, 0, 0, 0}
	}
	binary.BigEndian.PutUint64(b, math.Float64bits(value))
	_, err = w.Write(b[:8])
	return
}

func (p *binaryProtocol) WriteString(w io.Writer, value string) error {
	if len(value) <= len(p.writeBuf) {
		if err := p.WriteI32(w, int32(len(value))); err != nil {
			return err
		}
		n := copy(p.writeBuf, value)
		_, err := w.Write(p.writeBuf[:n])
		return err
	}
	return p.WriteBytes(w, []byte(value))
}

func (p *binaryProtocol) WriteBytes(w io.Writer, value []byte) error {
	if err := p.WriteI32(w, int32(len(value))); err != nil {
		return err
	}
	_, err := w.Write(value)
	return err
}

func (p *binaryProtocol) ReadMessageBegin(r io.Reader) (name string, messageType byte, seqid int32, err error) {
	size, e := p.ReadI32(r)
	if e != nil {
		err = e
		return
	}
	if size < 0 {
		version := uint32(size) & versionMask
		if version != version1 {
			err = ProtocolError{"BinaryProtocol", "bad version in ReadMessageBegin"}
			return
		}
		messageType = byte(uint32(size) & typeMask)
		if name, err = p.ReadString(r); err != nil {
			return
		}
	} else {
		if p.strictRead {
			err = ProtocolError{"BinaryProtocol", "no protocol version header"}
			return
		}
		if size > maxMessageNameSize {
			err = ProtocolError{"BinaryProtocol", "message name exceeds max size"}
			return
		}
		nameBytes := make([]byte, size)
		if _, err = r.Read(nameBytes); err != nil {
			return
		}
		name = string(nameBytes)
		if messageType, err = p.ReadByte(r); err != nil {
			return
		}
	}
	seqid, err = p.ReadI32(r)
	return
}

func (p *binaryProtocol) ReadMessageEnd(r io.Reader) error {
	return nil
}

func (p *binaryProtocol) ReadStructBegin(r io.Reader) error {
	return nil
}

func (p *binaryProtocol) ReadStructEnd(r io.Reader) error {
	return nil
}

func (p *binaryProtocol) ReadFieldBegin(r io.Reader) (fieldType byte, id int16, err error) {
	if fieldType, err = p.ReadByte(r); err != nil || fieldType == TypeStop {
		return
	}
	id, err = p.ReadI16(r)
	return
}

func (p *binaryProtocol) ReadFieldEnd(r io.Reader) error {
	return nil
}

func (p *binaryProtocol) ReadMapBegin(r io.Reader) (keyType byte, valueType byte, size int, err error) {
	if keyType, err = p.ReadByte(r); err != nil {
		return
	}
	if valueType, err = p.ReadByte(r); err != nil {
		return
	}
	var sz int32
	sz, err = p.ReadI32(r)
	size = int(sz)
	return
}

func (p *binaryProtocol) ReadMapEnd(r io.Reader) error {
	return nil
}

func (p *binaryProtocol) ReadListBegin(r io.Reader) (elementType byte, size int, err error) {
	if elementType, err = p.ReadByte(r); err != nil {
		return
	}
	var sz int32
	sz, err = p.ReadI32(r)
	size = int(sz)
	return
}

func (p *binaryProtocol) ReadListEnd(r io.Reader) error {
	return nil
}

func (p *binaryProtocol) ReadSetBegin(r io.Reader) (elementType byte, size int, err error) {
	if elementType, err = p.ReadByte(r); err != nil {
		return
	}
	var sz int32
	sz, err = p.ReadI32(r)
	size = int(sz)
	return
}

func (p *binaryProtocol) ReadSetEnd(r io.Reader) error {
	return nil
}

func (p *binaryProtocol) ReadBool(r io.Reader) (bool, error) {
	if b, e := p.ReadByte(r); e != nil {
		return false, e
	} else if b != 0 {
		return true, nil
	}
	return false, nil
}

func (p *binaryProtocol) ReadByte(r io.Reader) (value byte, err error) {
	_, err = io.ReadFull(r, p.readBuf[:1])
	value = p.readBuf[0]
	return
}

func (p *binaryProtocol) ReadI16(r io.Reader) (value int16, err error) {
	_, err = io.ReadFull(r, p.readBuf[:2])
	value = int16(binary.BigEndian.Uint16(p.readBuf))
	return
}

func (p *binaryProtocol) ReadI32(r io.Reader) (value int32, err error) {
	_, err = io.ReadFull(r, p.readBuf[:4])
	value = int32(binary.BigEndian.Uint32(p.readBuf))
	return
}

func (p *binaryProtocol) ReadI64(r io.Reader) (value int64, err error) {
	_, err = io.ReadFull(r, p.readBuf[:8])
	value = int64(binary.BigEndian.Uint64(p.readBuf))
	return
}

func (p *binaryProtocol) ReadDouble(r io.Reader) (value float64, err error) {
	_, err = io.ReadFull(r, p.readBuf[:8])
	value = math.Float64frombits(binary.BigEndian.Uint64(p.readBuf))
	return
}

func (p *binaryProtocol) ReadString(r io.Reader) (string, error) {
	ln, err := p.ReadI32(r)
	if err != nil || ln == 0 {
		return "", err
	}
	if ln < 0 {
		return "", ProtocolError{"BinaryProtocol", "negative length while reading string"}
	}
	b := p.readBuf
	if int(ln) > len(b) {
		b = make([]byte, ln)
	} else {
		b = b[:ln]
	}
	if _, err := io.ReadFull(r, b); err != nil {
		return "", err
	}
	return string(b), nil
}

func (p *binaryProtocol) ReadBytes(r io.Reader) ([]byte, error) {
	ln, err := p.ReadI32(r)
	if err != nil || ln == 0 {
		return nil, err
	}
	if ln < 0 {
		return nil, ProtocolError{"BinaryProtocol", "negative length while reading bytes"}
	}
	b := make([]byte, ln)
	if _, err := io.ReadFull(r, b); err != nil {
		return nil, err
	}
	return b, nil
}
