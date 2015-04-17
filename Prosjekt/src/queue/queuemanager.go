package queue

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"
	"time"
)

type Elevator struct {
	IPADDR   string
	POSITION int
	/*
	   Etg.			Pos. nr.	ElevInfo.POSITION
	    1 ......... 0.........0
	  	  ......... 1.........0/1
		2 ......... 2.........1
		  ......... 3.........1/2
		3 ......... 4.........2
		  ......... 5.........2/3
		4 ......... 6.........3
	*/

	DIRECTION int
	/*
		opp = 1, ned = -1, stillestående = 0
	*/

	DESTINATION int
	/*
		1. etg = 0
		2. etg = 1
		3. etg = 2
		4. etg = 3
	*/

	ORDER_MATRIX [][]int
	/* 			   opp    	 ned    inne i heis			Settes til 1 ved en ordre
	1.etg	[[  0         0         0]
	2.etg 	 [  0         0         0]
	3.etg 	 [  0         0         0]
	4.etg	 [  0         0         0]]
	osv.
	*/
}

type ElevInfo struct {
	IPADDR     string
	F_NEW_INFO bool

	F_DEAD_ELEV   bool
	F_BUTTONPRESS bool

	POSITION    int
	DIRECTION   int
	DESTINATION int

	BUTTON_TYPE int
	BUTTONFLOOR int
}

const (
	N_FLOORS    = 4
	N_POSITIONS = N_FLOORS + (N_FLOORS - 1)
)

// Indexen i map'en er ip-adressen til den aktuelle heisen
var Active_elevators = make(map[string]Elevator)

// IP-adressen til "denne" heisen
var my_ipaddr string

func InitQueuemanager(ipaddr string, c_from_elevManager chan []byte, c_to_statemachine chan int, c_peerListUpdate chan string) {
	my_ipaddr = ipaddr
	my_ordermatrix := make([][]int, N_FLOORS)
	for i := 0; i < N_FLOORS; i++ {
		my_ordermatrix[i] = []int{0, 0, 0}
	}
	new_elevator := Elevator{my_ipaddr, 0, 0, 0, my_ordermatrix}
	Active_elevators[my_ipaddr] = new_elevator
	fmt.Println("Elevator", Active_elevators[my_ipaddr].IPADDR, "online\n")

	go processNewInfo(c_from_elevManager, c_peerListUpdate)
	go checkQueue(c_to_statemachine)
	fmt.Printf("Queuemanager operational\n")
}

// Denne funkjsonen brukes kun ifm debugging
func SetElevator(ipaddr string, position int, direction int, destinasjon_pos int) {
	temp := Active_elevators[ipaddr]
	temp.POSITION = position
	temp.DIRECTION = direction
	temp.DESTINATION = destinasjon_pos
	Active_elevators[ipaddr] = temp

}

func AppendElevator(ipaddr string) {
	new_ordermatrix := make([][]int, N_FLOORS)
	for i := 0; i < N_FLOORS; i++ {
		new_ordermatrix[i] = []int{0, 0, 0}
	}
	new_elevator := Elevator{ipaddr, 0, 0, 0, new_ordermatrix}
	Active_elevators[ipaddr] = new_elevator
	fmt.Println("Elevator", Active_elevators[ipaddr].IPADDR, "online\n")
}

func PrintActiveElevators() {
	fmt.Printf("************************************************************\n")
	for i := range Active_elevators {
		fmt.Println("Elevator:", Active_elevators[i].IPADDR)
		fmt.Println("Position:", Active_elevators[i].POSITION, "Direction:", Active_elevators[i].DIRECTION, "Destination floor:", Active_elevators[i].DESTINATION)
		// fmt.Println("Direction:", Active_elevators[i].DIRECTION)
		// fmt.Println("Destination:", Active_elevators[i].DESTINATION)
		fmt.Printf("Orders:\n")
		for floor := 0; floor < N_FLOORS; floor++ {
			fmt.Println("Floor", floor+1, ":", Active_elevators[i].ORDER_MATRIX[floor])
		}
		fmt.Printf("\n")
	}
	fmt.Printf("************************************************************\n")
	// fmt.Println("\n")
}

