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
	c_io_floor := make(chan int)	// int floor

	c_peerUpdate := make(chan string) //IP-adress
	c_toNetwork := make(chan []byte)	//queuemanager.ElevInfo
	c_fromNetwork := make(chan []byte)	//queuemanager.ElevInfo

	c_router_info := make(chan []byte)	//queuemanager.ElevInfo

//	c_queMan_button := make(chan []byte) // This channel sets button lights in IO from queueManager
	c_queMan_dest := make(chan int)	// int dest

	c_SM_state := make(chan []byte)		//stateMachine.Output
	c_SM_output := make(chan []byte)	//stateMachine.Output
	c_SM_floor := make(chan int)	//int floor






}

/*
The router takes in info fromchannels,
and send it to those modules that need the update: 
IO button channel: toNet(myIP, info), queuemanager(myIP, info)
IO floor channel

*/
func router(my_ip string,  c_io_button <-chan []byte, c_SM_state <-chan []byte, c_toNetwork chan<- []byte, c_router_info chan<- []byte, c_SM_floor chan<- int) {

	var elevator queuemanager.ElevInfo
	var floor_input int	
	var state stateMachine.ElevState
	var buttonpress driver.Input

	myElevator := queuemanager.ElevInfo{my_ip, true, false, false,0,0,0,0,0};

	for{
		select{

			case e_button_input := <- c_io_button:
				json_err := json.Unmarshal(e_button_input, &buttonpress)
				if json_err != nil {
					fmt.Println("router unMarshal JSON error: ", json_err)
				}
				myElevator.F_BUTTONPRESS = true
				myElevator.ButtonType = buttonpress.BUTTON_TYPE
				myElevator.ButtonFloor = buttonpress.FLOOR
				sendButtonpress(myElevator, c_router_info)
				sendButtonpress(myElevator, c_toNetwork)

			case e_state := <- c_SM_State:
				json_err := json.Unmarshal(e_state, &state)
				if json_err != nil {
					fmt.Println("router unMarshal JSON error: ", json_err)
				}
				myElevator.POSITION = state.POSITION
				myElevator.DIRECTION = state.DIRECTION
				myElevator.DESTINATION = state.DESTINATION			


		}
	}
}



func sendElev(info ElevInfo, channel chan []byte) {
	encoded_output, err := json.Marshal(info)
	if err != nil {
		fmt.Println("SM JSON error: ", err)
	}
	channel <- encoded_output		
}

func sendSMOutput(output Output, channel chan []byte) {
	encoded_output, err := json.Marshal(output)
	if err != nil {
		fmt.Println("SM JSON error: ", err)
	}
	channel <- encoded_output	
}

func sendButtonpress(info Input, channel chan []byte) {
	encoded_output, err := json.Marshal(info)
	if err != nil {
		fmt.Println("SM JSON error: ", err)
	}
	channel <- encoded_output	
}


