package main

import (
	"encoding/base64"
	"strings"
)

const (
	mountpointType = "mountpoint"
	clientType     = "client"
	noType         = "no"
)

var backStr = [...]string{
	"HTTP/1.0 400 Bad Request\r\nServer: NTRIP NtripCaster\r\nContent-Type: text/html\r\nConnection: close\r\n",
	"ICY 200 OK\r\n",
	"ERROR - Bad Password\r\n",
	"HTTP/1.0 401 Unauthorized\r\n",
	"OK\r\n",
}

type verifyBack struct {
	userType       string
	mountPointName string
	clientName     string
	backStrIndex   int8
}

func verifyLogin(loginType userIfa, data []byte) *verifyBack {
	back := new(verifyBack)
	dataStr := string(data[:])
	if len(dataStr) < 15 {
		back.backStrIndex = 0
		goto exit
	}
	if dataStr[0:6] == "SOURCE" { //ntripServer
		ds := strings.Split(dataStr, "\r\n")
		if len(ds) < 4 {
			back.backStrIndex = 0
			goto exit
		}
		ds0s := strings.Split(ds[0], " ")
		if len(ds0s) != 3 {
			back.backStrIndex = 0
			goto exit
		}
		mountPoint := ds0s[2]
		passwd := ds0s[1]
		if len(ds[1]) < 14 || len(ds[2]) < 5 || ds[1][:14] != "Source-Agent: " || ds[2][:5] != "STR: " {
			back.backStrIndex = 0
			goto exit
		}
		sourceAgent := ds[1][14:]
		strs := ds[2][5:]

		if OK, err := loginType.mountpointLogin(mountPoint, passwd); OK {
			logger.Trace("verify pass tpye:%s , mountPoint:%s , Source-Agent:%s , STR:%s", mountpointType, mountPoint, sourceAgent, strs)
			back.userType = mountpointType
			back.mountPointName = mountPoint
			back.backStrIndex = 4
			return back
		} else {
			back.backStrIndex = 2
			logger.Trace("%s", err)
		}

	} else if dataStr[0:3] == "GET" { // ntripClient
		ds := strings.Split(dataStr, "\r\n")
		if len(ds) < 4 {
			back.backStrIndex = 0
			goto exit
		}
		ds0s := strings.Split(ds[0], " ")
		if len(ds0s) != 3 {
			goto exit
			back.backStrIndex = 0
		}
		if ds0s[0] != "GET" || ds0s[2] != "HTTP/1.0" {
			back.backStrIndex = 0
			goto exit
		}
		if len(ds0s[1]) < 2 || ds0s[1][0] != '/' {
			goto exit
			back.backStrIndex = 0
		}
		mountPoint := ds0s[1][1:]
		if len(ds[1]) < 12 || ds[1][:12] != "User-Agent: " {
			back.backStrIndex = 0
			goto exit
		}
		ds2s := strings.Split(ds[2], " ")
		if ds2s[0] != "Authorization:" || ds2s[1] != "Basic" {
			back.backStrIndex = 0
			goto exit
		}
		res, err := base64.StdEncoding.DecodeString(ds2s[2])
		if err != nil {
			back.backStrIndex = 0
			goto exit
		}
		userAndPasswdSp := strings.Split(string(res), ":")
		if len(userAndPasswdSp) != 2 {
			goto exit
		}
		clientName := userAndPasswdSp[0]
		passwd := userAndPasswdSp[1]
		if OK, err := loginType.clientLogin(mountPoint, clientName, passwd); OK {
			logger.Trace("verify pass tpye:%s , mountPoint:%s , clientName:%s", clientType, mountPoint, clientName)
			back.userType = clientType
			back.mountPointName = mountPoint
			back.clientName = clientName
			back.backStrIndex = 1
			return back
		} else {
			back.backStrIndex = 3
			logger.Trace("%s", err)
		}
	}
exit:
	return back

}
