package main

import ()

type statePutMysql struct {
}

func (state statePutMysql) setRoverSt(roverName, st string) {
	// fmt.Println("rover", roverName, st)
	setStatus("rover", roverName, st)

}

func (state statePutMysql) setMounpointSt(mountpointName, st string) {
	// fmt.Println("mountpoint", mountpointName, st)
	setStatus("mountpoint", mountpointName, st)

}

func (state statePutMysql) setAllDown() {
	setStatus("all", "", "")

}

func (state statePutMysql) initState(conf map[string]string) {
	host := conf["host"]
	port := conf["port"]
	user := conf["username"]
	passwd := conf["password"]
	mysqlInit(host, port, user, passwd)

}
