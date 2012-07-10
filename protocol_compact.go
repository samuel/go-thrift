package thrift

import (
	"encoding/binary"
	"io"
	// "math"
)

const (
	compactBufferSize      = 16
	compactProtocolId      = 0x82
	compactVersion         = 1
	compactVersionMask     = 0x1f
	compactTypeMask        = 0xe0
	compactTypeShiftAmount = 5
)

type CompactProtocol struct {
	Writer io.Writer
	Reader io.Reader
	buf    [compactBufferSize]byte
}

func (p *CompactProtocol) writeByte(value byte) (err error) {
	p.buf[0] = value
	_, err = p.Writer.Write(p.buf[:1])
	return
}

func (p *CompactProtocol) writeVarint(value int64) (err error) {
	n := binary.PutVarint(p.buf[:compactBufferSize], value)
	_, err = p.Writer.Write(p.buf[:n])
	return
}

func (p *CompactProtocol) WriteMessageBegin(name string, messageType byte, seqid int32) (err error) {
	if err = p.writeByte(compactProtocolId); err != nil {
		return
	}
	if err = p.writeByte(compactVersion | (messageType << compactTypeShiftAmount)); err != nil {
		return
	}
	if err = p.writeVarint(int64(seqid)); err != nil {
		return
	}
	err = p.WriteString(name)
	return
}

func (p *CompactProtocol) WriteMessageEnd() error {
	return nil
}

func (p *CompactProtocol) WriteStructBegin(name string) error {
	return nil
}

func (p *CompactProtocol) WriteStructEnd() error {
	return nil
}

// func (p *CompactProtocol) WriteFieldBegin(name string, fieldType byte, id int16) (err error) {
// 	if err = p.WriteByte(fieldType); err != nil {
// 		return
// 	}
// 	return p.WriteI16(id)
// }

// func (p *CompactProtocol) WriteFieldEnd() error {
// 	return nil
// }

// func (p *CompactProtocol) WriteFieldStop() error {
// 	return p.WriteByte(typeStop)
// }

// func (p *CompactProtocol) WriteMapBegin(keyType byte, valueType byte, size int) (err error) {
// 	if err = p.WriteByte(keyType); err != nil {
// 		return
// 	}
// 	if err = p.WriteByte(valueType); err != nil {
// 		return
// 	}
// 	return p.WriteI32(int32(size))
// }

// func (p *CompactProtocol) WriteMapEnd() error {
// 	return nil
// }

// func (p *CompactProtocol) WriteListBegin(elementType byte, size int) (err error) {
// 	if err = p.WriteByte(elementType); err != nil {
// 		return
// 	}
// 	return p.WriteI32(int32(size))
// }

// func (p *CompactProtocol) WriteListEnd() error {
// 	return nil
// }

// func (p *CompactProtocol) WriteSetBegin(elementType byte, size int) (err error) {
// 	if err = p.WriteByte(elementType); err != nil {
// 		return
// 	}
// 	return p.WriteI32(int32(size))
// }

// func (p *CompactProtocol) WriteSetEnd() error {
// 	return nil
// }

// func (p *CompactProtocol) WriteBool(value bool) error {
// 	if value {
// 		return p.WriteByte(1)
// 	}
// 	return p.WriteByte(0)
// }

// func (p *CompactProtocol) WriteByte(value byte) (err error) {
// 	p.buf[0] = value
// 	_, err = p.Writer.Write(p.buf[:1])
// 	return
// }

// func (p *CompactProtocol) WriteI16(value int16) (err error) {
// 	b := p.buf[:2]
// 	binary.BigEndian.PutUint16(b, uint16(value))
// 	_, err = p.Writer.Write(b)
// 	return
// }

// func (p *CompactProtocol) WriteI32(value int32) (err error) {
// 	b := p.buf[:4]
// 	binary.BigEndian.PutUint32(b, uint32(value))
// 	_, err = p.Writer.Write(b)
// 	return
// }

// func (p *CompactProtocol) WriteI64(value int64) (err error) {
// 	b := p.buf[:8]
// 	binary.BigEndian.PutUint64(b, uint64(value))
// 	_, err = p.Writer.Write(b)
// 	return
// }

// func (p *CompactProtocol) WriteDouble(value float64) (err error) {
// 	b := p.buf[:8]
// 	binary.BigEndian.PutUint64(b, math.Float64bits(value))
// 	_, err = p.Writer.Write(b)
// 	return
// }

func (p *CompactProtocol) WriteString(value string) (err error) {
	if err = p.writeVarint(int64(len(value))); err != nil {
		return
	}
	_, err = p.Writer.Write([]byte(value))
	return
}

