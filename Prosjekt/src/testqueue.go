package main

import (
	"./queue"
	// "fmt"
)

func main() {
	my_ipaddr := "111.11.111"
	other_ipaddr := "999.99.999"

	queue.InitQueuemanager(my_ipaddr)
	queue.AppendElevator(other_ipaddr)

	// fmt.Println(queue.Active_elevators[my_ipaddr])
	queue.PrintActiveElevators()

	position := 0
	direction := 0
	destination_floor := 0
	queue.SetElevator(my_ipaddr, position, direction, destination_floor)
	
	position = 3
	direction = 1
	destination_floor = 6
	queue.SetElevator(other_ipaddr, position, direction, destination_floor)

	queue.PrintActiveElevators()


	button_type := 0
	button_floor := 2
	queue.AppendOrder(button_type, button_floor)

	queue.PrintActiveElevators()


	// fmt.Println(queue.Active_elevators[my_ipaddr])
	// queue.PrintActiveElevators()

	// queue.RemoveElevator(ipaddr)
	// queue.PrintActiveElevators()
	// fmt.Println("Cost:", queue.CostFunction(my_ipaddr, order_floor, button_dir))

}
