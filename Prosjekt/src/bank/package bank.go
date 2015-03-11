package bank

import {
	"fmt"
}

/* Denne skal opprette ny goroutine for ny heis.
Hver routine har en for select for enten oppdatert info, eller timeout.
Ved oppdatert info, send struct til kømodul.
Ved timeout, send KILL-heis til kømodul.
Husk å også drepe goroutine

*/

func InitBank(Kanal fra main) {
	
	make nyInfoTilKømodulKanal

	for{
		select{
			case: fraMian <- mainChan{
				dekode info fraMain
				if (nyElev) {
					make nyElevKanal
					spawnElevcheck(nyElecKanal)
				} else if (oppdatertinfopåeksisterendeHeis){
					denDetGjelderSinKanal <- DuLever
					nyInfoTilKømodulKanal <- NyInfo
				} else{
					// Kun fått alive-melding, send på den det gjelder sin kanal at du lever
					denDetGjelderSinKanal <- DuLever
				}

			}
		}
	}
}


func spawnElevcheck(ny kanal for ny heis) {
	for{
		select{
			
		}
	} 
}