// func (p *CompactProtocol) ReadMessageBegin() (name string, messageType byte, seqid int32, err error) {
// 	size, e := p.ReadI32()
// 	if e != nil {
// 		err = e
// 		return
// 	}
// 	if size < 0 {
// 		version := uint32(size) & versionMask
// 		if version != version1 {
// 			err = ErrBadVersion
// 			return
// 		}
// 		messageType = byte(uint32(size) & typeMask)
// 		if name, err = p.ReadString(); err != nil {
// 			return
// 		}
// 	} else {
// 		if p.StrictRead {
// 			err = ErrNoProtocolVersionHeader
// 			return
// 		}
// 		nameBytes := make([]byte, size)
// 		if _, err = p.Reader.Read(nameBytes); err != nil {
// 			return
// 		}
// 		name = string(nameBytes)
// 		if messageType, err = p.ReadByte(); err != nil {
// 			return
// 		}
// 	}
// 	seqid, err = p.ReadI32()
// 	return
// }

// func (p *CompactProtocol) ReadMessageEnd() error {
// 	return nil
// }

// func (p *CompactProtocol) ReadStructBegin() error {
// 	return nil
// }

// func (p *CompactProtocol) ReadStructEnd() error {
// 	return nil
// }

// func (p *CompactProtocol) ReadFieldBegin() (fieldType byte, id int16, err error) {
// 	if fieldType, err = p.ReadByte(); err != nil || fieldType == typeStop {
// 		return
// 	}
// 	id, err = p.ReadI16()
// 	return
// }

// func (p *CompactProtocol) ReadFieldEnd() error {
// 	return nil
// }

// func (p *CompactProtocol) ReadMapBegin() (keyType byte, valueType byte, size int, err error) {
// 	if keyType, err = p.ReadByte(); err != nil {
// 		return
// 	}
// 	if valueType, err = p.ReadByte(); err != nil {
// 		return
// 	}
// 	var sz int32
// 	sz, err = p.ReadI32()
// 	size = int(sz)
// 	return
// }

// func (p *CompactProtocol) ReadMapEnd() error {
// 	return nil
// }

// func (p *CompactProtocol) ReadListBegin() (elementType byte, size int, err error) {
// 	if elementType, err = p.ReadByte(); err != nil {
// 		return
// 	}
// 	var sz int32
// 	sz, err = p.ReadI32()
// 	size = int(sz)
// 	return
// }

// func (p *CompactProtocol) ReadListEnd() error {
// 	return nil
// }

// func (p *CompactProtocol) ReadSetBegin() (elementType byte, size int, err error) {
// 	if elementType, err = p.ReadByte(); err != nil {
// 		return
// 	}
// 	var sz int32
// 	sz, err = p.ReadI32()
// 	size = int(sz)
// 	return
// }

// func (p *CompactProtocol) ReadSetEnd() error {
// 	return nil
// }

// func (p *CompactProtocol) ReadBool() (bool, error) {
// 	if b, e := p.ReadByte(); e != nil {
// 		return false, e
// 	} else if b != 0 {
// 		return true, nil
// 	}
// 	return false, nil
// }

// func (p *CompactProtocol) ReadByte() (value byte, err error) {
// 	_, err = io.ReadFull(p.Reader, p.buf[:1])
// 	value = p.buf[0]
// 	return
// }

// func (p *CompactProtocol) ReadI16() (value int16, err error) {
// 	b := p.buf[:2]
// 	_, err = io.ReadFull(p.Reader, b)
// 	value = int16(binary.BigEndian.Uint16(b))
// 	return
// }

// func (p *CompactProtocol) ReadI32() (value int32, err error) {
// 	b := p.buf[:4]
// 	_, err = io.ReadFull(p.Reader, b)
// 	value = int32(binary.BigEndian.Uint32(b))
// 	return
// }

// func (p *CompactProtocol) ReadI64() (value int64, err error) {
// 	b := p.buf[:8]
// 	_, err = io.ReadFull(p.Reader, b)
// 	value = int64(binary.BigEndian.Uint64(b))
// 	return
// }

// func (p *CompactProtocol) ReadDouble() (value float64, err error) {
// 	b := p.buf[:8]
// 	_, err = io.ReadFull(p.Reader, b)
// 	value = math.Float64frombits(binary.BigEndian.Uint64(b))
// 	return
// }

// func (p *CompactProtocol) ReadString() (string, error) {
// 	ln, err := p.ReadI32()
// 	if err != nil || ln == 0 {
// 		return "", err
// 	}
// 	var st []byte
// 	if ln <= bufferSize {
// 		st = p.buf[:ln]
// 	} else {
// 		st = make([]byte, ln)
// 	}
// 	if _, err := io.ReadFull(p.Reader, st); err != nil {
// 		return "", err
// 	}
// 	return string(st), nil
// }
