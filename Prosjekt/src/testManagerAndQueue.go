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

	c_mainchannel := make(chan []byte)
	c_peerUpdate := make(chan string)
	c_to_queuemanager := make(chan []byte)
	ipaddr := network.GetOwnIP()
	c_to_statemachine := make(chan int)
	c_pos_from_statemachine := make(chan int)
	c_dir_from_statemachine := make(chan int)

	testInfo := elevManager.ElevInfo{}
	testInfo.POSITION = 1

	go elevManager.InitBank(c_mainchannel, c_peerUpdate, c_to_queuemanager)
	queue.InitQueuemanager(ipaddr, c_to_queuemanager, c_to_statemachine, c_pos_from_statemachine, c_dir_from_statemachine)

	me := elevManager.ElevInfo{ipaddr, false, false, false, 0, 0, 0, 0, 0}
	time.Sleep(50 * time.Millisecond)
	encoded_message, err := json.Marshal(me)
	if err != nil {
		fmt.Println("error: ", err)
	}

	c_mainchannel <- encoded_message
	fmt.Printf("Sendt myself to elevBank \n")

	go SendShit2(c_mainchannel, c_peerUpdate)

	for {
		select {
		case queue_info := <-c_to_statemachine:
			fmt.Printf("Info to statemachine: ")
			fmt.Println(queue_info)

		case <-time.After(10000 * time.Millisecond):
			return

		}
	}

}

func SendShit(c_mainchannel chan []byte, c_peerUpdate chan string) {

	time.Sleep(2 * time.Second)
	testMelding := elevManager.ElevInfo{"0", false, false, false, 0, 0, 0, 1, 3}
	i := 0

	for {
		time.Sleep(500 * time.Millisecond)
		encoded_message, err2 := json.Marshal(testMelding)
		if err2 != nil {
			fmt.Println("error: ", err2)
		}

		c_mainchannel <- encoded_message
		fmt.Printf("Sendte %d melding til ip %s \n", i, testMelding.IPADDR)

		if i == 5 {
			time.Sleep(5 * time.Second)
			return
		}

		i++
		fmt.Println(i)
		if i == 2 {
			//testMelding.IPADDR = strconv.Itoa(i)
			testMelding.F_NEW_INFO = true
			testMelding.F_BUTTONPRESS = true
		} else if i == 3 {
			testMelding.F_NEW_INFO = false
		} else if i == 4 {
			testMelding.F_NEW_INFO = false
			fmt.Printf("Kill elev 0 \n")
			c_peerUpdate <- "0"
		} else {
			testMelding.F_NEW_INFO = false			
		}
	}
}

func SendShit2(c_mainchannel chan []byte, c_peerUpdate chan string) {
	time.Sleep(2 * time.Second)
	testMelding := elevManager.ElevInfo{"129.241.187.09", true, false, false, 0, 0, 0, 0, 0}
	encoded_message, err3 := json.Marshal(testMelding)
	if err3 != nil {
			fmt.Println("error: ", err3)
		}
	c_mainchannel <- encoded_message

	time.Sleep(2 * time.Second)
	testMelding2 := elevManager.ElevInfo{"129.241.187.124", true, false, true, 0, 0, 0, 1, 3}
	encoded_message, err4 := json.Marshal(testMelding2)
	if err4 != nil {
			fmt.Println("error: ", err4)
		}
	c_mainchannel <- encoded_message

	time.Sleep(2 * time.Second)
	testMelding3 := elevManager.ElevInfo{"129.241.187.09", true, false, false, 3, 0, 3, 0, 0}
	encoded_message, err5 := json.Marshal(testMelding3)
	if err5 != nil {
			fmt.Println("error: ", err5)
		}
	c_mainchannel <- encoded_message
}
