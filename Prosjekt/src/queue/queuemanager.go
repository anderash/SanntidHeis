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
	   Etg.		Pos. nr.	ElevInfo.POSITION
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
	/* 			   opp    	 ned    inne i heis
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

	F_BUTTONPRESS bool
	F_ACK_ORDER   bool

	POSITION    int
	DIRECTION   int
	DESTINATION int
	MOVING      bool

	BUTTON_TYPE int
	BUTTONFLOOR int
}

type Output struct {
	OUTPUT_TYPE int
	/*
		LIGHT_OUTPUT = 0
		MOTOR_OUTPUT = 1
	*/

	LIGHT_TYPE int
	/*
		BUTTON_LAMP = 0
		FLOOR_INDICATOR = 1
		DOOR_LAMP = 2
	*/

	BUTTON_TYPE int
	/*
		BUTTON_CALL_UP = 0
	    BUTTON_CALL_DOWN = 1
	    BUTTON_COMMAND = 2
	    NO_BUTTON = -1
	*/

	FLOOR int

	VALUE int
	/*
		on = 1
		off = 0
	*/

	OUTPUT_DIRECTION int
	/*
		UP = 1
		STOP = 0
		DOWN = -1
	*/
}

type ButtonLight struct {
	BUTTON_TYPE int
	BUTTONFLOOR int
	VALUE       int
}

const (
	N_FLOORS    = 4
	N_POSITIONS = N_FLOORS + (N_FLOORS - 1)
)

var Active_elevators = make(map[string]Elevator)

var my_ipaddr string
var internal_orders []byte
var order_file *os.File
var acknowledgeTimer *time.Timer

func InitQueuemanager(ipaddr string, c_router_info chan []byte, c_to_statemachine chan int, c_peerListUpdate chan string, c_queMan_output chan []byte, c_queMan_ack_order chan []byte) {
	my_ipaddr = ipaddr
	my_ordermatrix := make([][]int, N_FLOORS)

	for i := 0; i < N_FLOORS; i++ {
		my_ordermatrix[i] = []int{0, 0, 0}
	}

	enc_file, f_err := os.OpenFile("internal_orders.dat", os.O_RDWR, 0777)

	if f_err != nil {
		fmt.Println("Created order_file\n")
		enc_file, _ = os.Create("internal_orders.dat")
	}

	order_file = enc_file
	internal_orders = make([]byte, N_FLOORS)
	_, r_err := order_file.ReadAt(internal_orders, 0)
	fmt.Println(internal_orders)
	if r_err != nil {
		fmt.Println("read order_file error: ", r_err)
	}

	for i := 0; i < N_FLOORS; i++ {
		my_ordermatrix[i][0] = int(internal_orders[i])
		my_ordermatrix[i][1] = int(internal_orders[i])
		my_ordermatrix[i][2] = int(internal_orders[i])
		if int(internal_orders[i]) == 1 {
			button_output := Output{0, 0, 2, i, 1, -1}
			sendButtonLamp(button_output, c_queMan_output)
		}
	}

	new_elevator := Elevator{my_ipaddr, 0, 0, 0, my_ordermatrix}
	Active_elevators[my_ipaddr] = new_elevator
	fmt.Println("Elevator", Active_elevators[my_ipaddr].IPADDR, "online\n")

	go processNewInfo(c_router_info, c_peerListUpdate, c_queMan_output, c_queMan_ack_order)
	go findNewDestination(c_to_statemachine)
	fmt.Printf("Queuemanager operational\n")
}

func appendElevator(elev_info ElevInfo) {
	new_ordermatrix := make([][]int, N_FLOORS)

	for i := 0; i < N_FLOORS; i++ {
		new_ordermatrix[i] = []int{0, 0, 0}
	}

	new_elevator := Elevator{elev_info.IPADDR, elev_info.POSITION * 2, elev_info.DIRECTION, elev_info.DESTINATION, new_ordermatrix}
	Active_elevators[elev_info.IPADDR] = new_elevator
	fmt.Println("Elevator", Active_elevators[elev_info.IPADDR].IPADDR, "online\n")
}

func removeElevator(ipaddr string) {
	orders_to_dist := Active_elevators[ipaddr].ORDER_MATRIX
	delete(Active_elevators, ipaddr)

	for floor := 0; floor < N_FLOORS; floor++ {
		for button_type := 0; button_type < 2; button_type++ {
			if orders_to_dist[floor][button_type] == 1 {
				appendOrder(button_type, floor)
			}
		}
	}
	fmt.Println("Deleted", ipaddr, "\n")
}

