package thrift

import (
	"errors"
	"io"
	"net/rpc"
)

type serverCodec struct {
	transport io.ReadWriteCloser
	protocol  Protocol
}

// ServeConn runs the Thrift RPC server on a single connection. ServeConn blocks,
// serving the connection until the client hangs up. The caller typically invokes
// ServeConn in a go statement.
func ServeConn(conn io.ReadWriteCloser, protocol Protocol) {
	rpc.ServeCodec(NewServerCodec(conn, protocol))
}

// NewServerCodec returns a new rpc.ServerCodec using Thrift RPC on conn using the specified protocol.
func NewServerCodec(conn io.ReadWriteCloser, protocol Protocol) rpc.ServerCodec {
	return &serverCodec{conn, protocol}
}

func (c *serverCodec) ReadRequestHeader(request *rpc.Request) error {
	name, messageType, seq, err := c.protocol.ReadMessageBegin(c.transport)
	if err != nil {
		return err
	}
	request.ServiceMethod = name
	request.Seq = uint64(seq)

	if messageType != messageTypeCall { // Currently don't support one way
		return errors.New("Exception Call message type")
	}

	if err := skipValue(c.transport, c.protocol, typeStruct); err != nil {
		return err
	}

	if err := c.protocol.ReadMessageEnd(c.transport); err != nil {
		return err
	}

	// exc := &ApplicationException{}
	// x = TApplicationException(TApplicationException.UNKNOWN_METHOD, 'Unknown function %s' % (name))
	// oprot.writeMessageBegin(name, TMessageType.EXCEPTION, seqid)
	// x.write(oprot)
	// oprot.writeMessageEnd()
	// oprot.trans.flush()

	return nil
}

func (c *serverCodec) ReadRequestBody(thriftStruct interface{}) error {
	if thriftStruct == nil {
		panic("TODO: Skip body in ReadRequestBody")
	}

	if err := DecodeStruct(c.transport, c.protocol, thriftStruct); err != nil {
		return err
	}
	return c.protocol.ReadMessageEnd(c.transport)
}

func (c *serverCodec) WriteResponse(response *rpc.Response, thriftStruct interface{}) error {
	if err := c.protocol.WriteMessageBegin(c.transport, response.ServiceMethod, messageTypeReply, int32(response.Seq)); err != nil {
		return err
	}
	if response.Error != "" {
		thriftStruct = &ApplicationException{response.Error, ExceptionUnknown}
	}
	if err := EncodeStruct(c.transport, c.protocol, thriftStruct); err != nil {
		return err
	}
	if err := c.protocol.WriteMessageEnd(c.transport); err != nil {
		return err
	}
	if flusher, ok := c.transport.(Flusher); ok {
		return flusher.Flush()
	}
	return nil
}

func (c *serverCodec) Close() error {
	return nil
}
