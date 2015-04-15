package main

import (
	"./driver"
	"./elevManager"
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
	c_toNetwork := make(chan []byte)	//elevMan.ElevInfo
	c_fromNetwork := make(chan []byte)	//elevMan.ElevInfo

	c_router_info := make(chan []byte)	//elevMan.ElevInfo
	c_elevMan_info := make(chan []byte)	//elevMan.ElevInfo

//	c_queMan_button := make(chan []byte) // This channel sets button lights in IO from queueManager
	c_queMan_dest := make(chan int)	// int dest
	c_SM_output := make(chan []byte)	//stateMachine.Output
	c_SM_floor := make(chan int)	//int floor






}

/*
The router takes in info fromchannels,
and send it to those modules that need the update: 
IO button channel: toNet(myIP, info), elevMan(myIP, info)
IO floor channel

*/
func router(my_ip string, c_io_floor <-chan []byte, c_io_button <-chan []byte, c_SM_output <-chan []byte, c_toNetwork chan<- []byte, c_router_info chan<- []byte, c_SM_floor chan<- int) {

var elevator elevManager.ElevInfo
var floor_input int	
var SM_output stateMachine.Output
var buttonpress driver.Input

myElevator := elevManager.ElevInfo{my_ip, true, false, false,0,0,0,0,0};

for{
	select{
		case floor_input <- c_io_floor:
			myElevator.POSITION = floor_input
			sendElev(myElevator,c_router_info)
			sendElev(myElevator,c_toNetwork)

		case button_input := <- c_io_button:
			json_err := json.Unmarshal(button_input, &buttonpress)
			if json_err != nil {
				fmt.Println("router unMarshal JSON error: ", json_err)
			}
			myElevator.F_BUTTONPRESS = true
			myElevator.ButtonType = buttonpress.BUTTON_TYPE
			myElevator.ButtonFloor = buttonpress.FLOOR
			sendButtonpress(myElevator, c_router_info)
			sendButtonpress(myElevator, c_toNetwork)

			case 



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


