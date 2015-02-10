// main
package main

import (
	//"fmt"
	"./network"
	"runtime"
	)

func main() {
	melding := "I am alive"

	c_listen := make(chan []byte)
	c_broadcast := make(chan []byte)


	runtime.GOMAXPROCS(runtime.NumCPU())

	go network.UDPBroadcast(c_broadcast)
	go network.UDPListen(c_listen)

	
	c_broadcast <- []byte(melding)

	c := make(chan int)
	<- c

	//fmt.Printf("Antall bytes sendt: %i", nrBsendt)


}