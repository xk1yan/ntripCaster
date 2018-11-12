package main

import (
	"fmt"
	"net"
	"time"
)

func clientrun() {
	var testdataClient = "GET /yunzhihui HTTP/1.0\r\nUser-Agent: NTRIP sss\r\nAuthorization: Basic eXpocm92ZXIwMzoxMjM0NTY=\r\n"
	i := 0
	j := 0
	data := make([]byte, 1024)
	for {
		j = 0
		time.Sleep(time.Second * 2)
		fmt.Printf("client connect times:%d\r\n", i)
		conn, err := net.Dial("tcp", "127.0.0.1:2101")
		if err != nil {
			fmt.Println(err)
		}
		_ = conn.SetDeadline(time.Now().Add(time.Second * 4))
		if _, err := conn.Write([]byte(testdataClient)); err != nil {
			fmt.Printf("write login:%s", err)
			conn.Close()
			continue
		}
		lenn, err1 := conn.Read(data)
		if err != nil {
			fmt.Printf("login E:%s len:%d", err1, lenn)
			conn.Close()
			continue
		}
		go func() {
			// conn.Close()
			// conn.Close()
		}()
		fmt.Println("lenn:", lenn, err1)
		fmt.Println(len(data))
		if string(data[:lenn]) == backStr[1] {
			for {
				data := make([]byte, 1024)
				_ = conn.SetDeadline(time.Now().Add(time.Second * 100))
				lenn, err1 := conn.Read(data)
				if err1 != nil || lenn == 0 {
					fmt.Printf("read E:%s len:%d", err1, lenn)
					conn.Close()
					break
				}
				fmt.Println("client read data:", data[:lenn])
				if _, err := conn.Write([]byte("client data")); err != nil {
					fmt.Printf("write data error :%s", err)
					conn.Close()
					break
				}
				j++
				if j >= 2 {
					conn.Close()
					break
				}
			}
		} else {
			fmt.Println("client fail var error")
		}
		i++
	}

}
