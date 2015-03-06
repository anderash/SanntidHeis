// Main for testing network module
package main

import (
	"./network"
	"encoding/json"
	"fmt"
	"runtime"
	"time"
)

type Melding struct {
	Alive bool

	Message string

	Floor int
}

func main() {
	melding := Melding{true, "Hei", 4}
	encoded_melding, err2 := json.Marshal(melding)
	if err2 != nil {
		fmt.Println("error: ", err2)
	}

	var recieved Melding

	c_listen := make(chan []byte)
	c_broadcast := make(chan []byte)
	c_NrBytes := make(chan int)

	runtime.GOMAXPROCS(runtime.NumCPU())

	go network.UDPBroadcast(c_broadcast)
	go network.UDPListen(c_listen, c_NrBytes)

	for {
		c_broadcast <- []byte(encoded_melding)
		time.Sleep(1000 * time.Millisecond)
		listen_message := <-c_listen
		length := <- c_NrBytes
		stripped := listen_message[:length]
		err := json.Unmarshal(stripped, &recieved)
		if err != nil {
			fmt.Println("error: ", err)
		}

		fmt.Println("Alive = ", recieved.Alive, "Melding = ", recieved.Message, "Floor = ", recieved.Floor)
	}

	//fmt.Printf("Antall bytes sendt: %i", nrBsendt)

}
