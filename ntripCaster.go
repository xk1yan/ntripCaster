package main

import (
	"github.com/astaxie/beego/logs"
	"net"
)

var logger = logs.NewLogger(1000)

func init() {
	logger.SetLogger(logs.AdapterMultiFile, `{"filename":"logs/test.log","separate":["emergency", "alert", "critical", "error", "warning", "notice", "info", "debug"]}`)
	logger.EnableFuncCallDepth(true)

}

type ntripClient struct {
	name string
	con  *net.Conn
}

type ntripMountpoint struct {
	name    string
	con     *net.Conn
	clients []*ntripClient
}

func main() {
	var verobj verifyer
	ln, err := net.Listen("tcp", ":2101")
	if err != nil {
		// handle error
		logger.Emergency("listen Fail %s", err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
			continue
		}
		go verobj.handleConnection(conn)
	}

}
