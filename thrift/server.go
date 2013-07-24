// Copyright 2012 Samuel Stauffer. All rights reserved.
// Use of this source code is governed by a 3-clause BSD
// license that can be found in the LICENSE file.

package thrift

import (
	"errors"
	"io"
	"net/rpc"
	"strings"
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
	name = CamelCase(name)
	if strings.ContainsRune(name, '.') {
		request.ServiceMethod = name
	} else {
		request.ServiceMethod = "Thrift." + name
	}
	request.Seq = uint64(seq)

	if messageType != messageTypeCall { // Currently don't support one way
		return errors.New("thrift: exception Call message type")
	}

	return nil
}

func (c *serverCodec) ReadRequestBody(thriftStruct interface{}) error {
	if thriftStruct == nil {
		if err := SkipValue(c.transport, c.protocol, TypeStruct); err != nil {
			return err
		}
	} else {
		if err := DecodeStruct(c.transport, c.protocol, thriftStruct); err != nil {
			return err
		}
	}
	return c.protocol.ReadMessageEnd(c.transport)
}

func (c *serverCodec) WriteResponse(response *rpc.Response, thriftStruct interface{}) error {
	mtype := byte(messageTypeReply)
	if response.Error != "" {
		mtype = messageTypeException
		etype := int32(ExceptionInternalError)
		if strings.HasPrefix(response.Error, "rpc: can't find") {
			etype = ExceptionUnknownMethod
		}
		thriftStruct = &ApplicationException{response.Error, etype}
	}
	if err := c.protocol.WriteMessageBegin(c.transport, response.ServiceMethod, mtype, int32(response.Seq)); err != nil {
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

func (c *serverCodec) Close() error {
	return c.transport.Close()
}
