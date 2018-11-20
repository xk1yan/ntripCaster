package main

import (
	"xk/ntripCaster/testDataMaker"
)

type usersIn struct {
	mountpointPasswdMap map[string]string
	clientPasswdMap     map[string]string
}

func (users usersIn) mountpointLogin(mountpointName, passwd string) (res bool, err error) {
	if users.mountpointPasswdMap[mountpointName] == passwd {
		return true, nil
	}
	return false, nil
}

func (users usersIn) clientLogin(mountpointName, clientName, passwd string) (res bool, err error) {
	if _, ok := users.mountpointPasswdMap[mountpointName]; ok {
		if users.clientPasswdMap[clientName] == passwd {
			return true, nil
		}
	}
	return false, nil
}

func (users *usersIn) usersInit(conf map[string]string) (err error) {

	if users.clientPasswdMap == nil {
		users.clientPasswdMap = make(map[string]string)
	}
	if users.mountpointPasswdMap == nil {
		users.mountpointPasswdMap = make(map[string]string)
	}

	mps := testDataMaker.CreateMountPointAndClient(100, 100)
	for k, v := range mps {
		// fmt.Println(k)
		users.mountpointPasswdMap[k] = "123456"
		for _, c := range v {
			users.clientPasswdMap[c] = "123456"
		}

	}
	users.mountpointPasswdMap["mount1"] = "passwd1"
	users.mountpointPasswdMap["yunzhihui"] = "yzhbase01"
	users.mountpointPasswdMap["jzsk2base2"] = "123456"

	users.clientPasswdMap["client1"] = "passwd2"
	users.clientPasswdMap["yzhrover03"] = "123456"
	return nil
}
