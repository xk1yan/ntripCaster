package main

import (
	"fmt"
)

type statePutNone struct {
}

func (state statePutNone) setRoverSt(roverName, st string) {
	fmt.Println("rover", roverName, st)

}

func (state statePutNone) setMounpointSt(mountpointName, st string) {
	fmt.Println("mountpoint", mountpointName, st)

}

func (state statePutNone) setAllDown() {

}

func (state statePutNone) initState(conf map[string]string) {

}
