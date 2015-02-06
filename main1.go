// main
package main

import (
	"fmt"
	"network"
	)

func main() {
	melding := "I am alive"

	nrBsendt := UDPBroadcast(melding)
	fmt.Printf("Antall bytes sendt: %i", nrBsendt)
	

}