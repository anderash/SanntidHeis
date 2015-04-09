package main

import (
	"fmt"

	"./network"
	"encoding/json"
	"runtime"
	"time"
)

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

func main() {

	ownIP := network.GetOwnIP()

	c_toNetwork := make(chan []byte)
	c_fromNetwork := make(chan []byte)
	c_peerList := make(chan []string)	

	runtime.GOMAXPROCS(runtime.NumCPU())

	go network.UDPNetwork(c_toNetwork, c_fromNetwork, c_peerList)

	go SendShit(ownIP, c_toNetwork)

	var recievedMessage ElevInfo
	//var peerlist []string
	time.Sleep(100 * time.Millisecond)

	fmt.Printf("entering select statement\n")

	for{
		select{
			case listenMessage := <-c_fromNetwork:
				err1 := json.Unmarshal(listenMessage, &recievedMessage)
				if err1 != nil {
					fmt.Println("error1: ", err1)
				}
				fmt.Printf("IP: %s, floor: %d, dead: %t \n", recievedMessage.IPADDR, recievedMessage.POSITION, recievedMessage.F_DEAD_ELEV)
		case peerlist :=  <- c_peerList:
			for i := range peerlist{
				fmt.Printf("IP is: %s \n", peerlist[i])
			}
		case <-time.After(500 * time.Millisecond):
			fmt.Printf("Timeout! Did not get a new message\n")
			if(recievedMessage.POSITION == 1){
				time.Sleep(1000 * time.Millisecond)
				return
			}
		} 	
	}
}

func SendShit(ip string, c_toNetwork chan []byte) {

	message := ElevInfo{ip, false, false, false, 2, 0,1,0,0} 
	time.Sleep(400 * time.Millisecond)

	for{
		encoded_melding, err2 := json.Marshal(message)
		if err2 != nil {
			fmt.Println("error: ", err2)
		}
		fmt.Printf("Skriver toNetwork\n")
		c_toNetwork <- []byte(encoded_melding)		
		time.Sleep(500 * time.Millisecond)
	}
}