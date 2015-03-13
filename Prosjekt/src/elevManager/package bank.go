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

type main_Info {
	IPADDR string
	F_NEW_INFO bool

	POSITION int
	DIRECTION int
	DESTINATION int

	ButtonType int
	ButtonFloor int
}

type queueModuleInfo{
	IPADDR string
	F_DEAD_ELEV bool

	POSITION int
	DIRECTION int
	DESTINATION int

	ButtonType int
	ButtonFloor int
}


func InitBank(Kanal fra main) {
	
	make c_Kømodul

	for{
		select{
			case: fraMain <- mainChan{
				dekode info fraMain
				if (nyElev) {
					make nyElevKanal
					go spawnElevcheck(nyElecKanal)
					c_Kømodul <- NyHeis

				} else if (main_Info.FLAG_NEW_INFO == true){
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