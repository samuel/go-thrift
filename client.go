// Copyright 2012 Samuel Stauffer. All rights reserved.
// Use of this source code is governed by a 3-clause BSD
// license that can be found in the LICENSE file.

package thrift

import (
	"io"
	"net"
	"net/rpc"
)

// Implements rpc.ClientCodec
type clientCodec struct {
	transport io.ReadWriteCloser
	protocol  Protocol
}

// Dial connects to a Thrift RPC server at the specified network address using the specified protocol.
func Dial(network, address string, framed bool, protocol Protocol) (*rpc.Client, error) {
	conn, err := net.Dial(network, address)
	if err != nil {
		return nil, err
	}
	codec := &clientCodec{
		transport: conn,
		protocol:  protocol,
	}
	if framed {
		codec.transport = NewFramedReadWriteCloser(conn, DefaultMaxFrameSize)
	}
	return rpc.NewClientWithCodec(codec), nil
}

// NewClient returns a new rpc.Client to handle requests to the set of
// services at the other end of the connection.
func NewClient(conn io.ReadWriteCloser, protocol Protocol) *rpc.Client {
	return rpc.NewClientWithCodec(NewClientCodec(conn, protocol))
}

// NewClientCodec returns a new rpc.ClientCodec using Thrift RPC on conn using the specified protocol.
func NewClientCodec(conn io.ReadWriteCloser, protocol Protocol) rpc.ClientCodec {
	return &clientCodec{
		transport: conn,
		protocol:  protocol,
	}
}

func (c *clientCodec) WriteRequest(request *rpc.Request, thriftStruct interface{}) error {
	if err := c.protocol.WriteMessageBegin(c.transport, request.ServiceMethod, messageTypeCall, int32(request.Seq)); err != nil {
		return err
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

func (c *clientCodec) ReadResponseHeader(response *rpc.Response) error {
	name, messageType, seq, err := c.protocol.ReadMessageBegin(c.transport)
	if err != nil {
		return err
	}
	response.ServiceMethod = name
	response.Seq = uint64(seq)
	if messageType == messageTypeException {
		exception := &ApplicationException{}
		if err := DecodeStruct(c.transport, c.protocol, exception); err != nil {
			return err
		}
		response.Error = exception.String()
		return c.protocol.ReadMessageEnd(c.transport)
	}
	return nil
}

func (c *clientCodec) ReadResponseBody(thriftStruct interface{}) error {
	if thriftStruct == nil {
		// Should only get called if ReadResponseHeader set the Error value in
		// which case we've already read the body (ApplicationException)
		return nil
	}

	if err := DecodeStruct(c.transport, c.protocol, thriftStruct); err != nil {
		return err
	}

	return c.protocol.ReadMessageEnd(c.transport)
}

func (c *clientCodec) Close() error {
	return nil
}
