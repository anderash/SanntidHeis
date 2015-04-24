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
	"os"
)
const(
	N_FLOORS = 4
)

func main() {

	

	runtime.GOMAXPROCS(runtime.NumCPU())

	my_ipaddr := network.GetOwnIP()
	fmt.Printf("Versjon 3, ip %s \n", my_ipaddr)

	c_io_button := make(chan []byte) //driver.Input
	c_io_floor := make(chan int)     // int floor

	c_peerUpdate := make(chan string)  //IP-adress
	c_toNetwork := make(chan []byte)   //queue.ElevInfo
	c_fromNetwork := make(chan []byte) //queue.ElevInfo

	c_router_info := make(chan []byte) //queue.ElevInfo
	
	// c_queMan_button := make(chan []byte) // This channel sets button lights in IO from queueManager
	c_queMan_dest := make(chan int)      // int dest
	c_queMan_output := make(chan []byte) // This channel sets button lights in IO from queueManager
	c_queMan_ackOrder := make(chan []byte)   //queue.Elevinfo (Sends acknowledgment if order is handled by my IP for broadcasting)

	c_stMachine_output := make(chan []byte) //stateMachine.Output
	c_stMachine_state := make(chan []byte)  //stateMachine.Output
	c_forloop := make(chan bool)

	go router(my_ipaddr, c_fromNetwork, c_io_button, c_stMachine_state, c_queMan_ackOrder, c_toNetwork, c_router_info)

	driver.InitDriver(c_io_button, c_io_floor, c_stMachine_output, c_queMan_output)
	
	queue.InitQueuemanager(my_ipaddr, c_router_info, c_queMan_dest, c_peerUpdate, c_queMan_output, c_queMan_ackOrder)

	go network.UDPNetwork(c_toNetwork, c_fromNetwork, c_peerUpdate)

	stateMachine.InitStatemachine(c_queMan_dest, c_io_floor, c_stMachine_output, c_stMachine_state)

	<-c_forloop

}

/*
The router takes in info fromchannels,
and send it to those modules that need the update:
IO button channel: toNet(myIP, info), queue(myIP, info)
IO floor channel

*/
func router(my_ipaddr string, c_fromNetwork <-chan []byte, c_io_button <-chan []byte, c_stMachine_state <-chan []byte, c_queMan_ackOrder chan []byte, c_toNetwork chan<- []byte, c_router_info chan<- []byte) {

	var state stateMachine.ElevState
	var buttonpress driver.Input
	var order queue.ElevInfo

	program_timer := time.NewTimer(10 * time.Second)
	//doorTimer.Stop()
	my_ordermatrix := make([][]int, N_FLOORS)
	for i := 0; i < N_FLOORS; i++ {
		my_ordermatrix[i] = []int{0, 0, 0}
	}
	myElevator := queue.ElevInfo{my_ipaddr, true, false, false, false, 0, 0, 0, false, 0, 0, my_ordermatrix}

	for {
		select {

		case enc_button_input := <-c_io_button:
			fmt.Printf("Router: ButtonInput \n")
			json_err := json.Unmarshal(enc_button_input, &buttonpress)
			if json_err != nil {
				fmt.Println("router button unMarshal JSON error: ", json_err)
			}
			myElevator.F_NEW_INFO = true
			myElevator.F_BUTTONPRESS = true
			myElevator.BUTTON_TYPE = buttonpress.BUTTON_TYPE
			myElevator.BUTTONFLOOR = buttonpress.FLOOR
			fmt.Println(buttonpress)
			fmt.Printf("router input Trying to send\n")
			sendElev(myElevator, c_router_info)
			if buttonpress.BUTTON_TYPE != 2 { 	// Does not broadcast if internal button
				sendElev(myElevator, c_toNetwork)
				
			}
			myElevator.F_BUTTONPRESS = false

		case enc_state := <-c_stMachine_state:
			fmt.Printf("Router: StateInput \n")
			json_err := json.Unmarshal(enc_state, &state)
			if json_err != nil {
				fmt.Println("router floor unMarshal JSON error: ", json_err)
			}
			fmt.Println(state)
			myElevator.F_NEW_INFO = true
			myElevator.POSITION = state.POSITION
			myElevator.DIRECTION = state.DIRECTION
			myElevator.DESTINATION = state.DESTINATION
			myElevator.MOVING = state.MOVING
			fmt.Printf("router state Trying to send\n")
			sendElev(myElevator, c_router_info)
			sendElev(myElevator, c_toNetwork)
			fmt.Printf("router state sendt\n")
			program_timer.Reset(10*time.Second)

		case enc_netInfo := <-c_fromNetwork:
			fmt.Printf("router fromNet Trying to send\n")
			c_router_info <- enc_netInfo


		case enc_order := <- c_queMan_ackOrder:
			fmt.Printf("router ack recieved\n")
			json_err := json.Unmarshal(enc_order, &order)
			if json_err != nil {
				fmt.Println("router ack unMarshal JSON error: ", json_err)
			}
			myElevator.ORDER_MATRIX = order.ORDER_MATRIX
			fmt.Println(order.ORDER_MATRIX)
			fmt.Printf("router ack Trying to send\n")
			c_toNetwork <- enc_order
			

		case <-time.After(500 * time.Millisecond):
			myElevator.F_NEW_INFO = false
			sendElev(myElevator, c_toNetwork)

		case <- program_timer.C:
			if myElevator.DIRECTION == 0{
				fmt.Printf("10 sek since last state-update, in idle\n")
				program_timer.Reset(10*time.Second)
			}else{
				fmt.Printf("Encountered an error, crashing program. Call maintnaince \n")
				os.Exit(1)
			}
		}
	}
}

func sendElev(info queue.ElevInfo, channel chan<- []byte) {
	encoded_output, err := json.Marshal(info)
	if err != nil {
		fmt.Println("router sendElev JSON error: ", err)
	}
	channel <- encoded_output
}

