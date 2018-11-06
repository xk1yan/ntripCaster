package main

import (
	"encoding/base64"
	"fmt"
	"testing"
	"time"
)

var testdataServer = "SOURCE yunzhihui yzhbase01\r\nSource-Agent: NTRIP sss\r\nSTR: bbbb\r\n"
var testdataClient = "GET /jzsk2base2 HTTP/1.0\r\nUser-Agent: NTRIP sss\r\nAuthorization: Basic eXpocm92ZXIwMzoxMjM0NTY=\r\n"

func TestVerifyLogin(t *testing.T) {
	var v verifyer
	v.loginType = new(testlogin)
	v.loginType.updateUserMap()
	dataServer := []byte(testdataClient) //类型转化
	userType, mountpointName, clientName := v.verifyLogin(dataServer)
	fmt.Println("登录", userType, mountpointName, clientName)
	fmt.Println(time.Now())

}
func TestLogin(t *testing.T) {
	// var users loginIfa
	// fmt.Printf("%v\r\n", users)
	// users = new(testlogin)
	// fmt.Printf("%v\r\n", users)
	// _ = users.updateUserMap()
	// fmt.Printf("%v\r\n", users)
	// fmt.Println(users)

}
func TestTimes(t *testing.T) {
	s1 := "User-Agent: "
	s2 := "STR: "
	fmt.Println("s1 len:", len(s1))
	fmt.Println("s2 len:", len(s2))
	fmt.Println(s1[0:12])
	fmt.Println(s1[1:])
	if OK := 2 > 4; OK {
		fmt.Println("1111")
		fmt.Println("2222")
	}
	res, err := base64.StdEncoding.DecodeString("eXpocm92ZXIwMzoxMjM0NTY=")
	fmt.Println(string(res), err)
	// fmt.Println(time.Now())
	// fmt.Println(time.Now().Add(time.Second * 3))
}
