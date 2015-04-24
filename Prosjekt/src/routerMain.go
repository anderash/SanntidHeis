package main

import (
	"./driver"
	//	"./elevManager"
	"./network"
	"./queue"
	"./stateMachine"
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"time"
)

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	my_ipaddr := network.GetOwnIP()
	fmt.Printf("Versjon 4, ip %s \n", my_ipaddr)

	c_io_button := make(chan []byte) //Type: driver.Input; io -> router
	c_io_floor := make(chan int)     //Type: int 		; io -> stMachine

	c_peerUpdate := make(chan string)  //Type: string		; network ->
	c_toNetwork := make(chan []byte)   //Type: queue.ElevInfo
	c_fromNetwork := make(chan []byte) //Type: queue.ElevInfo

	c_router_info := make(chan []byte) //Type: queue.ElevInfo

	c_queMan_dest := make(chan int)        //Type: int dest
	c_queMan_output := make(chan []byte)   //Type: queue.Output This channel sets button lights in IO from queueManager
	c_queMan_ackOrder := make(chan []byte) //Type: queue.Elevinfo

	c_stMachine_output := make(chan []byte) //Type: stateMachine.Output
	c_stMachine_state := make(chan []byte)  //Type: stateMachine.ElevState
	c_forloop := make(chan bool)

	go router(my_ipaddr, c_fromNetwork, c_io_button, c_stMachine_state, c_queMan_ackOrder, c_toNetwork, c_router_info)

	driver.InitDriver(c_io_button, c_io_floor, c_stMachine_output, c_queMan_output)

	queue.InitQueuemanager(my_ipaddr, c_router_info, c_queMan_dest, c_peerUpdate, c_queMan_output, c_queMan_ackOrder)

	go network.UDPNetwork(c_toNetwork, c_fromNetwork, c_peerUpdate)

	stateMachine.InitStatemachine(c_queMan_dest, c_io_floor, c_stMachine_output, c_stMachine_state)

	<-c_forloop

}

func router(my_ipaddr string, c_fromNetwork <-chan []byte, c_io_button <-chan []byte, c_stMachine_state <-chan []byte, c_queMan_ackOrder chan []byte, c_toNetwork chan<- []byte, c_router_info chan<- []byte) {

	var state stateMachine.ElevState
	var buttonpress driver.Input

	program_timer := time.NewTimer(10 * time.Second)

	myElevator := queue.ElevInfo{my_ipaddr, true, false, false, 0, 0, 0, false, 0, 0}

	for {
		select {

		case enc_button_input := <-c_io_button:
			json_err := json.Unmarshal(enc_button_input, &buttonpress)
			if json_err != nil {
				fmt.Println("router button unMarshal JSON error: ", json_err)
			}
			myElevator.F_NEW_INFO = true
			myElevator.F_BUTTONPRESS = true
			myElevator.BUTTON_TYPE = buttonpress.BUTTON_TYPE
			myElevator.BUTTONFLOOR = buttonpress.FLOOR

			sendElev(myElevator, c_router_info)
			if buttonpress.BUTTON_TYPE != 2 { // Does not broadcast if internal button
				sendElev(myElevator, c_toNetwork)

			}
			myElevator.F_BUTTONPRESS = false

		case enc_state := <-c_stMachine_state:
			json_err := json.Unmarshal(enc_state, &state)
			if json_err != nil {
				fmt.Println("router state unMarshal JSON error: ", json_err)
			}
			myElevator.F_NEW_INFO = true
			myElevator.POSITION = state.POSITION
			myElevator.DIRECTION = state.DIRECTION
			myElevator.DESTINATION = state.DESTINATION
			myElevator.MOVING = state.MOVING
			sendElev(myElevator, c_router_info)
			sendElev(myElevator, c_toNetwork)
			program_timer.Reset(10 * time.Second)

		case netInfo := <-c_fromNetwork:
			c_router_info <- netInfo

		case enc_order := <-c_queMan_ackOrder:
			c_toNetwork <- enc_order

		case <-time.After(250 * time.Millisecond):
			myElevator.F_NEW_INFO = false
			sendElev(myElevator, c_toNetwork)

		case <-program_timer.C:
			if myElevator.DIRECTION == 0 {
				fmt.Printf("10 sek since last state-update, in idle\n")
				program_timer.Reset(10 * time.Second)
			} else {
				fmt.Printf("Encountered an error. Elevator standing. Crashing program. Call maintnaince \n")
				os.Exit(1)
			}
		}
	}
}

func sendElev(info queue.ElevInfo, channel chan<- []byte) {
	encoded_output, err := json.Marshal(info)
	if err != nil {
		fmt.Println("sendElev JSON error: ", err)
	}
	channel <- encoded_output
}
