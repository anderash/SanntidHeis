package elevManager

import (
	"encoding/json"
	"fmt"
	//"time"
)

/* Denne skal opprette ny goroutine for ny heis.
Hver routine har en for select for enten oppdatert info, eller timeout.
Ved oppdatert info, send struct til kømodul.
Ved timeout, send KILL-heis til kømodul.
Husk å også drepe goroutine

*/

type ElevInfo struct {
	IPADDR     string
	F_NEW_INFO bool

	F_DEAD_ELEV   bool
	F_BUTTONPRESS bool

	POSITION    int
	DIRECTION   int
	DESTINATION int

	ButtonType  int
	ButtonFloor int
}

func InitBank(c_from_main <-chan []byte, c_peerListUpdate chan string, c_to_queuemanager chan<- []byte) {

	var info_package ElevInfo

	bank := make(map[string]ElevInfo)

	for {
		select {
		case from_main := <-c_from_main:

			json_err := json.Unmarshal(from_main, &info_package)
			if json_err != nil {
				fmt.Println("error: ", json_err)
			}

			fmt.Printf("Info om IP %s \n", info_package.IPADDR)
			_, in_bank := bank[info_package.IPADDR]
			fmt.Printf("IP %s er %t i bank\n", info_package.IPADDR, in_bank)

			bank[info_package.IPADDR] = info_package

			for key, value := range bank {
				fmt.Printf("Key: %s Position: %i Dead: %t \n", key, value.POSITION, value.F_DEAD_ELEV)
			}

			if !in_bank {
				c_to_queuemanager <- from_main

			} else if info_package.F_NEW_INFO {
				fmt.Printf("New info for IP %s, its now on floor %d \n", info_package.IPADDR, info_package.POSITION)
				c_to_queuemanager <- from_main

			} else {
				fmt.Printf("No new info, just a alive-ping from IP %s \n", info_package.IPADDR)
			}

		case peerUpdate := <-c_peerListUpdate:
			fmt.Printf("Recieved a dead elevator call IP: %s \n", peerUpdate)
			tmp := bank[peerUpdate]
			tmp.F_DEAD_ELEV = true
			bank[peerUpdate] = tmp

			encoded_message, err2 := json.Marshal(bank[peerUpdate])
			if err2 != nil {
				fmt.Println("JSON error: ", err2)
			}
			c_to_queuemanager <- encoded_message
		}
	}
}

/*
func spawnElevcheck(c_mychan chan bool, my_IP string) {
	for {
		select {
		case <-c_mychan:
			fmt.Printf("I: %s am alive\n", my_IP)
		case <-time.After(3000 * time.Millisecond):
			fmt.Printf("I: %s died\n", my_IP)
			//c_kømodul <- JegDøde
			c_mychan <- false
			return
		}
	}
}
*/
