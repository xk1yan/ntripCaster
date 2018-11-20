package main

import (
	"fmt"
	"github.com/Unknwon/goconfig"
	"github.com/astaxie/beego/logs"
	"net"
	"sync"
	"time"
)

const (
	up   = "y"
	down = "n"
)

var logger = logs.NewLogger(1000)
var ntripMountpoints = make(map[string]*ntripMountpointsNode)

var configs = map[string]map[string]string{}

var loginType userIfa
var stateType statePutIfa

var nMsRMMutex sync.RWMutex

var dataPool = &sync.Pool{New: func() interface{} { return make([]byte, 1024) }}

type userIfa interface {
	mountpointLogin(mountpointName, passwd string) (res bool, err error)
	clientLogin(mountpointName, clientName, passwd string) (res bool, err error)
	usersInit(conf map[string]string) (err error)
}

type statePutIfa interface {
	setRoverSt(roverName, st string)
	setMounpointSt(mountpointName, st string)
	setAllDown()
	initState(conf map[string]string)
}

func readConfig() {
	conf, err := goconfig.LoadConfigFile("ntripCaster.conf")
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, v := range conf.GetSectionList() {
		v1, _ := conf.GetSection(v)
		configs[v] = v1
	}
}

func init() {
	var conf map[string]string

	logger.SetLogger(logs.AdapterMultiFile, `{"filename":"logs/test.log","separate":["emergency", "alert", "critical", "error", "warning", "notice", "info", "debug"]}`)
	logger.EnableFuncCallDepth(true)

	fmt.Println("init")
	readConfig()
	switch configs["ntripcaster"]["userstype"] {
	case "insert":
		loginType = new(usersIn)
		conf = configs["insert"]
	case "mysql":
		loginType = new(usersMysql)
		conf = configs["mysql"]
	default:
		panic("userstype Error")
	}
	loginType.usersInit(conf)

	switch configs["ntripcaster"]["statetpye"] {
	case "mysql":
		stateType = new(statePutMysql)
		stateType.initState(configs["mysql"])
	default:
		stateType = new(statePutNone)
	}
}

type ntripClientsNode struct {
	name    string
	conTime time.Time
	con     *net.Conn
	dataCh  chan *[]byte
}

type ntripMountpointsNode struct {
	name        string
	con         *net.Conn
	nodeRMMutex sync.RWMutex
	clients     map[string]*ntripClientsNode
}

func handleConnection(conn *net.Conn) {
	returnFlag := false
	_ = (*conn).SetDeadline(time.Now().Add(time.Second * 4))
	data := dataPool.Get().([]byte)
	lenn, err := (*conn).Read(data)
	if err != nil {
		logger.Warning("login E:%s", err)
		(*conn).Close()
		return
	}
	if lenn == 0 {
		(*conn).Close()
		return
	}

	dataa := data[:lenn]
	fmt.Println(time.Now())
	res := verifyLogin(loginType, dataa)
	dataPool.Put(data)
	// fmt.Println(time.Now())
	_, err = (*conn).Write([]byte(backStr[res.backStrIndex]))
	switch {
	case res.userType == mountpointType:
		node := ntripMountpointsNode{}
		node.name = res.mountPointName
		node.con = conn
		node.clients = make(map[string]*ntripClientsNode)

		nMsRMMutex.Lock()
		if _, ok := ntripMountpoints[res.mountPointName]; ok {
			(*conn).Close()
			returnFlag = true
		} else {
			ntripMountpoints[res.mountPointName] = &node
		}
		nMsRMMutex.Unlock()
		if returnFlag {
			return
		}
		stateType.setMounpointSt(res.mountPointName, up)
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
				(*ntripMountpoints[res.mountPointName].clients[res.clientName].con).Close()

			}
			ntripMountpoints[res.mountPointName].clients[res.clientName] = &cNode

		} else {
			(*conn).Close()
			returnFlag = true
		}
		nMsRMMutex.Unlock()
		if returnFlag {
			return
		}
		stateType.setRoverSt(res.clientName, up)
		clientRun(&cNode, res.mountPointName)
	default:
		(*conn).Close()
	}
}

func mountPointRun(mNode *ntripMountpointsNode, mountPointName string) {
	defer func() {
		nMsRMMutex.Lock()
		delete(ntripMountpoints, mountPointName)
		nMsRMMutex.Unlock()
		for k, _ := range mNode.clients {
			(*mNode.clients[k].con).Close()
			// (*v.con).Close()
			mNode.clients[k].dataCh <- nil // 关闭通道
			// v.dataCh <- nil
		}
		stateType.setMounpointSt(mountPointName, down)

	}()
	for {
		data := dataPool.Get().([]byte)
		_ = (*mNode.con).SetDeadline(time.Now().Add(time.Second * 20))
		lenn, err := (*mNode.con).Read(data)
		if err != nil || lenn == 0 {
			break
		} else {
			dataa := data[:lenn]
			for k, _ := range mNode.clients {
				if len(mNode.clients[k].dataCh) < 1 { //只要最新数据
					mNode.clients[k].dataCh <- &dataa
				}
			}
		}
		dataPool.Put(data)
	}
}

func clientRun(cNode *ntripClientsNode, mountPointName string) {
	sendDone, readDone := make(chan struct{}), make(chan struct{})
	defer func() {
		delNodeF := false //下线标志，使设置下线状态在锁外进行。保证锁时间可控
		nMsRMMutex.Lock()
		if _, ok := ntripMountpoints[mountPointName]; ok {
			if ntripMountpoints[mountPointName].clients[cNode.name] == cNode { //如果没有被覆盖
				delete(ntripMountpoints[mountPointName].clients, cNode.name)
				delNodeF = true
			}

		} else { //基准站已经释放
			delNodeF = true
		}
		nMsRMMutex.Unlock()
		if delNodeF {
			stateType.setRoverSt(cNode.name, down)
		}

	}()
	go func() { //sendData
		defer func() {
			(*cNode.con).Close()
			close(sendDone)
		}()

		for {
			data := <-cNode.dataCh
			if data == nil { //如果通道被关闭
				break
			}
			_ = (*cNode.con).SetWriteDeadline(time.Now().Add(time.Second * 3))
			_, err := (*cNode.con).Write(*data)
			if err != nil {
				break
			}
		}
	}()

	go func() { //readData GGA
		defer func() {
			(*cNode.con).Close()
			close(readDone)
		}()

		for {
			// rdata := make([]byte, 1024)
			rdata := dataPool.Get().([]byte)
			timeoutSet := time.Now().Add(time.Second * 60)
			_ = (*cNode.con).SetReadDeadline(timeoutSet)
			lenn, err := (*cNode.con).Read(rdata)
			if err != nil && time.Now().Before(timeoutSet) { //读数据忽略超时错误
				break
			}
			if lenn == 0 {
				continue
			}
			// rdataa := rdata[:lenn]
			dataPool.Put(rdata)
		}
	}()

	<-sendDone
	<-readDone
}

func loop() {
	defer func() {
		stateType.setAllDown()
	}()
	port := fmt.Sprintf(":%s", configs["ntripcaster"]["port"])
	// ln, err := net.Listen("tcp", ":2101")
	ln, err := net.Listen("tcp", port)
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
		go handleConnection(&conn)
	}

}

func main() {
	loop()

}
