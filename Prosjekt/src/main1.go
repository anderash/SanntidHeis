// main
package main

import (
	//"fmt"
	"network"
	"time"
	"runtime"
	)

func main() {
	melding := "I am alive"

	c_listen := make(chan []byte)
	c_broadcast := make(chan []byte)

	c_broadcast <- []byte(melding)

	runtime.GOMAXPROCS(runtime.NumCPU())

	go network.UDPBroadcast(c_broadcast)
	go network.UDPListen(c_listen)

	time.Sleep(100*time.Millisecond)

	//fmt.Printf("Antall bytes sendt: %i", nrBsendt)


}