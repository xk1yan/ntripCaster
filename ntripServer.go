package main

import (
	"fmt"
	"net"
	"time"
)

func serverrun() {
	var testdataServer = "SOURCE yunzhihui yzhbase01\r\nSource-Agent: NTRIP sss\r\nSTR: bbbb\r\n"
	i := 0
	j := 0
	data := make([]byte, 1024)
	for {
		j = 0
		time.Sleep(time.Second * 1)
		fmt.Printf("server connect times:%d\r\n", i)
		conn, err := net.Dial("tcp", "127.0.0.1:2101")
		if err != nil {
			fmt.Println(err)
		}
		_ = conn.SetDeadline(time.Now().Add(time.Second * 4))
		if _, err := conn.Write([]byte(testdataServer)); err != nil {
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
				_ = conn.SetDeadline(time.Now().Add(time.Second * 4))
				if _, err := conn.Write([]byte("22e1er2r")); err != nil {
					fmt.Printf("write data:%s", err)
					conn.Close()
					break
				}
				time.Sleep(time.Second * 1)
				j++
				if j >= 10 {
					conn.Close()
					break
				}
			}
		}
		i++

	}

}
