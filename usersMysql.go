package main

import (
	"fmt"
)

type usersMysql struct {
}

func (users usersMysql) mountpointLogin(mountpointName, passwd string) (res bool, err error) {
	res = mountpointVer(mountpointName, passwd)
	return res, nil
}

func (users usersMysql) clientLogin(mountpointName, clientName, passwd string) (res bool, err error) {
	res = clientVer(clientName, passwd)
	return res, nil
}

/*
conf file

[mysql]
username=root
password=123456
host=127.0.0.1
port=3306

*/
func (users *usersMysql) usersInit(conf map[string]string) (err error) {
	// var err error = nil
	host := conf["host"]
	port := conf["port"]
	user := conf["username"]
	passwd := conf["password"]
	fmt.Println(host, port, user, passwd)
	err = mysqlInit(host, port, user, passwd)

	return err
}
