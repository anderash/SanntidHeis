package main

import (
	"./elevManager"
	"./network"
	"./queue"
//	"encoding/json"
//	"fmt"
	"runtime"
	"time"
)

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

//	c_mainchannel := make(chan []byte)
	c_peerUpdate := make(chan string)
	c_to_queuemanager := make(chan []byte)

	ipaddr := network.GetOwnIP()
	c_to_statemachine := make(chan int)
	c_pos_from_statemachine := make(chan int)
	c_dir_from_statemachine := make(chan int)

	c_toNetwork := make(chan []byte)
	c_fromNetwork := make(chan []byte)

	go network.UDPNetwork(c_toNetwork, c_fromNetwork, c_peerUpdate)
	go elevManager.InitBank(c_fromNetwork, c_peerUpdate, c_to_queuemanager)
	queue.InitQueuemanager(ipaddr, c_to_queuemanager, c_to_statemachine, c_pos_from_statemachine, c_dir_from_statemachine)

	for{
		select{
			case <-time.After(5000 * time.Millisecond):
				return
		}
	}

}