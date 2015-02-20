// main
package main

import (
	"fmt"
	"./network"
	"runtime"
	"time"

	)

func main() {
	melding := "I am alive"

	c_listen := make(chan []byte)
	c_broadcast := make(chan []byte)


	runtime.GOMAXPROCS(runtime.NumCPU())

	go network.UDPBroadcast(c_broadcast)
	go network.UDPListen(c_listen)

	for{
		c_broadcast <- []byte(melding)
		time.Sleep(1000*time.Millisecond)
		listen_message := <- c_listen
		fmt.Printf("%s", string(listen_message)+"\n")
	}
	



	//fmt.Printf("Antall bytes sendt: %i", nrBsendt)


}