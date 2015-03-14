package main

import (
	"./queue"
	// "fmt"
)

func main() {
	my_ipaddr := "111.11.111"
	button_type := 0
	button_floor := 2

	// position := 4
	// direction := 0
	// destination_floor := 3
	// order_floor := 2
	// button_dir := "down"



	queue.InitQueuemanager(my_ipaddr)
	queue.AppendElevator("999.99.999")

	// fmt.Println(queue.Active_elevators[my_ipaddr])
	queue.PrintActiveElevators()
	//queue.SetElevator(ipaddr, position, direction, destination_floor)
	
	queue.AppendOrder(button_type, button_floor)

	// fmt.Println(queue.Active_elevators[my_ipaddr])
	// queue.PrintActiveElevators()

	// queue.RemoveElevator(ipaddr)
	queue.PrintActiveElevators()
	// fmt.Println("Cost:", queue.CostFunction(my_ipaddr, order_floor, button_dir))

}
