package main

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"net"
	"time"
	"xk/ntripCaster/testDataMaker"
)

var backStr = "ICY 200 OK\r\n"
var logger = logs.NewLogger(1000)

func init() {
	logger.SetLogger(logs.AdapterMultiFile, `{"filename":"logs/ntripServer.log","separate":["emergency", "alert", "critical", "error", "warning", "notice", "info", "debug"]}`)
	logger.EnableFuncCallDepth(true)
}

func serverrun(mountPoint, passwd, agent, nSTR string) {
	// var testdataServer = "SOURCE yunzhihui yzhbase01\r\nSource-Agent: NTRIP sss\r\nSTR: bbbb\r\n"
	i := 0
	loginData := fmt.Sprintf("SOURCE %s %s\r\nSource-Agent: %s\r\nSTR: %s\r\n", passwd, mountPoint, agent, nSTR)
	data := make([]byte, 1024)
	time.Sleep(time.Second * 1)
	fmt.Printf("server %s connect \r\n", mountPoint)
	conn, err := net.Dial("tcp", "127.0.0.1:2101")
	if err != nil {
		fmt.Println(err)
		return
	}
	_ = conn.SetDeadline(time.Now().Add(time.Second * 4))
	if _, err := conn.Write([]byte(loginData)); err != nil {
		fmt.Printf("write login:%s", err)
		conn.Close()
		return
	}
	lenn, err1 := conn.Read(data)
	if err != nil || lenn == 0 {
		fmt.Printf("login E:%s len:%d", err1, lenn)
		conn.Close()
		return
	}
	if string(data[:lenn]) == backStr {
		for {
			_ = conn.SetDeadline(time.Now().Add(time.Second * 4))
			sdata := testDataMaker.DataMaker(i)
			if _, err := conn.Write([]byte(sdata)); err != nil {
				fmt.Printf("m:%s,E:", mountPoint, err)
				logger.Error("m:%s,E:", mountPoint, err)
				conn.Close()
				break
			}
			time.Sleep(time.Second * 1)
			i++
		}
		fmt.Printf("m:%s  end", mountPoint)
	}
}

func main() {
	mps := testDataMaker.CreateMountPointAndClient(100, 100)
	for k, _ := range mps {
		// mountpoint := k
		go func(mountpoint string) {
			fmt.Println(mountpoint)
			serverrun(mountpoint, "123456", "NTRIP", "11")

		}(k)

	}
	time.Sleep(time.Second * 1200)
}
