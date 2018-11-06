package main

import (
	"encoding/base64"
	"fmt"
	"net"
	"strings"
	"time"
)

const (
	MountpointType = "mountpoint"
	ClientType     = "client"
	NoType         = "no"
)

type loginIfa interface {
	mountpointLogin(mountpointName, passwd string) (res bool, err error)
	clientLogin(mountpointName, clientName, passwd string) (res bool, err error)
	updateUserMap() (err error)
}
type verifyer struct {
	loginType loginIfa
}

func (v verifyer) handleConnection(conn net.Conn) {
	err := conn.SetReadDeadline(time.Now().Add(time.Second * 60))
	if err != nil {
		logger.Error("set time con timeout error :%s", err)
	}
	data := make([]byte, 1024)
	lenn, err := conn.Read(data)
	if err != nil {
		logger.Warning("login E:%s len:%d", err, lenn)
		conn.Close()
		return
	}
}

func (v verifyer) verifyLogin(data []byte) (userType, mountPointName, clientName string) {
	dataStr := string(data[:])
	fmt.Println(time.Second)
	if len(dataStr) < 15 {
		goto exit //  ==return NoType, "",""
	}
	if dataStr[0:6] == "SOURCE" { //ntripServer
		ds := strings.Split(dataStr, "\r\n")
		if len(ds) < 4 {
			goto exit
		}
		ds0s := strings.Split(ds[0], " ")
		if len(ds0s) != 3 {
			goto exit
		}
		mountPoint := ds0s[1]
		passwd := ds0s[2]
		if len(ds[1]) < 14 || len(ds[2]) < 5 || ds[1][:14] != "Source-Agent: " || ds[2][:5] != "STR: " {
			goto exit
		}
		sourceAgent := ds[1][14:]
		strs := ds[2][5:]

		if OK, err := v.loginType.mountpointLogin(mountPoint, passwd); OK {
			logger.Trace("verify pass tpye:%s , mountPoint:%s , Source-Agent:%s , STR:%s", MountpointType, mountPoint, sourceAgent, strs)
			return MountpointType, mountPoint, ""
		} else {
			logger.Trace("%s", err)
		}

	} else if dataStr[0:3] == "GET" { // ntripClient
		ds := strings.Split(dataStr, "\r\n")
		if len(ds) < 4 {
			goto exit //  ==return NoType, ""
		}
		ds0s := strings.Split(ds[0], " ")
		if len(ds0s) != 3 {
			goto exit
		}
		if ds0s[0] != "GET" || ds0s[2] != "HTTP/1.0" {
			goto exit
		}
		if len(ds0s[1]) < 2 || ds0s[1][0] != '/' {
			goto exit
		}
		mountPoint := ds0s[1][1:]
		if len(ds[1]) < 12 || ds[1][:12] != "User-Agent: " {
			goto exit
		}
		ds2s := strings.Split(ds[2], " ")
		if ds2s[0] != "Authorization:" || ds2s[1] != "Basic" {
			goto exit
		}
		res, err := base64.StdEncoding.DecodeString(ds2s[2])
		if err != nil {
			goto exit
		}
		userAndPasswdSp := strings.Split(string(res), ":")
		if len(userAndPasswdSp) != 2 {
			goto exit
		}
		clientName := userAndPasswdSp[0]
		passwd := userAndPasswdSp[1]
		if OK, err := v.loginType.clientLogin(mountPoint, clientName, passwd); OK {
			return ClientType, mountPoint, clientName
		} else {
			logger.Trace("%s", err)
		}
	}
exit:
	return NoType, "", ""

}
