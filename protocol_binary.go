package thrift

import (
	"encoding/binary"
	"errors"
	"io"
	"math"
)

const (
	bufferSize = 8
)

var (
	ErrBadVersion              = errors.New("Bad version in ReadMessageBegin")
	ErrNoProtocolVersionHeader = errors.New("No protocol version header")
)

type BinaryProtocol struct {
	Writer      io.Writer
	Reader      io.Reader
	StrictWrite bool
	StrictRead  bool
	buf         [bufferSize]byte
}

func (p *BinaryProtocol) WriteMessageBegin(name string, messageType byte, seqid int32) (err error) {
	if p.StrictWrite {
		if err = p.WriteI32(int32(version1 | uint32(messageType))); err != nil {
			return
		}
		if err = p.WriteString(name); err != nil {
			return
		}
		err = p.WriteI32(seqid)
	} else {
		if err = p.WriteString(name); err != nil {
			return
		}
		if err = p.WriteByte(messageType); err != nil {
			return
		}
		err = p.WriteI32(seqid)
	}
	return
}

func (p *BinaryProtocol) WriteMessageEnd() error {
	return nil
}

func (p *BinaryProtocol) WriteStructBegin(name string) error {
	return nil
}

func (p *BinaryProtocol) WriteStructEnd() error {
	return nil
}

func (p *BinaryProtocol) WriteFieldBegin(name string, fieldType byte, id int16) (err error) {
	if err = p.WriteByte(fieldType); err != nil {
		return
	}
	return p.WriteI16(id)
}

func (p *BinaryProtocol) WriteFieldEnd() error {
	return nil
}

func (p *BinaryProtocol) WriteFieldStop() error {
	return p.WriteByte(typeStop)
}

func (p *BinaryProtocol) WriteMapBegin(keyType byte, valueType byte, size int) (err error) {
	if err = p.WriteByte(keyType); err != nil {
		return
	}
	if err = p.WriteByte(valueType); err != nil {
		return
	}
	return p.WriteI32(int32(size))
}

func (p *BinaryProtocol) WriteMapEnd() error {
	return nil
}

func (p *BinaryProtocol) WriteListBegin(elementType byte, size int) (err error) {
	if err = p.WriteByte(elementType); err != nil {
		return
	}
	return p.WriteI32(int32(size))
}

func (p *BinaryProtocol) WriteListEnd() error {
	return nil
}

func (p *BinaryProtocol) WriteSetBegin(elementType byte, size int) (err error) {
	if err = p.WriteByte(elementType); err != nil {
		return
	}
	return p.WriteI32(int32(size))
}

func (p *BinaryProtocol) WriteSetEnd() error {
	return nil
}

func (p *BinaryProtocol) WriteBool(value bool) error {
	if value {
		return p.WriteByte(1)
	}
	return p.WriteByte(0)
}

func (p *BinaryProtocol) WriteByte(value byte) (err error) {
	p.buf[0] = value
	_, err = p.Writer.Write(p.buf[:1])
	return
}

func (p *BinaryProtocol) WriteI16(value int16) (err error) {
	b := p.buf[:2]
	binary.BigEndian.PutUint16(b, uint16(value))
	_, err = p.Writer.Write(b)
	return
}

func (p *BinaryProtocol) WriteI32(value int32) (err error) {
	b := p.buf[:4]
	binary.BigEndian.PutUint32(b, uint32(value))
	_, err = p.Writer.Write(b)
	return
}

func (p *BinaryProtocol) WriteI64(value int64) (err error) {
	b := p.buf[:8]
	binary.BigEndian.PutUint64(b, uint64(value))
	_, err = p.Writer.Write(b)
	return
}

func (p *BinaryProtocol) WriteDouble(value float64) (err error) {
	b := p.buf[:8]
	binary.BigEndian.PutUint64(b, math.Float64bits(value))
	_, err = p.Writer.Write(b)
	return
}

func (p *BinaryProtocol) WriteString(value string) (err error) {
	if err = p.WriteI32(int32(len(value))); err != nil {
		return
	}
	_, err = p.Writer.Write([]byte(value))
	return
}