func PrintActiveElevators2() {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 0, '\t', 0)
	ipstr, infostr, orderstr, orderstr1, orderstr2, orderstr3, orderstr4 := "", "", "", "", "", "", ""
	// infostr := ""
	// orderstr := ""
	for key, elev := range Active_elevators {
		ipstr += "Elevator: " + key + "\t"
		infostr += "Position: " + strconv.Itoa(elev.POSITION) + "   Direction: " + strconv.Itoa(elev.DIRECTION) + "   Destination: " + strconv.Itoa(elev.DESTINATION) + "\t"

		tempstr := "     "
		for i := 0; i < 3; i++ {
			tempstr += strconv.Itoa(elev.ORDER_MATRIX[3][i]) + "     "
		}
		orderstr1 += "Floor: 3" + tempstr + "\t"

		tempstr = "     "
		for i := 0; i < 3; i++ {
			tempstr += strconv.Itoa(elev.ORDER_MATRIX[2][i]) + "     "
		}
		orderstr2 += "Floor: 2" + tempstr + "\t"

		tempstr = "     "
		for i := 0; i < 3; i++ {
			tempstr += strconv.Itoa(elev.ORDER_MATRIX[1][i]) + "     "
		}
		orderstr3 += "Floor: 1" + tempstr + "\t"

		tempstr = "     "
		for i := 0; i < 3; i++ {
			tempstr += strconv.Itoa(elev.ORDER_MATRIX[0][i]) + "     "
		}
		orderstr4 += "Floor: 0" + tempstr + "\t"

		orderstr += "Orders:     OPP   NED   INNE" + "\t"
	}
	fmt.Fprintln(w, ipstr)
	fmt.Fprintln(w, infostr)
	fmt.Fprintln(w, orderstr)
	fmt.Fprintln(w, orderstr1)
	fmt.Fprintln(w, orderstr2)
	fmt.Fprintln(w, orderstr3)
	fmt.Fprintln(w, orderstr4)
	fmt.Printf("*********************************************************************************************\n")
	w.Flush()
	fmt.Printf("*********************************************************************************************\n")
}

func RemoveElevator(ipaddr string) {
	orders_to_dist := Active_elevators[ipaddr].ORDER_MATRIX
	delete(Active_elevators, ipaddr)
	for floor := 0; floor < N_FLOORS; floor++ {
		for button_type := 0; button_type < 2; button_type++ {
			if orders_to_dist[floor][button_type] == 1 {
				AppendOrder(button_type, floor)
			}
		}
	}
	fmt.Println("Deleted", ipaddr, "\n")
}

// Bruker kostfunksjonen for å legge til ny ordre
func AppendOrder(button_type int, button_floor int) {
	fmt.Printf("Appending order\n")
	var button_dir string
	var optimal_elevatorIP string
	// Setter først kost urimelig høyt
	cost := 100

	if button_type == 0 {
		button_dir = "up"
	} else if button_type == 1 {
		button_dir = "down"
	} else if button_type == 2 {
		temp_elev := Active_elevators[my_ipaddr]
		for i := 0; i < 3; i++ {
			temp_elev.ORDER_MATRIX[button_floor][i] = 1
		}
		// temp_elev.ORDER_MATRIX[button_floor][button_type] = 1

		Active_elevators[my_ipaddr] = temp_elev
		return
	}

	for ipaddr := range Active_elevators {
		// fmt.Println("Cost:", CostFunction(ipaddr, button_floor, button_dir))
		new_cost := CostFunction(ipaddr, button_floor, button_dir)
		if new_cost < cost {
			cost = new_cost
			optimal_elevatorIP = ipaddr
		} else if new_cost == cost {
			old_ip_num, _ := strconv.Atoi(optimal_elevatorIP[12:len(optimal_elevatorIP)])
			new_ip_num, _ := strconv.Atoi(ipaddr[12:len(ipaddr)])

			if new_ip_num < old_ip_num {
				optimal_elevatorIP = ipaddr
			}
		}
	}
	// fmt.Println("Cost:", cost)
	// legger inn ordre i køen til den optimale heisen
	temp_elev := Active_elevators[optimal_elevatorIP]
	temp_elev.ORDER_MATRIX[button_floor][button_type] = 1
	Active_elevators[optimal_elevatorIP] = temp_elev
}

