package queue

import(
	"fmt"
)

type  Elevator struct{
	IPADDR string
	POSITION int
	DIRECTION int
	ORDER_MATRIX [][]int

	// 			inne_i_heis    ned       opp				Settes til 1 ved en ordre
	// 1.etg	[[  0          0         0]
	// 2.etg 	[   0          0         0]
	// 3.etg 	[   0          0         0]
	// 4.etg	[   0          0         0]]
}


var active_elevators =  make(map[string]Elevator)


func InitQueuemanager(ipaddr string) {
	my_ordermatrix := make([][]int, 4)
	for i := 0; i < 4; i++{
		my_ordermatrix[i] = []int{0,0,0}
	}
	new_elevator := Elevator{ipaddr, 0, 0, my_ordermatrix} 
	active_elevators[ipaddr] = new_elevator
	fmt.Println("Elevator", active_elevators[ipaddr].IPADDR, "online")
}


func AppendElevator(ipaddr string) {
	new_ordermatrix := make([][]int, 4)
	for i := 0; i < 4; i++{
		new_ordermatrix[i] = []int{0,0,0}
	}
	new_elevator := Elevator{ipaddr, 0, 0, new_ordermatrix}
	active_elevators[ipaddr] = new_elevator
	fmt.Println("Elevator", active_elevators[ipaddr].IPADDR, "online")}



func PrintActiveElevators() {
	for i := range(active_elevators){
		fmt.Println("Elevator:",active_elevators[i].IPADDR)
		for floor := 0; floor < 4; floor++ {	
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