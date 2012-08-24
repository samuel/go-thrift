package thrift

import (
	"encoding/binary"
	"errors"
	"io"
	"math"
)

var (
	ErrBadVersion              = errors.New("Bad version in ReadMessageBegin")
	ErrNoProtocolVersionHeader = errors.New("No protocol version header")

	// BinaryProtocol with strictWrite=true and strictRead=false
	BinaryProtocol = NewBinaryProtocol(true, false)
)

type binaryProtocol struct {
	strictWrite bool
	strictRead  bool
}

func NewBinaryProtocol(strictWrite bool, strictRead bool) Protocol {
	return &binaryProtocol{strictWrite, strictRead}
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
	return p.WriteByte(w, typeStop)
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
	_, err := w.Write([]byte{value})
	return err
}

func (p *binaryProtocol) WriteI16(w io.Writer, value int16) (err error) {
	b := []byte{0, 0}
	binary.BigEndian.PutUint16(b, uint16(value))
	_, err = w.Write(b)
	return
}

func (p *binaryProtocol) WriteI32(w io.Writer, value int32) (err error) {
	b := []byte{0, 0, 0, 0}
	binary.BigEndian.PutUint32(b, uint32(value))
	_, err = w.Write(b)
	return
}

func (p *binaryProtocol) WriteI64(w io.Writer, value int64) (err error) {
	b := []byte{0, 0, 0, 0, 0, 0, 0, 0}
	binary.BigEndian.PutUint64(b, uint64(value))
	_, err = w.Write(b)
	return
}

func (p *binaryProtocol) WriteDouble(w io.Writer, value float64) (err error) {
	b := []byte{0, 0, 0, 0, 0, 0, 0, 0}
	binary.BigEndian.PutUint64(b, math.Float64bits(value))
	_, err = w.Write(b)
	return
}

func (p *binaryProtocol) WriteString(w io.Writer, value string) error {
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
			err = ErrBadVersion
			return
		}
		messageType = byte(uint32(size) & typeMask)
		if name, err = p.ReadString(r); err != nil {
			return
		}
	} else {
		if p.strictRead {
			err = ErrNoProtocolVersionHeader
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
	if fieldType, err = p.ReadByte(r); err != nil || fieldType == typeStop {
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
	b := []byte{0}
	_, err = io.ReadFull(r, b[:1])
	value = b[0]
	return
}

func (p *binaryProtocol) ReadI16(r io.Reader) (value int16, err error) {
	b := []byte{0, 0}
	_, err = io.ReadFull(r, b)
	value = int16(binary.BigEndian.Uint16(b))
	return
}

func (p *binaryProtocol) ReadI32(r io.Reader) (value int32, err error) {
	b := []byte{0, 0, 0, 0}
	_, err = io.ReadFull(r, b)
	value = int32(binary.BigEndian.Uint32(b))
	return
}

func (p *binaryProtocol) ReadI64(r io.Reader) (value int64, err error) {
	b := []byte{0, 0, 0, 0, 0, 0, 0, 0}
	_, err = io.ReadFull(r, b)
	value = int64(binary.BigEndian.Uint64(b))
	return
}

func (p *binaryProtocol) ReadDouble(r io.Reader) (value float64, err error) {
	b := []byte{0, 0, 0, 0, 0, 0, 0, 0}
	_, err = io.ReadFull(r, b)
	value = math.Float64frombits(binary.BigEndian.Uint64(b))
	return
}

func (p *binaryProtocol) ReadString(r io.Reader) (string, error) {
	bytes, err := p.ReadBytes(r)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (p *binaryProtocol) ReadBytes(r io.Reader) ([]byte, error) {
	ln, err := p.ReadI32(r)
	if err != nil || ln == 0 {
		return nil, err
	}
	st := make([]byte, ln)
	if _, err := io.ReadFull(r, st); err != nil {
		return nil, err
	}
	return st, nil
}
