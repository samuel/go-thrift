package thrift

// import (
// 	"encoding/binary"
// 	"errors"
// 	"fmt"
// 	"io"
// 	// "math"
// )

// const (
// 	compactBufferSize      = 16
// 	compactProtocolId      = 0x82
// 	compactVersion         = 1
// 	compactVersionMask     = 0x1f
// 	compactTypeMask        = 0xe0
// 	compactTypeShiftAmount = 5
// )

// const (
// 	csClear = iota
// 	csFieldWrite
// 	csValueWrite
// 	csContainerWrite
// 	csBoolWrite
// 	csFieldRead
// 	csContainerRead
// 	csValueRead
// 	csBoolRead
// )

// const (
// 	ctStop = iota
// 	ctTrue
// 	ctFalse
// 	ctByte
// 	ctI16
// 	ctI32
// 	ctI64
// 	ctDouble
// 	ctBinary
// 	ctList
// 	ctSet
// 	ctMap
// 	ctStruct
// )

// var (
// 	ctypes = map[byte]int16{
// 		TypeStop:   ctStop,
// 		TypeBool:   ctTrue,
// 		TypeByte:   ctByte,
// 		TypeI16:    ctI16,
// 		TypeI32:    ctI32,
// 		TypeI64:    ctI64,
// 		TypeDouble: ctDouble,
// 		TypeString: ctBinary,
// 		TypeStruct: ctStruct,
// 		TypeList:   ctList,
// 		TypeSet:    ctSet,
// 		TypeMap:    ctMap,
// 	}
// )

// type InvalidStateError struct {
// 	CurrentState  int
// 	ExpectedState int
// }

// func (e InvalidStateError) Error() string {
// 	return fmt.Sprintf("InvalidStateError current=%d expected=%d",
// 		e.CurrentState, e.ExpectedState)
// }

// type structState struct {
// 	state   int
// 	lastFid int16
// }

// type compactProtocol struct {
// 	state     int
// 	lastFid   int16
// 	boolFid   int16
// 	boolValue int
// 	structs   []structState
// }

// func NewCompactProtocol() Protocol {
// 	return &compactProtocol{
// 		state:     csClear,
// 		lastFid:   0,
// 		boolFid:   -1,
// 		boolValue: -1,
// 		structs:   make([]structState, 0, 8),
// 		// self.__containers = []
// 	}
// }

// func (p *compactProtocol) writeVarint(w io.Writer, value int64) (err error) {
// 	b := make([]byte, compactBufferSize)
// 	n := binary.PutVarint(b, value)
// 	_, err = w.Write(b[:n])
// 	return
// }

// func (p *compactProtocol) WriteMessageBegin(w io.Writer, name string, messageType byte, seqid int32) (err error) {
// 	if p.state != csClear {
// 		return InvalidStateError{p.state, csClear}
// 	}

// 	if err = p.WriteByte(w, compactProtocolId); err != nil {
// 		return
// 	}
// 	if err = p.WriteByte(w, compactVersion|(messageType<<compactTypeShiftAmount)); err != nil {
// 		return
// 	}
// 	if err = p.writeVarint(w, int64(seqid)); err != nil {
// 		return
// 	}
// 	err = p.WriteString(w, name)

// 	p.state = csValueWrite
// 	return
// }

// func (p *compactProtocol) WriteMessageEnd(w io.Writer) error {
// 	if p.state != csValueWrite {
// 		return InvalidStateError{p.state, csValueWrite}
// 	}
// 	p.state = csClear
// 	return nil
// }

// func (p *compactProtocol) WriteStructBegin(w io.Writer, name string) error {
// 	switch p.state {
// 	case csClear:
// 	case csContainerWrite:
// 	case csValueWrite:
// 	default:
// 		return InvalidStateError{p.state, csClear}
// 	}
// 	p.structs = append(p.structs, structState{p.state, p.lastFid})
// 	p.state = csFieldWrite
// 	p.lastFid = 0
// 	return nil
// }

// func (p *compactProtocol) WriteStructEnd(w io.Writer) error {
// 	if p.state != csFieldWrite {
// 		return InvalidStateError{p.state, csFieldWrite}
// 	}
// 	if len(p.structs) == 0 {
// 		return errors.New("Struct end without matching begin")
// 	}
// 	st := p.structs[len(p.structs)-1]
// 	p.structs = p.structs[:len(p.structs)-1]
// 	p.state, p.lastFid = st.state, p.lastFid
// 	return nil
// }

// func (p *compactProtocol) writeFieldHeader(compactType int16, fid int16) error {
// 	delta := fid - p.lastFid
// 	if delta > 0 && delta <= 15 {
// 		p.writeUByte(delta<<4 | compactType)
// 	} else {
// 		p.WriteByte(compactType)
// 		p.WriteI16(fid)
// 	}
// 	p.lastFid = fid
// }

// func (p *compactProtocol) WriteFieldBegin(name string, fieldType byte, id int16) (err error) {
// 	if p.state != csFieldWrite {
// 		return InvalidStateError{p.state, csFieldWrite}
// 	}

