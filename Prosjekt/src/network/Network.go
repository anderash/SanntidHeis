// NETWORK MODULE //
package network

import (
  "fmt"
  "net"
  "time"
  "os"
  "errors"
)

const (
	OwnIP  = "129.241.187.140"
	OwnPort  = "20001"
	Baddr  = "129.241.187.255"
	aliveInterval = 200 * time.Millisecond
)

var LocalIP string

func UDPNetwork(c_toNetwork <-chan []byte, c_fromNetwork chan<- []byte, c_peerListUpdate chan<- []string) 
{
	

	go udpBroadcast(c_toNetwork)
	go udpListen(c_fromNetwork, c_peerListUpdate)

	for{
		select{
		case toNetwork <- c_toNetwork:

		}
	}


}

func getOwnIP() string {
	if LocalIP == "" {

	}
}


func udpBroadcast(c_broadcast chan []byte) {
	
	raddr, err1 := net.ResolveUDPAddr("udp", Baddr+":"+OwnPort)

		if err1 != nil {
			fmt.Printf("Problemer med resolveUDPaddr")
			os.Exit(1)
		}

	socket, err2 := net.DialUDP ("udp", nil, raddr)

		if err2 != nil {
			fmt.Printf("Problemer med Dial")
			os.Exit(2)
		}	

	for {
		buffer := <- c_broadcast
		_ , err3 := socket.Write(buffer)
		//fmt.Printf("skrev %i bytes", n)

		if err3 != nil {
			fmt.Printf("Problemer med Write")
			os.Exit(3)
		}		
	}

}

func udpListen(c_listen chan []byte, c_NrBytes chan int, c_peerListUpdate chan <- []string){
	buffer := make([]byte, 1024)

	raddr, err1 := net.ResolveUDPAddr("udp4", Baddr+":"+OwnPort)
		
		if err1 != nil {
			fmt.Printf("Problemer med resolveUDPaddr")
			os.Exit(4)
		}

	socket, _ := net.ListenUDP("udp4", raddr)


	var listHasChanges bool
	var peerList []string


	for {
		
		nrBytes, remoteADDR, err4 := socket.ReadFromUDP(buffer)

		socket.SetReadDeadline(time.Now().Add(2*aliveInterval))
		listHasChanges = false


		if err4 != nil {
			fmt.Printf("Problemer med ReadFromUDP")
			os.Exit(5)
		}

		//fmt.Printf(string(buffer))
		c_listen <- buffer
		c_NrBytes <- nrBytes
		//time.Sleep(100*time.Millisecond)

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
