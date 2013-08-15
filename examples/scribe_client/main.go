package main

import (
	"fmt"
	"net"

	"github.com/shxsun/go-thrift/examples/scribe"
	"github.com/shxsun/go-thrift/thrift"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:1463")
	if err != nil {
		panic(err)
	}

	client := thrift.NewClient(thrift.NewFramedReadWriteCloser(conn, 0), thrift.NewBinaryProtocol(true, false), false)
	scr := scribe.ScribeClient{client}
	res, err := scr.Log([]*scribe.LogEntry{{"category", "message"}})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Response: %+v\n", res)
}
