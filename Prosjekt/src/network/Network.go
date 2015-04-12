// NETWORK MODULE //
package network

import (
	"fmt"
	. "net"
	"os"
	"sort"
	. "strings"
	"time"
)

const (
	OwnPort       = "20005"
	MsgPort       = "20005"
	Baddr         = "129.241.187.255"
	aliveInterval = 500 * time.Millisecond
	deadTimeout   = 1 * time.Second
)

var localIP string

func UDPNetwork(c_toNetwork <-chan []byte, c_fromNetwork chan<- []byte, c_peerListUpdate chan<- string) {
	localIP = GetOwnIP()
	fmt.Printf("GetOwnIP returns: %s \n", localIP)

	addr, err := ResolveUDPAddr("udp4", Baddr+":"+MsgPort)
	if err != nil {
		fmt.Printf("Problemer med resolveUDPaddr\n")
		os.Exit(1)
	}
	msgConn, _ := DialUDP("udp4", nil, addr)

	go udpListen(c_fromNetwork, c_peerListUpdate)

	for {
		select {
		case msg := <-c_toNetwork:
			msgConn.Write(msg)
			fmt.Printf("Sendt message \n")
		}

	}

}

func GetOwnIP() string {
	if localIP == "" {
		addr, _ := ResolveTCPAddr("tcp4", "google.com:80")
		conn, _ := DialTCP("tcp4", nil, addr)
		localIP = IPString(conn.LocalAddr())
	}
	return localIP
}

func IPString(addr Addr) string {
	return Split(addr.String(), ":")[0]
}

/*
func udpBroadcast(c_toNetwork <-chan []byte) {

	raddr, err1 := ResolveUDPAddr("udp4", Baddr+":"+OwnPort)

	if err1 != nil {
		fmt.Printf("Problemer med resolveUDPaddr")
		os.Exit(1)
	}
	fmt.Printf("Trying to dialUDP\n")
	socket, err2 := DialUDP("udp4", nil, raddr)

	if err2 != nil {
		fmt.Printf("Problemer med Dial\n")
		os.Exit(2)
	}

	for {
		buffer := <-c_toNetwork
		fmt.Printf("Trying to Write\n")
		_, err3 := socket.Write(buffer)
		//fmt.Printf("skrev %i bytes", n)

		if err3 != nil {
			fmt.Printf("Problemer med Write")
			os.Exit(3)
		}
	}

}
*/

func udpListen(c_fromNetwork chan<- []byte, c_peerListUpdate chan<- string) {
	buffer := make([]byte, 1024)

	raddr, err1 := ResolveUDPAddr("udp4", ":"+OwnPort)

	if err1 != nil {
		fmt.Printf("Problemer med resolveUDPaddr")
		os.Exit(4)
	}

	socket, _ := ListenUDP("udp4", raddr)
	fmt.Printf("Created listenSocket\n")

	lastSeen := make(map[string]time.Time)
	var listHasChanges bool
	var peerList []string

	for {
		socket.SetReadDeadline(time.Now().Add(2 * aliveInterval))
		nrBytes, remoteADDR, err := socket.ReadFromUDP(buffer)		
		

		listHasChanges = false
		peerList = nil
		if err == nil {
			_, inList := lastSeen[IPString(remoteADDR)]
			if !inList {
				listHasChanges = true
			}
			lastSeen[IPString(remoteADDR)] = time.Now()
		} else{
			fmt.Println("Network error:", err)
		}


		for key, value := range lastSeen {
			if time.Now().Sub(value) > deadTimeout {
				delete(lastSeen, key)
				fmt.Printf("Timeout on elev %s \n", key)
				c_peerListUpdate <- key
				listHasChanges = true
			}

		}
		if listHasChanges {
			for key := range lastSeen {
				peerList = append(peerList, key)
			}
			sort.Strings(peerList)
			//c_peerListUpdate <- peerList
		}
		if err == nil{
			stripped := buffer[:nrBytes]
			c_fromNetwork <- stripped
			//c_NrBytes <- nrBytes
		}
	}

}

/*
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
*/
