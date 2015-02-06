// main
package main

import (
	"fmt"
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

	go UDPBroadcast(c_broadcast)
	go UDPLIsten(c_listen)

	time.Sleep(100*time.Millisecond)

	fmt.Printf("Antall bytes sendt: %i", nrBsendt)


}