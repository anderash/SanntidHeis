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

func UDPBroadcast(data str) {
	buffer := make([]byte(data))

	raddr, err1 := net.resolveUDPaddr("udp", Baddr+":"+OwnPort)

		if err != nil {
		fmt.Printf("Addresse dritt")
		os.Exit(1)
		}

	socket, err2 := net.DialUDP ("udp", nil, raddr)

	n, _ = socket.Write(buffer)

	

}

func UDPListen(){
	buffer := make([]byte)


}