func (p *BinaryProtocol) ReadMessageBegin() (name string, messageType byte, seqid int32, err error) {
	size, e := p.ReadI32()
	if e != nil {
		err = e
		return
	}
	if size < 0 {
		version := uint32(size) & versionMask
		if version != version1 {
			err = ErrBadVersion
			return
		}
		messageType = byte(uint32(size) & typeMask)
		if name, err = p.ReadString(); err != nil {
			return
		}
	} else {
		if p.StrictRead {
			err = ErrNoProtocolVersionHeader
			return
		}
		nameBytes := make([]byte, size)
		if _, err = p.Reader.Read(nameBytes); err != nil {
			return
		}
		name = string(nameBytes)
		if messageType, err = p.ReadByte(); err != nil {
			return
		}
	}
	seqid, err = p.ReadI32()
	return
}

func (p *BinaryProtocol) ReadMessageEnd() error {
	return nil
}

func (p *BinaryProtocol) ReadStructBegin() error {
	return nil
}

func (p *BinaryProtocol) ReadStructEnd() error {
	return nil
}

func (p *BinaryProtocol) ReadFieldBegin() (fieldType byte, id int16, err error) {
	if fieldType, err = p.ReadByte(); err != nil || fieldType == typeStop {
		return
	}
	id, err = p.ReadI16()
	return
}

func (p *BinaryProtocol) ReadFieldEnd() error {
	return nil
}

func (p *BinaryProtocol) ReadMapBegin() (keyType byte, valueType byte, size int, err error) {
	if keyType, err = p.ReadByte(); err != nil {
		return
	}
	if valueType, err = p.ReadByte(); err != nil {
		return
	}
	var sz int32
	sz, err = p.ReadI32()
	size = int(sz)
	return
}

func (p *BinaryProtocol) ReadMapEnd() error {
	return nil
}

func (p *BinaryProtocol) ReadListBegin() (elementType byte, size int, err error) {
	if elementType, err = p.ReadByte(); err != nil {
		return
	}
	var sz int32
	sz, err = p.ReadI32()
	size = int(sz)
	return
}

func (p *BinaryProtocol) ReadListEnd() error {
	return nil
}

func (p *BinaryProtocol) ReadSetBegin() (elementType byte, size int, err error) {
	if elementType, err = p.ReadByte(); err != nil {
		return
	}
	var sz int32
	sz, err = p.ReadI32()
	size = int(sz)
	return
}

func (p *BinaryProtocol) ReadSetEnd() error {
	return nil
}

func (p *BinaryProtocol) ReadBool() (bool, error) {
	if b, e := p.ReadByte(); e != nil {
		return false, e
	} else if b != 0 {
		return true, nil
	}
	return false, nil
}

func (p *BinaryProtocol) ReadByte() (value byte, err error) {
	_, err = io.ReadFull(p.Reader, p.buf[:1])
	value = p.buf[0]
	return
}

func (p *BinaryProtocol) ReadI16() (value int16, err error) {
	b := p.buf[:2]
	_, err = io.ReadFull(p.Reader, b)
	value = int16(binary.BigEndian.Uint16(b))
	return
}

func (p *BinaryProtocol) ReadI32() (value int32, err error) {
	b := p.buf[:4]
	_, err = io.ReadFull(p.Reader, b)
	value = int32(binary.BigEndian.Uint32(b))
	return
}

func (p *BinaryProtocol) ReadI64() (value int64, err error) {
	b := p.buf[:8]
	_, err = io.ReadFull(p.Reader, b)
	value = int64(binary.BigEndian.Uint64(b))
	return
}

func (p *BinaryProtocol) ReadDouble() (value float64, err error) {
	b := p.buf[:8]
	_, err = io.ReadFull(p.Reader, b)
	value = math.Float64frombits(binary.BigEndian.Uint64(b))
	return
}

func (p *BinaryProtocol) ReadString() (string, error) {
	ln, err := p.ReadI32()
	if err != nil || ln == 0 {
		return "", err
	}
	var st []byte
	if ln <= bufferSize {
		st = p.buf[:ln]
	} else {
		st = make([]byte, ln)
	}
	if _, err := io.ReadFull(p.Reader, st); err != nil {
		return "", err
	}
	return string(st), nil
}