func appendOrder(button_type int, button_floor int) string {
	var button_dir string
	var optimal_elevatorIP string
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

		internal_orders[button_floor] = byte(1)

		_, w_err := order_file.WriteAt(internal_orders, 0)
		if w_err != nil {
			fmt.Println("write error:", w_err)
		}

		Active_elevators[my_ipaddr] = temp_elev
		return "nil"
	}

	for ipaddr := range Active_elevators {
		new_cost := costFunction(ipaddr, button_floor, button_dir)

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

	temp_elev := Active_elevators[optimal_elevatorIP]
	temp_elev.ORDER_MATRIX[button_floor][button_type] = 1
	Active_elevators[optimal_elevatorIP] = temp_elev

	return optimal_elevatorIP
}

func deleteOrder(ipaddr string, floor int) {
	temp_elev := Active_elevators[ipaddr]
	for i := 0; i < 3; i++ {
		temp_elev.ORDER_MATRIX[floor][i] = 0
	}
	Active_elevators[ipaddr] = temp_elev

	internal_orders[floor] = byte(0)
	_, w_err := order_file.WriteAt(internal_orders, 0)
	if w_err != nil {
		fmt.Println("write error:", w_err)
	}
}

func costFunction(elevator_ip string, order_floor int, button_dir string) int {
	cost := 0
	current_elevator := Active_elevators[elevator_ip]

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

func processNewInfo(c_router_info chan []byte, c_peerListUpdate chan string, c_queMan_output chan []byte, c_queMan_ack_order chan []byte) { //, c_pos_from_statemachine chan int, c_dir_from_statemachine chan int){
	var elev_info ElevInfo
	var last_info ElevInfo
	var button_order ElevInfo

	acknowledgeTimer = time.NewTimer(250 * time.Millisecond)
	acknowledgeTimer.Stop()

	for {
		select {
		case encoded_elev_info := <-c_router_info:
			err := json.Unmarshal(encoded_elev_info, &elev_info)
			if err != nil {
				fmt.Println("queMan Unmarshal error: ", err)
			}
			if _, in_list := Active_elevators[elev_info.IPADDR]; !in_list {
				appendElevator(elev_info)
			}
			if elev_info.F_NEW_INFO && (elev_info != last_info) {
				updateActiveElevators(elev_info)

				if elev_info.F_BUTTONPRESS {
					button_order = handleButtonpress(elev_info, c_queMan_output, c_queMan_ack_order)

				} else if elev_info.POSITION == elev_info.DESTINATION {
					deleteOrder(elev_info.IPADDR, elev_info.POSITION)

					for i := 0; i < 3; i++ {
						button_output := Output{0, 0, i, elev_info.POSITION, 0, -1}
						sendButtonLamp(button_output, c_queMan_output)
					}
				}
				last_info = elev_info

			} else {
				last_info = elev_info
			}

		case peerUpdate := <-c_peerListUpdate:
			removeElevator(peerUpdate)

		case <-acknowledgeTimer.C:
			if button_order.IPADDR == my_ipaddr {
				temp_elev := Active_elevators[my_ipaddr]
				temp_elev.ORDER_MATRIX[button_order.BUTTONFLOOR][button_order.BUTTON_TYPE] = 1
				Active_elevators[my_ipaddr] = temp_elev
				fmt.Printf("Acknowledge deadline reached. Taking order\n")
			}
		}
	}
}

func findNewDestination(c_to_statemachine chan int) {
	var dest int
	var pos_floor int
	for {
		switch {
		case Active_elevators[my_ipaddr].DIRECTION == 1:
			pos := Active_elevators[my_ipaddr].POSITION
			if Active_elevators[my_ipaddr].POSITION%2 == 0 {
				pos_floor = pos / 2
			} else {
				pos_floor = (pos+1)/2 - 1
			}
			for i := pos_floor + 1; i < N_FLOORS; i++ {
				if Active_elevators[my_ipaddr].ORDER_MATRIX[i][0] == 1 && i < Active_elevators[my_ipaddr].DESTINATION {
					dest = i
					temp_elev := Active_elevators[my_ipaddr]
					temp_elev.DESTINATION = dest
					Active_elevators[my_ipaddr] = temp_elev
					c_to_statemachine <- dest
					printActiveElevators()
					break

				} else if pos_floor == Active_elevators[my_ipaddr].DESTINATION && Active_elevators[my_ipaddr].ORDER_MATRIX[i][0] == 1 && i > pos_floor {
					dest = i
					temp_elev := Active_elevators[my_ipaddr]
					temp_elev.DESTINATION = dest
					Active_elevators[my_ipaddr] = temp_elev
					c_to_statemachine <- dest
					printActiveElevators()
					break
				}
			}

		case Active_elevators[my_ipaddr].DIRECTION == -1:
			pos := Active_elevators[my_ipaddr].POSITION
			if Active_elevators[my_ipaddr].POSITION%2 == 0 {
				pos_floor = pos / 2
			} else {
				pos_floor = (pos + 1) / 2
			}
			for i := pos_floor - 1; i >= 0; i-- {
				if Active_elevators[my_ipaddr].ORDER_MATRIX[i][1] == 1 && i > Active_elevators[my_ipaddr].DESTINATION {
					dest = i
					temp_elev := Active_elevators[my_ipaddr]
					temp_elev.DESTINATION = dest
					Active_elevators[my_ipaddr] = temp_elev
					c_to_statemachine <- dest
					printActiveElevators()
					break

				} else if pos_floor == Active_elevators[my_ipaddr].DESTINATION && Active_elevators[my_ipaddr].ORDER_MATRIX[i][1] == 1 && i < pos_floor {
					dest = i
					temp_elev := Active_elevators[my_ipaddr]
					temp_elev.DESTINATION = dest
					Active_elevators[my_ipaddr] = temp_elev
					c_to_statemachine <- dest
					printActiveElevators()
					break
				}
			}

		case Active_elevators[my_ipaddr].DIRECTION == 0:
			pos_floor = Active_elevators[my_ipaddr].POSITION / 2
			for i := 0; i < (N_FLOORS); i++ {
				if Active_elevators[my_ipaddr].ORDER_MATRIX[i][0] == 1 {
					dest = i

					temp_elev := Active_elevators[my_ipaddr]
					temp_elev.DESTINATION = dest
					Active_elevators[my_ipaddr] = temp_elev
					c_to_statemachine <- dest

					printActiveElevators()
					break

				} else if Active_elevators[my_ipaddr].ORDER_MATRIX[i][1] == 1 {
					dest = i
					temp_elev := Active_elevators[my_ipaddr]
					temp_elev.DESTINATION = dest
					Active_elevators[my_ipaddr] = temp_elev
					c_to_statemachine <- dest
					printActiveElevators()
					break
				}
			}
		}
		time.Sleep(10 * time.Millisecond)
	}

}

func sendButtonLamp(button_output Output, channel chan []byte) {
	encoded_output, err2 := json.Marshal(button_output)
	if err2 != nil {
		fmt.Println("QM button lamp JSON error: ", err2)
	}
	channel <- encoded_output
}

func sendElev(info ElevInfo, channel chan<- []byte) {
	encoded_output, err := json.Marshal(info)
	if err != nil {
		fmt.Println("SM JSON error: ", err)
	}
	channel <- encoded_output

}

func handleButtonpress(elev_info ElevInfo, c_queMan_output chan []byte, c_queMan_ack_order chan []byte) ElevInfo {
	var button_order ElevInfo

	if !elev_info.F_ACK_ORDER {
		optimal_ip := appendOrder(elev_info.BUTTON_TYPE, elev_info.BUTTONFLOOR)
		button_output := Output{0, 0, elev_info.BUTTON_TYPE, elev_info.BUTTONFLOOR, 1, -1}
		sendButtonLamp(button_output, c_queMan_output)
		switch {
		case elev_info.IPADDR == my_ipaddr:
			if optimal_ip != my_ipaddr {
				acknowledgeTimer.Reset(250 * time.Millisecond)
				button_order = elev_info
				button_order.IPADDR = optimal_ip
			}

		case elev_info.IPADDR != my_ipaddr:
			if optimal_ip == my_ipaddr {
				button_order = elev_info
				button_order.F_ACK_ORDER = true
				sendElev(button_order, c_queMan_ack_order)
				fmt.Printf("Got an optimal order, taking %d \n", button_order.BUTTONFLOOR)
			}
		}

	} else {
		acknowledgeTimer.Stop()
		button_order = elev_info
	}
	return button_order
}

func updateActiveElevators(elev_info ElevInfo) {
	temp_elev := Active_elevators[elev_info.IPADDR]
	temp_elev.DIRECTION = elev_info.DIRECTION

	if elev_info.MOVING {
		temp_elev.POSITION = elev_info.POSITION*2 + elev_info.DIRECTION
		if temp_elev.POSITION == -1 {
			temp_elev.POSITION = 0
		} else if temp_elev.POSITION == N_FLOORS*2-1 {
			temp_elev.POSITION = elev_info.POSITION * 2
		}
	} else {
		temp_elev.POSITION = elev_info.POSITION * 2
	}

	if elev_info.IPADDR != my_ipaddr {
		temp_elev.DESTINATION = elev_info.DESTINATION
	}

	Active_elevators[elev_info.IPADDR] = temp_elev
}

// For debugging purposes only
func printActiveElevators() {
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
