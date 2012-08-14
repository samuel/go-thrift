package main

import (
	"fmt"
	"net"
	"net/rpc"
	"github.com/samuel/go-thrift"
)

type ResultCode int32

var (
	resultCodeOk       ResultCode = 0
	resultCodeTryLater ResultCode = 1
)

func (rc ResultCode) String() string {
	switch rc {
	case resultCodeOk:
		return "Ok"
	case resultCodeTryLater:
		return "TryLater"
	}
	return fmt.Sprintf("Unknown(%d)", rc)
}

type LogEntry struct {
	Category string `thrift:"1,required"`
	Message  string `thrift:"2,required"`
}

func (e *LogEntry) String() string {
	return fmt.Sprintf("%+v", *e)
}

type ScribeLogRequest struct {
	Messages []*LogEntry `thrift:"1,required"`
}

type ScribeLogResponse struct {
	Result ResultCode `thrift:"0,required"`
}

// type ScribeService interface {
// 	Log(*ScribeLogRequest) (ResultCode, error)
// }

type ScribeService int

func (s *ScribeService) Log(req *ScribeLogRequest, res *ScribeLogResponse) error {
	fmt.Printf("REQ: %+v\n", req)
	res.Result = resultCodeOk
	return nil
}

func main() {
	scribeService := new(ScribeService)
	rpc.RegisterName("Thrift", scribeService)

	ln, err := net.Listen("tcp", ":1463")
	if err != nil {
		panic(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("ERROR: %+v\n", err)
			continue
		}
		fmt.Printf("New connection %+v\n", conn)
		go rpc.ServeCodec(thrift.NewServerCodec(thrift.NewFramedReadWriteCloser(conn, 0), thrift.BinaryProtocol))
	}
}
