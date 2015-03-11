package queue

import(
	"fmt"
)

type  Elevator struct{
	IPADDR string
	POSITION int
	/*
	   Etg.			Pos. nr.
	    1 ............ 1
	  	  ............ 2
		2 ............ 3
		  ............ 4
		3 ............ 5
		  ............ 6
		4 ............ 7
	*/

	DIRECTION int
	/*
		opp = 1, ned = -1, stillestående = 0 
	*/

	DESTINATION int

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
var active_elevators =  make(map[string]Elevator)



func InitQueuemanager(ipaddr string) {
	my_ordermatrix := make([][]int, N_FLOORS)
	for i := 0; i < N_FLOORS; i++{
		my_ordermatrix[i] = []int{0,0,0}
	}
	new_elevator := Elevator{ipaddr, 0, 0, 0, my_ordermatrix} 
	active_elevators[ipaddr] = new_elevator
	fmt.Println("Elevator", active_elevators[ipaddr].IPADDR, "online")
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

func AppendOrder() {

	
}

func costFunction(elevator_ip string, order_floor int, button_dir string) int{
	cost := 0
	current_elevator := active_elevators[elevator_ip]

	//Omregner etg. nr. til posisjonsnr. (Ihht. structen Elevator)
	order_floor_pos := order_floor + (order_floor - 1)
	dest_pos := current_elevator.DESTINATION + (current_elevator.DESTINATION - 1)

	// Sjekker alle utfall hvor bestillingsknapp "opp" er trykket
	if button_dir == "up" {
		
		if  current_elevator.DIRECTION == 1 && dest_pos >=  order_floor_pos {
			cost = dest_pos - current_elevator.POSITION
		}
	}

	return cost
}
