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


	var recieved Melding

	c_listen := make(chan []byte)
	c_broadcast := make(chan []byte)
	c_NrBytes := make(chan int)

	runtime.GOMAXPROCS(runtime.NumCPU())

	go network.UDPBroadcast(c_broadcast)
	go network.UDPListen(c_listen, c_NrBytes)

	for {
		if (melding.Floor > 0){
			encoded_melding, err2 := json.Marshal(melding)
			if err2 != nil {
				fmt.Println("error: ", err2)
			}
			c_broadcast <- []byte(encoded_melding)
			time.Sleep(1000 * time.Millisecond)
		} 
		select{
			case listen_message := <-c_lis ,ten:
				length := <- c_NrBytes
				stripped := listen_message[:length]
				err := json.Unmarshal(stripped, &recieved)
				if err != nil {
					fmt.Println("error: ", err)
				}

				fmt.Println("Alive = ", recieved.Alive, "Melding = ", recieved.Message, "Floor = ", recieved.Floor)
				melding.Floor = melding.Floor - 1
			case <-time.After(3000 * time.Millisecond):
				fmt.Printf("Timeout! Did not get a new message")
		}	
		
		
	}

	//fmt.Printf("Antall bytes sendt: %i", nrBsendt)

}
