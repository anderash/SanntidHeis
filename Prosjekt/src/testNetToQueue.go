package main

import (
	"./elevManager"
	"./network"
	"./queue"
	"encoding/json"
	"fmt"
	"runtime"
	"time"
)

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

//	c_mainchannel := make(chan []byte)
	c_peerUpdate := make(chan string)
	c_queMan_info := make(chan []byte)

	ipaddr := network.GetOwnIP()
	c_queMan_destination := make(chan int)
	c_SM_position := make(chan int)
	c_SM_direction := make(chan int)

	c_toNetwork := make(chan []byte)
	c_fromNetwork := make(chan []byte)

	go network.UDPNetwork(c_toNetwork, c_fromNetwork, c_peerUpdate)
	go elevManager.InitBank(c_fromNetwork, c_peerUpdate, c_queMan_info)
	queue.InitQueuemanager(ipaddr, c_queMan_info, c_queMan_destination, c_SM_position, c_SM_direction)

	go AliveRoutine(ipaddr, c_toNetwork)

	for{
		select{
			case <-time.After(500 * time.Millisecond):
				
		}
	}

}

func AliveRoutine(ip string, c_toNetwork chan []byte) {

	message := elevManager.ElevInfo{ip, false, false, false, 0,0,0,0,0} 
	time.Sleep(500 * time.Millisecond)

	for{	
		encoded_melding, err2 := json.Marshal(message)
			if err2 != nil {
				fmt.Println("AliveRoutin JSON error: ", err2)
			}
		c_toNetwork <- []byte(encoded_melding)
		time.Sleep(500 * time.Millisecond)		
	}
}