package main

import (
	"./queue"
	"fmt"
)

func main() {
	ipaddr := "123.45.345"
	position := 4
	direction := 0
	destination_floor := 3
	order_floor := 2
	button_dir := "down"

	queue.InitQueuemanager(ipaddr)

	queue.PrintActiveElevators()
	queue.SetElevator(ipaddr, position, direction, destination_floor)

	queue.PrintActiveElevators()

	// queue.RemoveElevator(ipaddr)
	// queue.PrintActiveElevators()
	fmt.Println("Cost:", queue.CostFunction(ipaddr, order_floor, button_dir))

}
