package main

import (
	"fmt"
	"net"

	"./network"
	"encoding/json"
	"runtime"
	"time"
	"errors"
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

	ip, err := externalIP()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("1. test of IP is: %s \n", ip)


	c_toNetwork := make(chan []byte)
	c_fromNetwork := make(chan []byte)
	c_peerList := make(chan []string)

	runtime.GOMAXPROCS(runtime.NumCPU())

	go network.UDPNetwork(c_toNetwork, c_fromNetwork, c_peerList)

	message := ElevInfo{ip, true, false, false, 3,-1,1,0,0} 
	var recievedMessage ElevInfo
	//var peerlist []string
	
	for{
		if (message.POSITION > 0){
			encoded_melding, err2 := json.Marshal(message)
			if err2 != nil {
				fmt.Println("error: ", err2)
			}
			fmt.Printf("Skriver toNetwork\n")
			c_toNetwork <- []byte(encoded_melding)
			fmt.Printf("Position is now: %d\n", message.POSITION)
			message.POSITION = message.POSITION - 1
			//time.Sleep(1000 * time.Millisecond)
		} else {
			return
		}
		fmt.Printf("entering select statement\n")
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
		case <-time.After(900 * time.Millisecond):
			fmt.Printf("Timeout! Did not get a new message\n")
		} 


	}
}

func externalIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String(), nil
		}
	}
	return "", errors.New("are you connected to the network?")
}