func deleteOrder(ipaddr string, floor int) {
	temp_elev := Active_elevators[ipaddr]
	for i := 0; i < 3; i++ {
		temp_elev.ORDER_MATRIX[floor][i] = 0
	}
	Active_elevators[ipaddr] = temp_elev

}

func CostFunction(elevator_ip string, order_floor int, button_dir string) int {
	cost := 0
	current_elevator := Active_elevators[elevator_ip]

	//Omregner etg. nr. til posisjonsnr. (Ihht. structen Elevator)
	order_floor_pos := order_floor * 2
	dest_pos := current_elevator.DESTINATION * 2

	switch {
	case current_elevator.DIRECTION == 0:
		if current_elevator.POSITION >= order_floor_pos {
			cost = current_elevator.POSITION - order_floor_pos
		} else {
			cost = order_floor_pos - current_elevator.POSITION
		}

	case button_dir == "up" && current_elevator.DIRECTION == 1:
		if current_elevator.POSITION <= order_floor_pos {

			if dest_pos >= order_floor_pos {
				cost = order_floor_pos - current_elevator.POSITION
			} else {
				// + 3 sek for dør-åpen-ventetid før man kjører videre mot bestilling
				cost = order_floor_pos - current_elevator.POSITION + 3
			}
		} else {
			cost = dest_pos - current_elevator.POSITION + 3 + dest_pos - order_floor_pos
		}

	case button_dir == "up" && current_elevator.DIRECTION == -1:
		cost = current_elevator.POSITION - dest_pos + 3 + order_floor_pos - dest_pos

	case button_dir == "down" && current_elevator.DIRECTION == -1:
		if current_elevator.POSITION >= order_floor_pos {
			if dest_pos <= order_floor_pos {
				cost = current_elevator.POSITION - order_floor_pos
			} else {
				cost = current_elevator.POSITION - order_floor_pos + 3
			}
		} else {
			cost = current_elevator.POSITION - dest_pos + 3 + order_floor_pos - dest_pos
		}

	case button_dir == "down" && current_elevator.DIRECTION == 1:
		cost = dest_pos - current_elevator.POSITION + 3 + dest_pos - order_floor_pos

	}

	return cost
}

// Får inn ny info fra heisManager (evt. timeout). Mottar pos og dir fra tilstandsmaskin.
// Må teste om deleteOrder() funker
func processNewInfo(c_from_elevManager chan []byte, c_peerListUpdate chan string) { //, c_pos_from_statemachine chan int, c_dir_from_statemachine chan int){
	var elev_info ElevInfo
	var last_info ElevInfo
	for {
		select {
		case encoded_elev_info := <-c_from_elevManager:			
			err := json.Unmarshal(encoded_elev_info, &elev_info)
			if err != nil {
				fmt.Println("error: ", err)
			}
			if _, ok := Active_elevators[elev_info.IPADDR]; !ok {
				AppendElevator(elev_info.IPADDR)
			}
			if elev_info.F_NEW_INFO && (elev_info != last_info) {

				temp_elev := Active_elevators[elev_info.IPADDR]
				temp_elev.POSITION = elev_info.POSITION*2
				temp_elev.DIRECTION = elev_info.DIRECTION

				// Sørger for at egen heis bare oppdaterer dest gjennom checkQueue()
				if elev_info.IPADDR != my_ipaddr {					
					temp_elev.DIRECTION = elev_info.DIRECTION
					temp_elev.DESTINATION = elev_info.DESTINATION
				}

				Active_elevators[elev_info.IPADDR] = temp_elev

				//elev_info.POSITION har samme sytax som Destination
				// Trenger kanskje ikke slette for hver gang vi får ny info?
				
				
				if elev_info.F_BUTTONPRESS == true {
					AppendOrder(elev_info.BUTTON_TYPE, elev_info.BUTTONFLOOR)
				}else if elev_info.POSITION == elev_info.DESTINATION {
					deleteOrder(elev_info.IPADDR, elev_info.POSITION)
					fmt.Printf("queue: Order completed, deleting\n")
				}

				
				last_info = elev_info

				PrintActiveElevators2()
			}else{
				last_info = elev_info
				fmt.Printf("queue: No new info\n")
			}

			if elev_info.F_DEAD_ELEV == true {
				RemoveElevator(elev_info.IPADDR)
			}

		case peerUpdate := <-c_peerListUpdate:
			RemoveElevator(peerUpdate)

			// All info kommer fra Router:
			/*
				case pos := <- c_pos_from_statemachine:
					temp_elev := Active_elevators[my_ipaddr]
					temp_elev.POSITION = pos
					Active_elevators[my_ipaddr] = temp_elev

					if (pos+2)/2 -1 == Active_elevators[my_ipaddr].DESTINATION {
						deleteOrder(my_ipaddr, (pos+2)/2 -1)
					}

					PrintActiveElevators2()

				case dir := <- c_dir_from_statemachine:
					temp_elev := Active_elevators[my_ipaddr]
					temp_elev.DIRECTION = dir
					Active_elevators[my_ipaddr] = temp_elev

					PrintActiveElevators2()
			*/
		}

	}
}

