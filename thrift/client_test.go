// Copyright 2012 Samuel Stauffer. All rights reserved.
// Use of this source code is governed by a 3-clause BSD
// license that can be found in the LICENSE file.

package thrift

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/rpc"
	// "runtime"
	"sync"
	"testing"
)

var (
	serverAddr, newServerAddr string
	once, newOnce             sync.Once
)

type TestRequest struct {
	Value int32 `thrift:"1,required"`
}

type TestResponse struct {
	Value int32 `thrift:"0,required"`
}

type TestService int

func (s *TestService) Success(req *TestRequest, res *TestResponse) error {
	res.Value = req.Value
	return nil
}

func (s *TestService) Fail(req *TestRequest, res *TestResponse) error {
	res.Value = req.Value
	return errors.New("fail")
}

func listenTCP() (net.Listener, string) {
	l, e := net.Listen("tcp", "127.0.0.1:0") // any available address
	if e != nil {
		log.Fatalf("net.Listen tcp :0: %v", e)
	}
	return l, l.Addr().String()
}

func startServer() {
	rpc.RegisterName("Thrift", new(TestService))

	var l net.Listener
	l, serverAddr = listenTCP()
	log.Println("Test RPC server listening on", serverAddr)
	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				panic(err)
				continue
			}
			go rpc.ServeCodec(NewServerCodec(NewFramedReadWriteCloser(conn, 0), NewBinaryProtocol(true, false)))
		}
	}()
}

func TestRPCClientSuccess(t *testing.T) {
	once.Do(startServer)

	c, err := Dial("tcp", serverAddr, true, NewBinaryProtocol(true, false), false)
	if err != nil {
		t.Fatalf("NewClient returned error: %+v", err)
	}
	req := &TestRequest{123}
	res := &TestResponse{789}
	if err := c.Call("Success", req, res); err != nil {
		t.Fatalf("Client.Call returned error: %+v", err)
	}
	if res.Value != req.Value {
		t.Fatalf("Response value wrong: %d != %d", res.Value, req.Value)
	}
}

func TestRPCClientFail(t *testing.T) {
	once.Do(startServer)

	c, err := Dial("tcp", serverAddr, true, NewBinaryProtocol(true, false), false)
	if err != nil {
		t.Fatalf("NewClient returned error: %+v", err)
	}
	req := &TestRequest{123}
	res := &TestResponse{789}
	if err := c.Call("Fail", req, res); err == nil {
		t.Fatalf("Client.Call didn't return an error as it should")
	} else if err.Error() != "Internal Error: fail" {
		t.Fatalf("Expected 'fail' error but got '%s'", err)
	}

	// Make sure an exception doesn't cause future requests to fail

	if err := c.Call("Success", req, res); err != nil {
		t.Fatalf("Client.Call returned error: %+v", err)
	}
	if res.Value != req.Value {
		t.Fatalf("Response value wrong: %d != %d", res.Value, req.Value)
	}
}

func TestRPCMallocCount(t *testing.T) {
	once.Do(startServer)

	c, err := Dial("tcp", serverAddr, true, NewBinaryProtocol(true, false), false)
	if err != nil {
		t.Fatalf("NewClient returned error: %+v", err)
	}
	req := &TestRequest{123}
	res := &TestResponse{789}
	allocs := testing.AllocsPerRun(100, func() {
		if err := c.Call("Success", req, res); err != nil {
			t.Fatalf("Client.Call returned error: %+v", err)
		}
	})
	fmt.Printf("mallocs per thrift.rpc round trip: %d\n", int(allocs))
}
