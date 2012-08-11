package thrift

import (
	"fmt"
	"net/rpc"
	"testing"
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

type ScribeClient struct {
	Client *rpc.Client
}

func (s *ScribeClient) Log(messages []*LogEntry) (ResultCode, error) {
	req := &ScribeLogRequest{messages}
	res := &ScribeLogResponse{}
	err := s.Client.Call("Log", req, res)
	return res.Result, err
}

func TestClient(t *testing.T) {
	c, err := Dial("tcp", "localhost:1463", true, DefaultBinaryProtocol)
	if err != nil {
		t.Fatalf("NewClient returned error: %+v", err)
	}
	scribe := ScribeClient{c}
	rc, err := scribe.Log([]*LogEntry{&LogEntry{"category", "message"}})
	if err != nil {
		t.Fatalf("scribe.Log returned error: %+v", err)
	}
	fmt.Printf("%+v\n", rc)
	// req := &ScribeLogRequest{[]*LogEntry{&LogEntry{"category", "message"}}}
	// res := &ScribeLogResponse{123}
	// if err := c.Call("Log", req, res); err != nil {
	// 	t.Fatalf("Client.Call returned error: %+v", err)
	// }
	// fmt.Printf("%+v\n", res)
}
