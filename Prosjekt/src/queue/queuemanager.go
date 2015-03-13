package queue

import(
	"fmt"
)

type  Elevator struct{
	IPADDR string
	POSITION int
	/*
	   Etg.			Pos. nr.
	    1 ............ 0
	  	  ............ 1
		2 ............ 2
		  ............ 3
		3 ............ 4
		  ............ 5
		4 ............ 6
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
	/* 			inne_i_heis    ned       opp				Settes til 1 ved en ordre
	   1.etg	[[  0          0         0]
	   2.etg 	 [  0          0         0]
	   3.etg 	 [  0          0         0]
	   4.etg	 [  0          0         0]]
	   osv.
	*/
}

const(
	N_FLOORS = 4
 	N_POSITIONS = N_FLOORS + (N_FLOORS-1)
)

// Indexen i map'en er ip-adressen til den aktuelle heisen
var active_elevators = make(map[string]Elevator)




func InitQueuemanager(ipaddr string) {
	my_ordermatrix := make([][]int, N_FLOORS)
	for i := 0; i < N_FLOORS; i++{
		my_ordermatrix[i] = []int{0,0,0}
	}
	new_elevator := Elevator{ipaddr, 0, 0, 0, my_ordermatrix} 
	active_elevators[ipaddr] = new_elevator
	fmt.Println("Elevator", active_elevators[ipaddr].IPADDR, "online\n")
}


// Denne funkjsonen brukes kun ifm debugging
func SetElevator(ipaddr string, position int, direction int, destinasjon_pos int){
	temp := active_elevators[ipaddr]
	temp.POSITION = position
	temp.DIRECTION = direction
	temp.DESTINATION = destinasjon_pos
	active_elevators[ipaddr] = temp

}

func AppendElevator(ipaddr string) {
	new_ordermatrix := make([][]int, N_FLOORS) 
	for i := 0; i < N_FLOORS; i++{
		new_ordermatrix[i] = []int{0,0,0}
	}
	new_elevator := Elevator{ipaddr, 0, 0, 0, new_ordermatrix}
	active_elevators[ipaddr] = new_elevator
	fmt.Println("Elevator", active_elevators[ipaddr].IPADDR, "online")}



func PrintActiveElevators() {
	for i := range(active_elevators){
		fmt.Println("Elevator:",active_elevators[i].IPADDR)
		fmt.Println("Position:", active_elevators[i].POSITION)
		fmt.Println("Direction:", active_elevators[i].DIRECTION)
		fmt.Println("Destination:", active_elevators[i].DESTINATION)
		fmt.Printf("Orders:\n")
		for floor := 0; floor < N_FLOORS; floor++ {	
			fmt.Println("Floor", floor + 1, ":", active_elevators[i].ORDER_MATRIX[floor])
		}
		fmt.Println("\n")	
	}
}


// Trenger også å distribuere alle ordrene til heisen som skal slettes til de andre heisene
func  RemoveElevator(ipaddr string) {
	delete(active_elevators, ipaddr)
	fmt.Println("Deleting", ipaddr, "\n")
}

// Bruker kostfunksjonen for å legge til ny ordre
func AppendOrder(button_type int, button_floor int) {
	if button_type == 0 {
		button_dir := "up"
	} else if button_type == 1 {
		button_dir := "down"
	} else {

	}

	cost := 100

	for i := range(active_elevators){
		if 
	}	
}



func CostFunction(elevator_ip string, order_floor int, button_dir string) int{
	cost := 0
	current_elevator := active_elevators[elevator_ip]

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
			if current_elevator.DESTINATION >= order_floor_pos {
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
			if current_elevator.DESTINATION <= order_floor_pos {
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


func ProsseserNyinfo( ){ //Tar inn kanal_fra_heis
/*
	Oppdaterer ny info fra de andre heisene
*/

}