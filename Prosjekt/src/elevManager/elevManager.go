package elevManager

import {
	"fmt"
}

/* Denne skal opprette ny goroutine for ny heis.
Hver routine har en for select for enten oppdatert info, eller timeout.
Ved oppdatert info, send struct til kømodul.
Ved timeout, send KILL-heis til kømodul.
Husk å også drepe goroutine

*/

type ElevInfo {
	IPADDR string
	F_NEW_INFO bool

	F_DEAD_ELEV bool

	POSITION int
	DIRECTION int
	DESTINATION int

	ButtonType int
	ButtonFloor int
}



func InitBank(c_from_main chan []byte) {
	
	c_to_queuemanager := make(chan []byte)
	var info_package ElevInfo
	var info_to_quemanager queueModuleInfo

	bank := make(map[string]ElevInfo)

	for{
		select{
			case: from_main <- c_from_main{
				json_err := json.Unmarshal(from_main, &info_package)
				if json_err != nil {
					fmt.Println("error: ", err)
				}

				elev, in_bank := bank[info_package.IPADDR]
				if (!in_bank) {
					c_newchannel := make(chan bool)
					go spawnElevcheck(c_newchannel)


					c_to_queuemanager <- NyHeis

				} else if (ElevInfo.FLAG_NEW_INFO == true){
					denDetGjelderSinKanal <- DuLever
					c_Kømodul <- NyInfo

				} else{
					// Kun fått alive-melding, send på den det gjelder sin kanal at du lever
					denDetGjelderSinKanal <- DuLever
				}

			}
		}
	}
}


func spawnElevcheck(minKanal) {
	for{
		select{
			case: <- minKanal
				//Jeg er i live
			case: <- timeout
				//Jeg døde
				 c_kømodul <- JegDøde
				 return
		}
	} 
}