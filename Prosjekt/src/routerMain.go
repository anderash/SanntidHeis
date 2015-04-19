package main

import (
	"./driver"
	//	"./elevManager"
	"./network"
	"./queue"
	"./stateMachine"
	"encoding/json"
	"fmt"
	"runtime"
	"time"
)

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	my_ipaddr := network.GetOwnIP()

	c_io_button := make(chan []byte) //driver.Input
	c_io_floor := make(chan int)     // int floor

	c_peerUpdate := make(chan string)  //IP-adress
	c_toNetwork := make(chan []byte)   //queue.ElevInfo
	c_fromNetwork := make(chan []byte) //queue.ElevInfo

	c_router_info := make(chan []byte) //queue.ElevInfo

	// c_queMan_button := make(chan []byte) // This channel sets button lights in IO from queueManager
	c_queMan_dest := make(chan int)      // int dest
	c_queMan_output := make(chan []byte) // This channel sets button lights in IO from queueManager

	c_SM_output := make(chan []byte) //stateMachine.Output
	c_SM_state := make(chan []byte)  //stateMachine.Output
	c_forloop := make(chan bool)

	go router(my_ipaddr, c_fromNetwork, c_io_button, c_SM_state, c_toNetwork, c_router_info)

	queue.InitQueuemanager(my_ipaddr, c_router_info, c_queMan_dest, c_peerUpdate, c_queMan_output)

	driver.InitDriver(c_io_button, c_io_floor, c_SM_output, c_queMan_output)

	go network.UDPNetwork(c_toNetwork, c_fromNetwork, c_peerUpdate)

	stateMachine.InitStatemachine(c_queMan_dest, c_io_floor, c_SM_output, c_SM_state)

	<- c_forloop

}

/*
The router takes in info fromchannels,
and send it to those modules that need the update:
IO button channel: toNet(myIP, info), queue(myIP, info)
IO floor channel

*/
func router(my_ipaddr string, c_fromNetwork <- chan []byte, c_io_button <-chan []byte, c_SM_state <-chan []byte, c_toNetwork chan<- []byte, c_router_info chan<- []byte) {

	var state stateMachine.ElevState
	var buttonpress driver.Input

	myElevator := queue.ElevInfo{my_ipaddr, true, false, false, 0, 0, 0, 0, 0}

	for {
		select {

		case e_button_input := <-c_io_button:
			fmt.Printf("Router: ButtonInput \n")
			json_err := json.Unmarshal(e_button_input, &buttonpress)
			if json_err != nil {
				fmt.Println("router unMarshal JSON error: ", json_err)
			}
			myElevator.F_NEW_INFO = true
			myElevator.F_BUTTONPRESS = true
			myElevator.BUTTON_TYPE = buttonpress.BUTTON_TYPE
			myElevator.BUTTONFLOOR = buttonpress.FLOOR
			sendElev(myElevator, c_router_info)
			if (buttonpress.BUTTON_TYPE != 2){ 		//Sender ikke på nett om det er en intern knapp
				sendElev(myElevator, c_toNetwork)
			}			
			myElevator.F_BUTTONPRESS = false

		case e_state := <-c_SM_state:
			fmt.Printf("Router: StateInput \n")
			json_err := json.Unmarshal(e_state, &state)
			if json_err != nil {
				fmt.Println("router unMarshal JSON error: ", json_err)
			}
			fmt.Println(state)
			myElevator.F_NEW_INFO = true
			myElevator.POSITION = state.POSITION
			myElevator.DIRECTION = state.DIRECTION
			myElevator.DESTINATION = state.DESTINATION
			sendElev(myElevator, c_router_info)
			sendElev(myElevator, c_toNetwork)

		case netInfo := <- c_fromNetwork:
			c_router_info <- netInfo
		// Send Alive-Ping
		case <-time.After(500 * time.Millisecond):
//			fmt.Printf("Router: Ping \n")
			myElevator.F_NEW_INFO = false
			sendElev(myElevator, c_toNetwork)

		}
	}
}

func sendElev(info queue.ElevInfo, channel chan<- []byte) {
	encoded_output, err := json.Marshal(info)
	if err != nil {
		fmt.Println("SM JSON error: ", err)
	}
	channel <- encoded_output
}

func sendSMOutput(output stateMachine.Output, channel chan<- []byte) {
	encoded_output, err := json.Marshal(output)
	if err != nil {
		fmt.Println("SM JSON error: ", err)
	}
	channel <- encoded_output
}

func sendButtonpress(info driver.Input, channel chan<- []byte) {
	encoded_output, err := json.Marshal(info)
	if err != nil {
		fmt.Println("SM JSON error: ", err)
	}
	channel <- encoded_output
}
