package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"net"
	"os"
	"strconv"
	"time"
	"xk/ntripCaster/testDataMaker"
)

var logger = logs.NewLogger(1000)

func init() {
	logger.SetLogger(logs.AdapterMultiFile, `{"filename":"logs/ntripClient.log","separate":["emergency", "alert", "critical", "error", "warning", "notice", "info", "debug"]}`)
	logger.EnableFuncCallDepth(true)
}

var backStr = "ICY 200 OK\r\n"

func clientrun(mountPoint, username, passwd string) {
	n := -1
	userAndPasswd := fmt.Sprintf("%s:%s", username, passwd)
	userBase64 := base64.StdEncoding.EncodeToString([]byte(userAndPasswd))
	loginData := fmt.Sprintf("GET /%s HTTP/1.0\r\nUser-Agent: NTRIP sss\r\nAuthorization: Basic %s\r\n\r\n", mountPoint, userBase64)
	data := make([]byte, 1024)
	fmt.Printf("client %s connect mountPoint %s\r\n", username, mountPoint)
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
		fmt.Printf("login E:%s len:%d\r\n", err1, lenn)
		logger.Alert("LOGIN m:%s,c:%s,E:%s", mountPoint, username, err)
		conn.Close()
		return
	}
	if string(data[:lenn]) == backStr {
		go func() {
			time.Sleep(time.Second * 20)
			if _, err := conn.Write([]byte("GGA data")); err != nil {
				fmt.Printf("write data error :%s", err)
				conn.Close()
			}
		}()
		for {
			data := make([]byte, 1024)
			_ = conn.SetDeadline(time.Now().Add(time.Second * 10))
			lenn, err1 := conn.Read(data)
			if err1 != nil || lenn == 0 {
				fmt.Printf("read E:%s len:%d\r\n", err1, lenn)
				logger.Error("m:%s,c:%s,E:%s", mountPoint, username, err1)
				conn.Close()
				break
			}
			data = data[:lenn]
			rn, r := varData(data)
			if !r {
				logger.Warning("m:%s,c:%s,Incomplete data", mountPoint, username)
				n = -1
				continue
			}
			if n == -1 {
				n = rn
			} else if n+1 != rn {
				logger.Warning("m:%s,c:%s,Data sequence discontinuity", mountPoint, username)

			}
			n = rn
			// fmt.Println("client read datalen:", len(data[:lenn]))
		}
	} else {
		fmt.Println("client fail var error")
	}

}

func varData(data []byte) (int, bool) {
	b := make(map[string]string)
	num := 0
	err := json.Unmarshal(data, &b)
	// fmt.Println(b, err)
	if err != nil || len(b) != 1 {
		return -1, false
	}
	for k, v := range b {
		n, err := strconv.Atoi(k)
		if err != nil {
			return -1, false
		}
		res := testDataMaker.KeyMaker(n)
		// fmt.Println(res)
		if res != v {
			return -1, false
		}
		num = n
		break
	}
	return num, true
}

func testUserInsert() {
	mps := testDataMaker.CreateMountPointAndClient(100, 50)
	n := 0
	fmt.Println(os.Args)
	fmt.Println(os.Args[1])
	mm := os.Args[1]

	// time.Sleep(time.Second * 1000)
	for k, v := range mps {
		if k != mm {
			continue
		}
		for _, c := range v {
			n++
			if n%10 == 0 {
				time.Sleep(time.Millisecond * 500)
			}
			go func(mountpoint, client string) {
				fmt.Println(mountpoint, client)
				clientrun(mountpoint, client, "123456")

			}(k, c)
		}

	}
	time.Sleep(time.Second * 1000)

}

func testUserMysql() {
	for {
		clientrun("DMDP_F1", "car001", "654321")
	}

}
func main() {

}
