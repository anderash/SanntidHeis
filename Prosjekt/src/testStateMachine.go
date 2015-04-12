package main

import (
	"./elevManager"
	"./queue"
	"./stateMachine"
	"./network"
	"fmt"
	"runtime"
	"time"
	"encoding/json"
)

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())
	my_ipaddr := network.GetOwnIP()

	c_peerUpdate := make(chan string)
	c_to_queuemanager := make(chan []byte)
	c_dest_to_statemachine := make(chan int)
	c_pos_from_statemachine := make(chan int)
	c_dir_from_statemachine := make(chan int)
	c_io_floor := make(chan int)
	c_toNetwork := make(chan []byte)
	c_fromNetwork := make(chan []byte)
	c_SM_output := make(chan []byte)

	go network.UDPNetwork(c_toNetwork, c_fromNetwork, c_peerUpdate)
	go elevManager.InitBank(c_fromNetwork, c_peerUpdate, c_to_queuemanager)

	queue.InitQueuemanager(my_ipaddr, c_to_queuemanager, c_dest_to_statemachine, c_pos_from_statemachine, c_dir_from_statemachine)
	stateMachine.InitStateMachine(c_dest_to_statemachine, c_io_floor, c_SM_output)

	go AliveRoutine(my_ipaddr, c_toNetwork)
/*
	time.Sleep(5 * time.Second)



	testMelding := elevManager.ElevInfo{my_ipaddr, true, false, true, 0, 0, 0, 1, 3}
	encoded_message, err := json.Marshal(testMelding)
	if err != nil {
		fmt.Println("error: ", err)
	}
	c_fromNetwork <- encoded_message
	fmt.Printf("ORDRE SENDT\n")

*/
	for{
		select{
			case <-time.After(500 * time.Millisecond):
				
		}
	}
}


func AliveRoutine(ip string, c_toNetwork chan []byte) {

	message := elevManager.ElevInfo{ip, false, false, false, 0,0,0,1,3} 
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