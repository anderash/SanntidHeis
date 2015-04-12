package main

import (
	"./elevManager"
	"./queue"
	"./stateMachine"
	"./network"
	"./driver"
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
	c_io_input := make (chan []byte)


	go elevManager.InitBank(c_fromNetwork, c_peerUpdate, c_to_queuemanager)

	queue.InitQueuemanager(my_ipaddr, c_to_queuemanager, c_dest_to_statemachine, c_pos_from_statemachine, c_dir_from_statemachine)
	
	
	driver.InitDriver(c_io_input, c_SM_output, c_io_floor)


	stateMachine.InitStateMachine(c_dest_to_statemachine, c_io_floor, c_SM_output)

	go network.UDPNetwork(c_toNetwork, c_fromNetwork, c_peerUpdate)
	go AliveRoutine(my_ipaddr, c_toNetwork)
	
/*
	time.Sleep(5 * time.Second)




	
*/
	var decoded_input driver.Input
	for{
		select{
		case ioInput := <- c_io_input:
			err := json.Unmarshal(ioInput, &decoded_input)
			if err != nil{
				fmt.Println("error: ", err)
			}
			fmt.Println("Input:", decoded_input.INPUT_TYPE, "Button type:", decoded_input.BUTTON_TYPE, "Floor:", decoded_input.FLOOR)
			if decoded_input.INPUT_TYPE == driver.FLOOR_SENSOR {

			} else {
				testMelding := elevManager.ElevInfo{my_ipaddr, true, false, true, 0, 0, 0, decoded_input.BUTTON_TYPE, decoded_input.FLOOR}
				encoded_message, err := json.Marshal(testMelding)
				if err != nil {
					fmt.Println("error: ", err)
				}
				c_fromNetwork <- encoded_message
				fmt.Printf("ORDRE SENDT\n")
			}


		case <-time.After(500 * time.Millisecond):
				
		}
	}
}


func AliveRoutine(ip string, c_toNetwork chan []byte) {
	fmt.Printf("Started aliveroutine")
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