// 	if fieldType == TypeBool {
// 		p.state = csBoolWrite
// 		p.boolFid = id
// 	} else {
// 		p.state = csValueWrite
// 		p.writeFieldHeader(ctypes[fieldType], id)
// 	}

// 	return nil
// }

// func (p *compactProtocol) WriteFieldEnd() error {
// 	if p.state != csValueWrite && p.state != csBoolWrite {
// 		return InvalidStateError{p.state, csValueWrite}
// 	}
// 	p.state = csFieldWrite
// 	return nil
// }

// func (p *compactProtocol) WriteFieldStop(w io.Writer) error {
// 	return p.WriteByte(w, TypeStop)
// }

// func (p *compactProtocol) writeCollectionBegin(w io.Writer, etype int, size int) error {
// 	if p.state != csValueWrite && p.state != csContainerWrite {
// 		return InvalidStateError{p.state, csValueWrite}
// 	}

// 	if size <= 14 {
// 		p.writeUByte(size<<4 | ctypes[etype])
// 	} else {
// 		p.writeUByte(0xf0 | ctypes[etype])
// 		p.writeSize(size)
// 	}
// }

// func (p *compactProtocol) writeMapBegin(w io.Writer, ktype int, vtype int, size int) error {
// 	if p.state != csValueWrite && p.state != csContainerWrite {
// 		return InvalidStateError{p.state, csValueWrite}
// 	}

// 	if size == 0 {
// 		p.writeByte(0)
// 	} else {
// 		p.writeSize(size)
// 		p.writeUByte(ctypes[ktype]<<4 | ctypes[vtype])
// 	}
// 	p.containers = append(p.containers, state)
// 	p.state = csContainerWrite
// }

// // func (p *compactProtocol) writeCollectionEnd(w io.Writer, )

// // def writeCollectionEnd(self):
// //   assert self.state == CONTAINER_WRITE, self.state
// //   self.state = self.__containers.pop()
// // writeMapEnd = writeCollectionEnd
// // writeSetEnd = writeCollectionEnd
// // writeListEnd = writeCollectionEnd

// // def __writeI16(self, i16):
// //   self.__writeVarint(makeZigZag(i16, 16))

// // def __writeSize(self, i32):
// //   self.__writeVarint(i32)

// // func (p *compactProtocol) WriteMapBegin(keyType byte, valueType byte, size int) (err error) {
// // 	if err = p.WriteByte(keyType); err != nil {
// // 		return
// // 	}
// // 	if err = p.WriteByte(valueType); err != nil {
// // 		return
// // 	}
// // 	return p.WriteI32(int32(size))
// // }

// // func (p *compactProtocol) WriteMapEnd() error {
// // 	return nil
// // }

// // func (p *compactProtocol) WriteListBegin(elementType byte, size int) (err error) {
// // 	if err = p.WriteByte(elementType); err != nil {
// // 		return
// // 	}
// // 	return p.WriteI32(int32(size))
// // }

// // func (p *compactProtocol) WriteListEnd() error {
// // 	return nil
// // }

// // func (p *compactProtocol) WriteSetBegin(elementType byte, size int) (err error) {
// // 	if err = p.WriteByte(elementType); err != nil {
// // 		return
// // 	}
// // 	return p.WriteI32(int32(size))
// // }

// // func (p *compactProtocol) WriteSetEnd() error {
// // 	return nil
// // }

// // func (p *compactProtocol) WriteBool(value bool) error {
// // 	if value {
// // 		return p.WriteByte(1)
// // 	}
// // 	return p.WriteByte(0)
// // }

// func (p *compactProtocol) WriteByte(w io.Writer, value byte) (err error) {
// 	b := []byte{value}
// 	_, err = w.Write(b[:1])
// 	return
// }

// // func (p *compactProtocol) WriteI16(value int16) (err error) {
// // 	b := p.buf[:2]
// // 	binary.BigEndian.PutUint16(b, uint16(value))
// // 	_, err = p.Writer.Write(b)
// // 	return
// // }

// // func (p *compactProtocol) WriteI32(value int32) (err error) {
// // 	b := p.buf[:4]
// // 	binary.BigEndian.PutUint32(b, uint32(value))
// // 	_, err = p.Writer.Write(b)
// // 	return
// // }

// // func (p *compactProtocol) WriteI64(value int64) (err error) {
// // 	b := p.buf[:8]
// // 	binary.BigEndian.PutUint64(b, uint64(value))
// // 	_, err = p.Writer.Write(b)
// // 	return
// // }

// // func (p *compactProtocol) WriteDouble(value float64) (err error) {
// // 	b := p.buf[:8]
// // 	binary.BigEndian.PutUint64(b, math.Float64bits(value))
// // 	_, err = p.Writer.Write(b)
// // 	return
// // }

// func (p *compactProtocol) WriteString(w io.Writer, value string) (err error) {
// 	if err = p.writeVarint(w, int64(len(value))); err != nil {
// 		return
// 	}
// 	_, err = w.Write([]byte(value))
// 	return
// }

