package thrift

import (
	"net"
	"net/rpc"
)

// Implements rpc.ClientCodec
type ClientCodec struct {
	protocol Protocol
}

func Dial(network, address string) (*rpc.Client, error) {
	conn, err := net.Dial(network, address)
	if err != nil {
		return nil, err
	}
	protocol := &BinaryProtocol{Writer: conn, Reader: conn, StrictWrite: true, StrictRead: false}
	codec := &ClientCodec{
		// netConn:  conn,
		protocol: protocol,
	}
	return rpc.NewClientWithCodec(codec), nil
}

func (c *ClientCodec) WriteRequest(request *rpc.Request, thriftStruct interface{}) error {
	if err := c.protocol.WriteMessageBegin(request.ServiceMethod, messageTypeCall, int32(request.Seq)); err != nil {
		return err
	}
	if err := EncodeStruct(c.protocol, thriftStruct); err != nil {
		return err
	}
	return c.protocol.WriteMessageEnd()
}

func (c *ClientCodec) ReadResponseHeader(response *rpc.Response) error {
	name, messageType, seq, err := c.protocol.ReadMessageBegin()
	if err != nil {
		return err
	}
	response.ServiceMethod = name
	response.Seq = uint64(seq)
	if messageType == messageTypeException {
		exception := &ApplicationException{}
		if err := DecodeStruct(c.protocol, exception); err != nil {
			return err
		}
		response.Error = exception.String()
		return c.protocol.ReadMessageEnd()
	}
	return nil
}

func (c *ClientCodec) ReadResponseBody(thriftStruct interface{}) error {
	if thriftStruct == nil {
		return nil
	}

	if err := DecodeStruct(c.protocol, thriftStruct); err != nil {
		return err
	}

	return c.protocol.ReadMessageEnd()
}

func (c *ClientCodec) Close() error {
	return nil
}
