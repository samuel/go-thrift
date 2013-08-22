// Copyright 2012 Samuel Stauffer. All rights reserved.
// Use of this source code is governed by a 3-clause BSD
// license that can be found in the LICENSE file.

package thrift

import (
	"errors"
	"io"
	"net"
	"net/rpc"
)

// Implements rpc.ClientCodec
type clientCodec struct {
	transport      io.ReadWriteCloser
	protocol       Protocol
	onewayRequests chan pendingRequest
	twowayRequests chan pendingRequest
	enableOneway   bool
}

type pendingRequest struct {
	method string
	seq    uint64
}

type oneway interface {
	Oneway() bool
}

var (
	ErrTooManyPendingRequests = errors.New("thrift.client: too many pending requests")
	ErrOnewayNotEnabled       = errors.New("thrift.client: one way support not enabled on codec")
)

const maxPendingRequests = 1000

// Dial connects to a Thrift RPC server at the specified network address using the specified protocol.
func Dial(network, address string, framed bool, protocol Protocol, supportOnewayRequests bool) (*rpc.Client, error) {
	conn, err := net.Dial(network, address)
	if err != nil {
		return nil, err
	}
	codec := &clientCodec{
		transport: conn,
		protocol:  protocol,
	}
	if supportOnewayRequests {
		codec.enableOneway = true
		codec.onewayRequests = make(chan pendingRequest, maxPendingRequests)
		codec.twowayRequests = make(chan pendingRequest, maxPendingRequests)
	}
	if framed {
		codec.transport = NewFramedReadWriteCloser(conn, DefaultMaxFrameSize)
	}
	return rpc.NewClientWithCodec(codec), nil
}

// NewClient returns a new rpc.Client to handle requests to the set of
// services at the other end of the connection.
func NewClient(conn io.ReadWriteCloser, protocol Protocol, supportOnewayRequests bool) *rpc.Client {
	return rpc.NewClientWithCodec(NewClientCodec(conn, protocol, supportOnewayRequests))
}

// NewClientCodec returns a new rpc.ClientCodec using Thrift RPC on conn using the specified protocol.
func NewClientCodec(conn io.ReadWriteCloser, protocol Protocol, supportOnewayRequests bool) rpc.ClientCodec {
	c := &clientCodec{
		transport: conn,
		protocol:  protocol,
	}
	if supportOnewayRequests {
		c.enableOneway = true
		c.onewayRequests = make(chan pendingRequest, maxPendingRequests)
		c.twowayRequests = make(chan pendingRequest, maxPendingRequests)
	}
	return c
}

func (c *clientCodec) WriteRequest(request *rpc.Request, thriftStruct interface{}) error {
	if err := c.protocol.WriteMessageBegin(c.transport, request.ServiceMethod, MessageTypeCall, int32(request.Seq)); err != nil {
		return err
	}
	if err := EncodeStruct(c.transport, c.protocol, thriftStruct); err != nil {
		return err
	}
	if err := c.protocol.WriteMessageEnd(c.transport); err != nil {
		return err
	}
	var err error
	if flusher, ok := c.transport.(Flusher); ok {
		err = flusher.Flush()
	}
	if err == nil {
		ow := false
		if o, ok := thriftStruct.(oneway); ok {
			ow = o.Oneway()
		}
		if c.enableOneway {
			if ow {
				select {
				case c.onewayRequests <- pendingRequest{request.ServiceMethod, request.Seq}:
				default:
					err = ErrTooManyPendingRequests
				}
			} else {
				select {
				case c.twowayRequests <- pendingRequest{request.ServiceMethod, request.Seq}:
				default:
					err = ErrTooManyPendingRequests
				}
			}
		} else if ow {
			return ErrOnewayNotEnabled
		}
	}
	return err
}

func (c *clientCodec) ReadResponseHeader(response *rpc.Response) error {
	if c.enableOneway {
		select {
		case ow := <-c.onewayRequests:
			response.ServiceMethod = ow.method
			response.Seq = ow.seq
			return nil
		case _ = <-c.twowayRequests:
		}
	}

	name, messageType, seq, err := c.protocol.ReadMessageBegin(c.transport)
	if err != nil {
		return err
	}
	response.ServiceMethod = name
	response.Seq = uint64(seq)
	if messageType == MessageTypeException {
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
	return c.transport.Close()
}
