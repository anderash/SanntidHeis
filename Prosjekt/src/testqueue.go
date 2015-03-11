package main

import(
	"./queue"
	// "fmt"
)

func main() {
	ipaddr := "123.45.345"
	queue.InitQueuemanager(ipaddr)
	queue.AppendElevator("999.99.999")
	queue.PrintActiveElevators()

	queue.RemoveElevator(ipaddr)
	queue.PrintActiveElevators()


}