// Sjekker (ikke lenger hele tiden, men hvert 10 ms) hele tiden køen, oppdaterer next destination og sender denne til tilstandsmaskin.
func checkQueue(c_to_statemachine chan int) {
	var dest int
	var pos_floor int
	for {
		// elev := Active_elevators[my_ipaddr]
		time.Sleep(10 * time.Millisecond)
		switch {
		case Active_elevators[my_ipaddr].DIRECTION == 1:
			pos := Active_elevators[my_ipaddr].POSITION
			if Active_elevators[my_ipaddr].POSITION%2 == 0 {
				pos_floor = (pos+2)/2 - 1
			} else {
				pos_floor = ((pos+1)+2)/2 - 1
			}
			for i := pos_floor; i < (N_FLOORS - 1); i++ {
				if Active_elevators[my_ipaddr].ORDER_MATRIX[i][0] == 1 && i != Active_elevators[my_ipaddr].DESTINATION {
					dest = i
					fmt.Println("queue: New destination floor: ", dest)
					temp_elev := Active_elevators[my_ipaddr]
					temp_elev.DESTINATION = dest
					Active_elevators[my_ipaddr] = temp_elev
					c_to_statemachine <- dest
					time.Sleep(10 * time.Millisecond)

					PrintActiveElevators2()
				}
			}

		case Active_elevators[my_ipaddr].DIRECTION == -1:
			pos := Active_elevators[my_ipaddr].POSITION
			if Active_elevators[my_ipaddr].POSITION%2 == 0 {
				pos_floor = (pos+2)/2 - 1
			} else {
				pos_floor = ((pos-1)+2)/2 - 1
			}
			for i := pos_floor; i >= 0; i-- {
				if Active_elevators[my_ipaddr].ORDER_MATRIX[i][1] == 1 && i != Active_elevators[my_ipaddr].DESTINATION {
					dest = i
					fmt.Println("queue: New destination floor: ", dest)
					temp_elev := Active_elevators[my_ipaddr]
					temp_elev.DESTINATION = dest
					Active_elevators[my_ipaddr] = temp_elev
					c_to_statemachine <- dest
					time.Sleep(10 * time.Millisecond)
					PrintActiveElevators2()

				}
			}

		case Active_elevators[my_ipaddr].DIRECTION == 0:
			pos_floor = (Active_elevators[my_ipaddr].POSITION+2)/2 - 1
			for i := 0; i < (N_FLOORS); i++ {
				if Active_elevators[my_ipaddr].ORDER_MATRIX[i][0] == 1 {//&& i != Active_elevators[my_ipaddr].DESTINATION {
					dest = i
					// fmt.Println("New destination floor: ", dest)
					temp_elev := Active_elevators[my_ipaddr]
					temp_elev.DESTINATION = dest
					Active_elevators[my_ipaddr] = temp_elev
					c_to_statemachine <- dest
					time.Sleep(10 * time.Millisecond)
					PrintActiveElevators2()

				} else if Active_elevators[my_ipaddr].ORDER_MATRIX[i][1] == 1 {//&& i != Active_elevators[my_ipaddr].DESTINATION {
					dest = i
					// fmt.Println("New destination floor: ", dest)
					temp_elev := Active_elevators[my_ipaddr]
					temp_elev.DESTINATION = dest
					Active_elevators[my_ipaddr] = temp_elev
					c_to_statemachine <- dest
					time.Sleep(10 * time.Millisecond)
					PrintActiveElevators2()
				}
			}
		}
	}
}

func sendButtonLamp() {
	
}