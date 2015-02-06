// main
package main

import (
	//"fmt"
	"network"
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

	c := make(chan int)
	<- c

	//fmt.Printf("Antall bytes sendt: %i", nrBsendt)


}