// // func (p *compactProtocol) ReadMessageBegin() (name string, messageType byte, seqid int32, err error) {
// // 	size, e := p.ReadI32()
// // 	if e != nil {
// // 		err = e
// // 		return
// // 	}
// // 	if size < 0 {
// // 		version := uint32(size) & versionMask
// // 		if version != version1 {
// // 			err = ErrBadVersion
// // 			return
// // 		}
// // 		messageType = byte(uint32(size) & TypeMask)
// // 		if name, err = p.ReadString(); err != nil {
// // 			return
// // 		}
// // 	} else {
// // 		if p.StrictRead {
// // 			err = ErrNoProtocolVersionHeader
// // 			return
// // 		}
// // 		nameBytes := make([]byte, size)
// // 		if _, err = p.Reader.Read(nameBytes); err != nil {
// // 			return
// // 		}
// // 		name = string(nameBytes)
// // 		if messageType, err = p.ReadByte(); err != nil {
// // 			return
// // 		}
// // 	}
// // 	seqid, err = p.ReadI32()
// // 	return
// // }

// // func (p *compactProtocol) ReadMessageEnd() error {
// // 	return nil
// // }

// // func (p *compactProtocol) ReadStructBegin() error {
// // 	return nil
// // }

// // func (p *compactProtocol) ReadStructEnd() error {
// // 	return nil
// // }

// // func (p *compactProtocol) ReadFieldBegin() (fieldType byte, id int16, err error) {
// // 	if fieldType, err = p.ReadByte(); err != nil || fieldType == TypeStop {
// // 		return
// // 	}
// // 	id, err = p.ReadI16()
// // 	return
// // }

// // func (p *compactProtocol) ReadFieldEnd() error {
// // 	return nil
// // }

// // func (p *compactProtocol) ReadMapBegin() (keyType byte, valueType byte, size int, err error) {
// // 	if keyType, err = p.ReadByte(); err != nil {
// // 		return
// // 	}
// // 	if valueType, err = p.ReadByte(); err != nil {
// // 		return
// // 	}
// // 	var sz int32
// // 	sz, err = p.ReadI32()
// // 	size = int(sz)
// // 	return
// // }

// // func (p *compactProtocol) ReadMapEnd() error {
// // 	return nil
// // }

// // func (p *compactProtocol) ReadListBegin() (elementType byte, size int, err error) {
// // 	if elementType, err = p.ReadByte(); err != nil {
// // 		return
// // 	}
// // 	var sz int32
// // 	sz, err = p.ReadI32()
// // 	size = int(sz)
// // 	return
// // }

// // func (p *compactProtocol) ReadListEnd() error {
// // 	return nil
// // }

// // func (p *compactProtocol) ReadSetBegin() (elementType byte, size int, err error) {
// // 	if elementType, err = p.ReadByte(); err != nil {
// // 		return
// // 	}
// // 	var sz int32
// // 	sz, err = p.ReadI32()
// // 	size = int(sz)
// // 	return
// // }

// // func (p *compactProtocol) ReadSetEnd() error {
// // 	return nil
// // }

// // func (p *compactProtocol) ReadBool() (bool, error) {
// // 	if b, e := p.ReadByte(); e != nil {
// // 		return false, e
// // 	} else if b != 0 {
// // 		return true, nil
// // 	}
// // 	return false, nil
// // }

// // func (p *compactProtocol) ReadByte() (value byte, err error) {
// // 	_, err = io.ReadFull(p.Reader, p.buf[:1])
// // 	value = p.buf[0]
// // 	return
// // }

// // func (p *compactProtocol) ReadI16() (value int16, err error) {
// // 	b := p.buf[:2]
// // 	_, err = io.ReadFull(p.Reader, b)
// // 	value = int16(binary.BigEndian.Uint16(b))
// // 	return
// // }

// // func (p *compactProtocol) ReadI32() (value int32, err error) {
// // 	b := p.buf[:4]
// // 	_, err = io.ReadFull(p.Reader, b)
// // 	value = int32(binary.BigEndian.Uint32(b))
// // 	return
// // }

// // func (p *compactProtocol) ReadI64() (value int64, err error) {
// // 	b := p.buf[:8]
// // 	_, err = io.ReadFull(p.Reader, b)
// // 	value = int64(binary.BigEndian.Uint64(b))
// // 	return
// // }

// // func (p *compactProtocol) ReadDouble() (value float64, err error) {
// // 	b := p.buf[:8]
// // 	_, err = io.ReadFull(p.Reader, b)
// // 	value = math.Float64frombits(binary.BigEndian.Uint64(b))
// // 	return
// // }

// // func (p *compactProtocol) ReadString() (string, error) {
// // 	ln, err := p.ReadI32()
// // 	if err != nil || ln == 0 {
// // 		return "", err
// // 	}
// // 	var st []byte
// // 	if ln <= bufferSize {
// // 		st = p.buf[:ln]
// // 	} else {
// // 		st = make([]byte, ln)
// // 	}
// // 	if _, err := io.ReadFull(p.Reader, st); err != nil {
// // 		return "", err
// // 	}
// // 	return string(st), nil
// // }
