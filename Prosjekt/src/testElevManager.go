package main

import (
	"./elevManager"
	"encoding/json"
	"fmt"
	"runtime"
)

type ElevInfo struct {
	IPADDR     string
	F_NEW_INFO bool

	F_DEAD_ELEV bool

	POSITION    int
	DIRECTION   int
	DESTINATION int

	ButtonType  int
	ButtonFloor int
}

func main() {
	infoPackage := make(chan []byte)
	testMelding := main_info{"123.456.789", true, false, 3, -1, 1, 1, 1}

	runtime.GOMAXPROCS(runtime.NumCPU())
	go InitBank(infoPackage)

}
