package elevManager

import (
	"encoding/json"
	"fmt"
	//"time"
)

/* 
Tar ny info fra main og sjekker om det er noe ny info.
Sender isåfall videre til Queue Manager for oppdatering
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

func InitBank(c_from_main <-chan []byte, c_peerListUpdate chan string, c_elevMan_info chan<- []byte) {

	var info_package ElevInfo

	bank := make(map[string]ElevInfo)

	for {
		select {
		case from_main := <-c_from_main:

			json_err := json.Unmarshal(from_main, &info_package)
			if json_err != nil {
				fmt.Println("elevMan unMarshal JSON error: ", json_err)
			}

			_, in_bank := bank[info_package.IPADDR]

			bank[info_package.IPADDR] = info_package

			if !in_bank {
				c_elevMan_info <- from_main

			} else if info_package.F_NEW_INFO {
				fmt.Printf("New info for IP %s \n", info_package.IPADDR)
				c_elevMan_info <- from_main

			} else {
				fmt.Printf("Alive-ping from IP %s \n", info_package.IPADDR)
			}

		case peerUpdate := <-c_peerListUpdate:
			fmt.Printf("Recieved a dead elevator call IP: %s \n", peerUpdate)
			tmp := bank[peerUpdate]
			tmp.F_DEAD_ELEV = true
			bank[peerUpdate] = tmp

			encoded_message, err2 := json.Marshal(bank[peerUpdate])
			if err2 != nil {
				fmt.Println("elevMan Marshal JSON error: ", err2)
			}
			c_elevMan_info <- encoded_message
		}
	}
}
