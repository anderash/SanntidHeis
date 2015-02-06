// NETWORK MODULE //
package network

import (
  "fmt"
  "net"
  "time"
  "os"
)

constant (
	OwnIP := "129.241.187.140"
	OwnPort := "20001"
	Baddr := "129.241.187.255"
)

func UDPBroadcast(c_broadcast chan []byte) {
	buffer := make([]byte(<- c_broadcast))

	raddr, err1 := net.resolveUDPaddr("udp", Baddr+":"+OwnPort)

		if err != nil {
		fmt.Printf("Problemer med resolveUDPaddr")
		os.Exit(1)
		}

	socket, err2 := net.DialUDP ("udp", nil, raddr)

		if err2 != nil {
		fmt.Printf("Problemer med Dial")
		os.Exit(2)
		}	

	n, _ = socket.Write(buffer)


}

func UDPListen(c_listen chan []byte){
	buffer := make([]byte)

	raddr, err1 := net.resolveUDPaddr("udp", nil, Baddr+":"OwnPort)
		
		if err1 != nil {
		fmt.Printf("Problemer med resolveUDPaddr")
		os.Exit(3)
		}

	socket, _ := net.ListenUDP("udp4", raddr)

	for {
		n, _, _ = socket.ReadFromUDP(buffer)
		c_listen <- buffer

	}


}