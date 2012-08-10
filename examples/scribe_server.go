package main

import (
	"net"
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

type ScribeLogRequest struct {
	Messages []*LogEntry `thrift:"1,required"`
}

type ScribeLogResponse struct {
	Result ResultCode `thrift:"0,required"`
}

type ScribeService interface {
	Log(*ScribeLogRequest) (ResultCode, error)
}

type ScribeService int

func (s *ScribeService) Log(req *ScribeLogRequest, res *ScribeLogResponse) error {
	req := &ScribeLogRequest{messages}
	res := &ScribeLogResponse{}
	err := s.Client.Call("Log", req, res)
	return res.Result, err
}

func main() {
	scribeService := new(ScribeService)
	rpc.Register(scribeService)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":1234")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
}
