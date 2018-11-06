package main

import ()

type testlogin struct {
	mountpointPasswdMap map[string]string
	clientPasswdMap     map[string]string
}

func (t testlogin) mountpointLogin(mountpointName, passwd string) (res bool, err error) {
	if t.mountpointPasswdMap[mountpointName] == passwd {
		return true, nil
	}
	return false, nil
}

func (t testlogin) clientLogin(mountpointName, clientName, passwd string) (res bool, err error) {
	if _, ok := t.mountpointPasswdMap[mountpointName]; ok {
		if t.clientPasswdMap[clientName] == passwd {
			return true, nil
		}
	}
	return false, nil
}

func (t *testlogin) updateUserMap() (err error) {

	if t.clientPasswdMap == nil {
		t.clientPasswdMap = make(map[string]string)
	}
	if t.mountpointPasswdMap == nil {
		t.mountpointPasswdMap = make(map[string]string)
	}
	t.mountpointPasswdMap["mount1"] = "passwd1"
	t.mountpointPasswdMap["yunzhihui"] = "yzhbase01"
	t.mountpointPasswdMap["jzsk2base2"] = "123456"

	t.clientPasswdMap["client1"] = "passwd2"
	t.clientPasswdMap["yzhrover03"] = "123456"

	return nil
}
