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

	c_io_button := make(chan []byte)
	c_io_floor := make(chan int)

	c_peerUpdate := make(chan string)
	c_toNetwork := make(chan []byte)
	c_fromNetwork := make(chan []byte)

	c_router_info := make(chan []byte)
	c_elevMan_info := make(chan []byte)

	c_queMan_dest := make(chan int)
	c_SM_output := make(chan []byte)
	c_SM_floor := make(chan int)





}

func router(my_ip string, c_io_floor <-chan []byte, c_io_button <-chan []byte, c_SM_output <-chan []byte, c_toNetwork chan<- []byte, c_router_info chan<- []byte, c_SM_floor chan<- int) {

var elevator elevManager.ElevInfo
var myElevator elevManager.ElevInfo
var floorInput int	
var SM_output stateMachine.Output

for{
	select{
		case floor_input <- c_io_floor:



	}
}
}

func getInput() {
	
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

