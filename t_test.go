package main

import (
	"encoding/base64"
	// "encoding/json"
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"
)

var testdataServer = "SOURCE yunzhihui yzhbase01\r\nSource-Agent: NTRIP sss\r\nSTR: bbbb\r\n"
var testdataClient = "GET /jzsk2base2 HTTP/1.0\r\nUser-Agent: NTRIP sss\r\nAuthorization: Basic eXpocm92ZXIwMzoxMjM0NTY=\r\n"

func TestVerifyLogin(t *testing.T) {
	t.Skip("skip VerifyLogin")
	loginType := new(usersIn)
	loginType.updateUserMap()
	dataServer := []byte(testdataClient) //类型转化
	res := verifyLogin(loginType, dataServer)
	fmt.Println("登录", res.userType, res.mountPointName, res.clientName, res.backStrIndex, backStr[res.backStrIndex])
	fmt.Println(time.Now())

}
func TestCaster(t *testing.T) {
	t.Skip("skip Caster")
	t.Parallel()
	loop()
}

func TestServer(t *testing.T) {
	t.Skip("skip Server")
	t.Parallel()
	// serverrun()
	fmt.Println(time.Now().Add(time.Second * 3))

}
func TestClient(t *testing.T) {
	t.Skip("skip Client")
	t.Parallel()
	// clientrun("1", "2", "3")

}
func TestSlice(t *testing.T) {
	t.Skip("skip ")
	var rwlk sync.RWMutex
	rwlk.Lock()
	rwlk.Unlock()
	var a = make(map[string]*ntripMountpointsNode)
	node := ntripMountpointsNode{}
	node.name = "w1"
	node.clients = make(map[string]*ntripClientsNode)

	node2 := ntripMountpointsNode{}
	node2.name = "w2"
	node2.clients = make(map[string]*ntripClientsNode)
	// var lk sync.RWMutex
	// node.nodeRMMutex = lk
	// fmt.Printf("%p\r\n", node)

	a["w1"] = &node
	fmt.Println(a["w1"])
	n1 := a["w1"]
	if n1 == a["w1"] {
		fmt.Println("eeeeee", n1)
	}
	a["w1"] = &node2
	if n1 == a["w1"] {
		fmt.Println("22222222", n1)
	}
	fmt.Println(a["w1"])

	fmt.Println("........")
	for k, v := range a {
		fmt.Println(k, v)
	}
	// nMnode, ok := a["w1"]
	// nMnode.name = "w3"
	// fmt.Println(a, ok)
	// fmt.Println(nMnode, ok)
}
func TestTimes(t *testing.T) {
	t.Skip("time ")
	ch := make(chan *[]byte, 1)
	// data := make([]byte, 1024)
	fmt.Printf("ch len:%d\r\n", len(ch))
	fmt.Println("00")
	ch <- nil
	ch <- nil
	fmt.Printf("ch len:%d\r\n", len(ch))
	fmt.Println("11")
	close(ch)
	// ch <- nil
	fmt.Printf("ch len:%d\r\n", len(ch))
	fmt.Println("22")
	a := <-ch
	fmt.Printf("%T,%v\r\n", a, a)
	b := <-ch
	fmt.Printf("%T,%v\r\n", b, b)
}
func TestBase64(t *testing.T) {
	t.Skip(" ")
	ss := "abc1231l:qqeqeqwe123"
	// res, err := base64.StdEncoding.DecodeString(ds2s[2])
	res := base64.StdEncoding.EncodeToString([]byte(ss))
	fmt.Println(res)
}
func TestDataMaker(t *testing.T) {
	// t.Skip(" ")

	// a := createMountPointAndClient(2, 2)
	// fmt.Println(a)
	// data, _ := json.Marshal(a)
	// fmt.Println(string(data))
	// b := make(map[string][]string)
	// json.Unmarshal(data, &b)
	// fmt.Println(b)
	// c := dataMaker(100)
	// fmt.Println(string(c))
	// fmt.Println(len(c))
	// s := "9223384124222222299"
	s := "9223372036854775807"
	b, err := strconv.Atoi(s)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%T,%d\r\n", b, b)
	d := make(map[string]string)
	d["12"] = "qeqweqwe"
	for k, v := range d {
		fmt.Println(k, v)
	}
	fmt.Println(len(d))
}
