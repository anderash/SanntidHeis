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
	//infoPackage := make(chan []byte)
	testMelding := elevManager.ElevInfo{"0", false, false, 2, -1, 1, 1, 1}

	runtime.GOMAXPROCS(runtime.NumCPU())

	c_mainchannel := make(chan []byte)
	c_to_queuemanager := make(chan []byte)

	testInfo := elevManager.ElevInfo{}

	go elevManager.InitBank(c_mainchannel, c_to_queuemanager)

	i := 0

	for {
		select {
		case queue_info := <-c_to_queuemanager:
			json_err := json.Unmarshal(queue_info, &testInfo)
			if json_err != nil {
				fmt.Println("error: ", json_err)
			}
			fmt.Printf("Fikk info på queueChan om IP %s \n", testInfo.IPADDR)
		case <-time.After(3000 * time.Millisecond):
			fmt.Printf("Ingen ny queueinfo\n")

		}

		encoded_message, err2 := json.Marshal(testMelding)
		if err2 != nil {
			fmt.Println("error: ", err2)
		}

		c_mainchannel <- encoded_message
		fmt.Printf("Sendte melding til ip %s \n", testMelding.IPADDR)

		if i == 5 {
			fmt.Printf("IPADDR: %s, Position: %d \n", testInfo.IPADDR, testInfo.POSITION)
			time.Sleep(5 * time.Second)
			return
		}

		i++
		if i < 2 {
			testMelding.IPADDR = strconv.Itoa(i)
		} else if i == 3 {
			testMelding.F_NEW_INFO = true
			testMelding.POSITION = 2
		} else {
			testMelding.F_NEW_INFO = false
		}

		time.Sleep(1000 * time.Millisecond)

	}

}
