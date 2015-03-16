package elevManager

import (
	"fmt"
	"time"
	"encoding/json"
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

	F_DEAD_ELEV bool
	F_BUTTONPRESS bool	

	POSITION    int
	DIRECTION   int
	DESTINATION int

	ButtonType  int
	ButtonFloor int
}

func InitBank(c_from_main chan []byte, c_to_queuemanager chan []byte) {

	var info_package ElevInfo

	bank := make(map[string]ElevInfo)
	elev_channels := make(map[string]chan bool, 1)

	for {
		select {
		case from_main := <- c_from_main:

			json_err := json.Unmarshal(from_main, &info_package)
			if json_err != nil {
				fmt.Println("error: ", json_err)
			}

			fmt.Printf("Info om IP %s \n", info_package.IPADDR)
			_, in_bank := bank[info_package.IPADDR]
			fmt.Printf("IP %s er %t i bank\n", info_package.IPADDR, in_bank)

			bank[info_package.IPADDR] = info_package

			for key, value := range(bank) {
			    fmt.Println("Key:", key, "Value:", value.IPADDR, "\n")
			}

			if !in_bank {
				c_newchannel := make(chan bool)
				go spawnElevcheck(c_newchannel, info_package.IPADDR)

				elev_channels[info_package.IPADDR] = c_newchannel
				fmt.Printf("Made new goroutine for elevator" + info_package.IPADDR + "\n")

				encoded_message, err2 := json.Marshal(info_package)
				if err2 != nil {
					fmt.Println("error: ", err2)
				}
				c_to_queuemanager <- encoded_message

			} else if info_package.F_NEW_INFO {
				fmt.Printf("New info for IP %s, its now on %d floor", info_package.IPADDR, info_package.POSITION)

				elev_channels[info_package.IPADDR] <- true

				encoded_message, err2 := json.Marshal(info_package)
				if err2 != nil {
					fmt.Println("error: ", err2)
				}
				c_to_queuemanager <- encoded_message

			} else {
				// Kun fått alive-melding, send på den det gjelder sin kanal at du lever
				elev_channels[info_package.IPADDR] <- true
			}
/*
		for IP, channel := range elev_channels{
			case check <- channel:
				fmt.Printf("Checked channel %s \n", IP)
				if check == false{
					delete(bank, "IP")
				}

		}
*/
		}
	}
}

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
