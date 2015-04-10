package main

import (
	"./elevManager"
	"encoding/json"
	"fmt"
	"runtime"
	"strconv"
	"time"
)

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	c_mainchannel := make(chan []byte)
	c_peerUpdate := make(chan string)
	c_to_queuemanager := make(chan []byte)

	testInfo := elevManager.ElevInfo{}
	testInfo.POSITION = 1

	go elevManager.InitBank(c_mainchannel, c_peerUpdate, c_to_queuemanager)

	go SendShit(c_mainchannel, c_peerUpdate)

	for {
		select {
		case queue_info := <-c_to_queuemanager:
			json_err := json.Unmarshal(queue_info, &testInfo)
			if json_err != nil {
				fmt.Println("error: ", json_err)
			}
			fmt.Printf("Fikk info på queueChan om IP %s \n", testInfo.IPADDR)
			if testInfo.F_DEAD_ELEV {
				fmt.Printf("IP %s is dead! \n", testInfo.IPADDR)
			}
		case <-time.After(200 * time.Millisecond):
			fmt.Printf("Ingen ny queueinfo\n")
			if testInfo.POSITION == 0 {
				fmt.Printf("IPADDR: %s, Position: %d \n", testInfo.IPADDR, testInfo.POSITION)
				return
			}

		}
	}

}

func SendShit(c_mainchannel chan []byte, c_peerUpdate chan string) {
	time.Sleep(500 * time.Millisecond)
	testMelding := elevManager.ElevInfo{"0", false, false, false, 3, -1, 1, 1, 1}
	i := 0

	for {
		time.Sleep(500 * time.Millisecond)
		encoded_message, err2 := json.Marshal(testMelding)
		if err2 != nil {
			fmt.Println("error: ", err2)
		}

		c_mainchannel <- encoded_message
		fmt.Printf("Sendte melding til ip %s \n", testMelding.IPADDR)

		if i == 5 {
			time.Sleep(5 * time.Second)
			return
		}

		i++
		fmt.Println(i)
		if i < 2 {
			testMelding.IPADDR = strconv.Itoa(i)
		} else if i <= 3 {
			testMelding.F_NEW_INFO = true
			testMelding.POSITION = i
		} else if i == 4 {
			testMelding.IPADDR = "0"
			testMelding.F_NEW_INFO = false
			fmt.Printf("Kill elev 1 \n")
			c_peerUpdate <- "1"
		} else {
			testMelding.F_NEW_INFO = true
			testMelding.POSITION = 0
		}
	}

}
