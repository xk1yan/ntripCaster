package main

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"net"
	"sync"
	"time"
)

var logger = logs.NewLogger(1000)
var ntripMountpoints = make(map[string]*ntripMountpointsNode)

func init() {
	logger.SetLogger(logs.AdapterMultiFile, `{"filename":"logs/test.log","separate":["emergency", "alert", "critical", "error", "warning", "notice", "info", "debug"]}`)
	logger.EnableFuncCallDepth(true)

	fmt.Println("init")
	loginType = new(usersIn)
	loginType.updateUserMap()
	// fmt.Println(loginType)
}

type ntripClientsNode struct {
	name    string
	conTime time.Time
	con     net.Conn
	dataCh  chan *[]byte
}

type ntripMountpointsNode struct {
	name        string
	con         net.Conn
	nodeRMMutex sync.RWMutex
	clients     map[string]*ntripClientsNode
}

var loginType userIfa

var nMsRMMutex sync.RWMutex

func handleConnection(conn net.Conn) {
	_ = conn.SetDeadline(time.Now().Add(time.Second * 4))
	data := make([]byte, 1024)
	lenn, err := conn.Read(data)
	if err != nil {
		logger.Warning("login E:%s", err)
		conn.Close()
		return
	}
	if lenn == 0 {
		conn.Close()
	}

	data = data[:lenn]
	// fmt.Println(time.Now())
	res := verifyLogin(loginType, data)
	// fmt.Println(time.Now())
	_, err = conn.Write([]byte(backStr[res.backStrIndex]))
	switch {
	case res.userType == mountpointType:
		node := ntripMountpointsNode{}
		node.name = res.mountPointName
		node.con = conn
		node.clients = make(map[string]*ntripClientsNode)

		nMsRMMutex.Lock()
		if _, ok := ntripMountpoints[res.mountPointName]; ok {
			conn.Close()
		} else {
			ntripMountpoints[res.mountPointName] = &node
		}
		nMsRMMutex.Unlock()
		mountPointRun(&node, res.mountPointName)
	case res.userType == clientType:
		cNode := ntripClientsNode{}
		cNode.con = conn
		cNode.name = res.clientName
		cNode.conTime = time.Now()
		cNode.dataCh = make(chan *[]byte, 3) //为两种nil形式关闭通道留空间,防止写通道阻塞

		nMsRMMutex.Lock()
		if _, ok := ntripMountpoints[res.mountPointName]; ok {
			//会覆盖旧的 覆盖前关通道，关连接
			if _, okk := ntripMountpoints[res.mountPointName].clients[res.clientName]; okk {
				ntripMountpoints[res.mountPointName].clients[res.clientName].dataCh <- nil
				ntripMountpoints[res.mountPointName].clients[res.clientName].con.Close()

			}
			ntripMountpoints[res.mountPointName].clients[res.clientName] = &cNode

		}
		nMsRMMutex.Unlock()
		clientRun(&cNode, res.mountPointName)
	default:
		conn.Close()
	}
}

func mountPointRun(mNode *ntripMountpointsNode, mountPointName string) {
	defer func() {
		nMsRMMutex.Lock()
		delete(ntripMountpoints, mountPointName)
		nMsRMMutex.Unlock()
		for _, v := range mNode.clients {
			v.con.Close()
			// close(v.dataCh)
			v.dataCh <- nil // 关闭通道
			// delete(mNode.clients, k)
		}

	}()
	for {
		data := make([]byte, 1024)
		_ = mNode.con.SetDeadline(time.Now().Add(time.Second * 20))
		lenn, err := mNode.con.Read(data)
		if err != nil || lenn == 0 {
			break
		} else {
			data = data[:lenn]
			for _, v := range mNode.clients {
				if len(v.dataCh) < 1 { //只要最新数据
					v.dataCh <- &data
				}
			}
		}
	}
}

func clientRun(cNode *ntripClientsNode, mountPointName string) {
	sendDone, readDone := make(chan struct{}), make(chan struct{})
	defer func() {
		nMsRMMutex.Lock()
		if _, ok := ntripMountpoints[mountPointName]; ok {
			if ntripMountpoints[mountPointName].clients[cNode.name] != cNode { //如果没有被覆盖
				delete(ntripMountpoints[mountPointName].clients, cNode.name)
			}

		}
		nMsRMMutex.Unlock()
	}()
	go func() { //sendData
		defer func() {
			cNode.con.Close()
			close(sendDone)
		}()

		for {
			data := <-cNode.dataCh
			if data == nil { //如果通道被关闭
				break
			}
			_ = cNode.con.SetWriteDeadline(time.Now().Add(time.Second * 3))
			_, err := cNode.con.Write(*data)
			if err != nil {
				break
			}
		}
	}()

	go func() { //readData GGA
		defer func() {
			cNode.con.Close()
			close(readDone)
		}()

		for {
			rdata := make([]byte, 1024)
			timeoutSet := time.Now().Add(time.Second * 60)
			_ = cNode.con.SetReadDeadline(timeoutSet)
			lenn, err := cNode.con.Read(rdata)
			if err != nil && time.Now().Before(timeoutSet) { //读数据忽略超时错误
				break
			}
			if lenn == 0 {
				continue
			}
			rdata = rdata[:lenn]
		}
	}()

	<-sendDone
	<-readDone
}

func loop() {
	ln, err := net.Listen("tcp", ":2101")
	if err != nil {
		// handle error
		logger.Emergency("listen Fail: %s", err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
			continue
		}
		go handleConnection(conn)
	}

}

func main() {
	loop()

}
