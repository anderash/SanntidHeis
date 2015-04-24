// NETWORK MODULE //
package network

import (
	"encoding/json"
	"fmt"
	. "net"
	"os"
	"sort"
	. "strings"
	"time"
)

const (
	OwnPort       = "20004"
	MsgPort       = "20004"
	Baddr         = "129.241.187.255"
	aliveInterval = 250 * time.Millisecond
	deadTimeout   = 2 * time.Second
)

type ElevInfo struct {
	IPADDR     string
	F_NEW_INFO bool

	F_DEAD_ELEV   bool
	F_BUTTONPRESS bool
	F_ACK_ORDER   bool

	POSITION    int
	DIRECTION   int
	DESTINATION int
	MOVING      bool

	ButtonType  int
	ButtonFloor int
}

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
	var info_package ElevInfo

	for {
		socket.SetReadDeadline(time.Now().Add(8 * aliveInterval))
		nrBytes, remoteADDR, err := socket.ReadFromUDP(buffer)

		listHasChanges = false
		peerList = nil
		if err == nil {
			_, inList := lastSeen[IPString(remoteADDR)]
			if !inList {
				listHasChanges = true
			}
			lastSeen[IPString(remoteADDR)] = time.Now()
		} else {
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
		}
		if err == nil {
			stripped_info := buffer[:nrBytes]

			json_err := json.Unmarshal(stripped_info, &info_package)
			if json_err != nil {
				fmt.Println("network unMarshal JSON error: ", json_err)
			}
			if info_package.F_NEW_INFO || listHasChanges {
				fmt.Printf("New info arrived from IP %s \n", info_package.IPADDR)
				c_fromNetwork <- stripped_info
			}
		}

	}
}
