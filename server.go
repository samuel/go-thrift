package thrift

import (
	"errors"
	"net/rpc"
)

type ServerCodec struct {
	protocol Protocol
}

func NewServerCodec(protocol Protocol) rpc.ServerCodec {
	return &ServerCodec{protocol}
}

func (c *ServerCodec) ReadRequestHeader(request *rpc.Request) error {
	name, messageType, seq, err := c.protocol.ReadMessageBegin()
	if err != nil {
		return err
	}
	request.ServiceMethod = name
	request.Seq = uint64(seq)

	if messageType != messageTypeCall { // Currently don't support one way
		return errors.New("Exception Call message type")
	}

	// iprot.skip(TType.STRUCT)
	// iprot.readMessageEnd()
	// x = TApplicationException(TApplicationException.UNKNOWN_METHOD, 'Unknown function %s' % (name))
	// oprot.writeMessageBegin(name, TMessageType.EXCEPTION, seqid)
	// x.write(oprot)
	// oprot.writeMessageEnd()
	// oprot.trans.flush()

	return nil
}

func (c *ServerCodec) ReadRequestBody(thriftStruct interface{}) error {
	if thriftStruct == nil {
		panic("TODO: Skip body in ReadRequestBody")
	}

	if err := DecodeStruct(c.protocol, thriftStruct); err != nil {
		return err
	}
	return c.protocol.ReadMessageEnd()
}

func (c *ServerCodec) WriteResponse(response *rpc.Response, thriftStruct interface{}) error {
	if err := c.protocol.WriteMessageBegin(response.ServiceMethod, messageTypeReply, int32(response.Seq)); err != nil {
		return err
	}
	if response.Error != "" {
		thriftStruct = &ApplicationException{response.Error, ExceptionUnknown}
	}
	if err := EncodeStruct(c.protocol, thriftStruct); err != nil {
		return err
	}
	return c.protocol.WriteMessageEnd()
}

func (c *ServerCodec) Close() error {
	return nil
}
