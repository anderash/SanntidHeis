package main

import (
	"fmt"
	"net"
	"os"
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

	host, _ := os.Hostname()
	addrs, err := net.LookupIP(host)
	for _, addr := range addrs {
   	 	if ipv4 := addr.To4(); ipv4 != nil {
    	    fmt.Println("IPv4: ", ipv4)
    	}   
	}

	ip, err := externalIP()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("1. test of IP is: %s", ip)


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
			fmt.Printf("Skriver toNetwork")
			c_toNetwork <- []byte(encoded_melding)
			time.Sleep(1000 * time.Millisecond)
		}
		select{
		case listenMessage := <-c_fromNetwork:
			err := json.Unmarshal(listenMessage, &recievedMessage)
			if err != nil {
				fmt.Println("error: ", err)
			}
			fmt.Printf("IP: %s, floor: %d, dead: %t", recievedMessage.IPADDR, recievedMessage.POSITION, recievedMessage.F_DEAD_ELEV)
			message.POSITION = message.POSITION - 1
			peerlist:=  <- c_peerList
			for i := range peerlist{
				fmt.Printf("IP is: %s \n", peerlist[i])
			}
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